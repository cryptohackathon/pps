package subcommands

import (
	"fmt"
	"math/big"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/ZenGo-X/fe-hackaton-demo/internal/data"
	"github.com/ZenGo-X/fe-hackaton-demo/internal/gofe"
	"github.com/ZenGo-X/fe-hackaton-demo/internal/recipient"
	"github.com/ZenGo-X/fe-hackaton-demo/internal/rounds"
)

var (
	searchArgs struct {
		party, from, to int
	}

	Search = cli.Command{
		Action: search,
		Name:   "search",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "party",
				Usage:       "Number of party",
				Required:    true,
				Destination: &searchArgs.party,
			},
			&cli.IntFlag{
				Name:        "from",
				Usage:       "Number of round when party was online last time",
				Required:    true,
				Destination: &searchArgs.from,
			},
			&cli.IntFlag{
				Name:        "to",
				Usage:       "Number of round when party went online",
				Destination: &searchArgs.to,
				DefaultText: "last round",
			},
		},
	}
)

func search(_ *cli.Context) error {
	party, err := recipient.LoadRecipient("stand/parties", searchArgs.party)
	if err != nil {
		return errors.Wrap(err, "load party secret")
	}

	repo, err := rounds.OpenRepository("stand/repo")
	if err != nil {
		return errors.Wrap(err, "open repository")
	}

	mpk, err := repo.GetMPK()
	if err != nil {
		return errors.Wrap(err, "couldn't retrieve MPK")
	}

	t1 := searchArgs.from
	var v1 *big.Int
	if t1 > 0 {
		ciphertext1, err := repo.GetRound(t1)
		if err != nil {
			return errors.Wrapf(err, "retrieve round %d", t1)
		}
		v1, err = gofe.Decrypt(mpk, party.Secret, ciphertext1)
		if err != nil {
			return errors.Wrapf(err, "decrypt ciphertext from round %d", t1)
		}
	} else {
		v1 = big.NewInt(0)
	}

	t2 := searchArgs.to
	var ciphertext2 *data.Ciphertext
	if searchArgs.to != 0 {
		ciphertext2, err = repo.GetRound(t2)
		if err != nil {
			return errors.Wrapf(err, "retrieve round %d", t2)
		}
	} else {
		t2, ciphertext2, err = repo.GetLastRound()
		if err != nil {
			return errors.Wrap(err, "retrieve last round")
		}
	}
	v2, err := gofe.Decrypt(mpk, party.Secret, ciphertext2)
	if err != nil {
		return errors.Wrapf(err, "decrypt ciphertext from round %d", t2)
	}

	if v1.Cmp(v2) == 0 {
		fmt.Printf("Party received no signal within rounds [%d;%d]\n", t1, t2)
		return nil
	}

	fmt.Println("Party received signal(s)!")
	for {
		fmt.Printf("Searching received signal within rounds [%d;%d]\n", t1, t2)
		ti, err := findFirstSignal(party, repo, mpk, t1, v1, t2)
		if err != nil {
			return errors.Wrap(err, "search failed")
		}
		fmt.Printf("Received signal at round %d!\n", ti+1)

		ciphertext, err := repo.GetRound(ti + 1)
		if err != nil {
			return errors.Wrapf(err, "cannot retrieve round %d", ti+1)
		}
		vi, err := gofe.Decrypt(mpk, party.Secret, ciphertext)
		if err != nil {
			return errors.Wrapf(err, "decrypt ciphertext from round %d", ti+1)
		}

		if vi.Cmp(v2) == 0 {
			fmt.Println("No more signals available")
			return nil
		}
		fmt.Println("More signals available!")
		t1 = ti + 1
		v1 = vi
	}
}

func findFirstSignal(party *recipient.Party, repo *rounds.Repository, mpk data.MPK, t1 int, v1 *big.Int, t2 int) (int, error) {
	if t1 == t2 {
		return t1, nil
	}

	m := (t1 + t2) / 2
	if ((t1 + t2) & 1) == 1 {
		m += 1
	}
	ciphertext, err := repo.GetRound(m)
	if err != nil {
		return 0, errors.Wrapf(err, "cannot retrieve round %d", m)
	}
	vm, err := gofe.Decrypt(mpk, party.Secret, ciphertext)
	if err != nil {
		return 0, errors.Wrapf(err, "decrypt ciphertext from round %d", m)
	}

	if v1.Cmp(vm) == 0 {
		fmt.Printf("Accessing round %d... v_%d == v_%d\n", m, t1, m)
		return findFirstSignal(party, repo, mpk, m, vm, t2)
	} else {
		fmt.Printf("Accessing round %d... v_%d != v_%d\n", m, t1, m)
		return findFirstSignal(party, repo, mpk, t1, v1, m-1)
	}
}

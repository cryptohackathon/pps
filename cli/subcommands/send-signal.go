package subcommands

import (
	"fmt"
	"math/big"

	gofe "github.com/fentec-project/gofe/data"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	gofe2 "github.com/ZenGo-X/fe-hackaton-demo/internal/gofe"
	"github.com/ZenGo-X/fe-hackaton-demo/internal/rounds"
)

var (
	recipientParty int

	SendSignal = cli.Command{
		Action: sendSignal,
		Name:   "send-signal",
		Usage:  "Sends encrypted signal to recipient",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "party",
				Usage:       "Recipient of the signal `j` (1 <= j <= N)",
				Destination: &recipientParty,
				Required:    true,
			},
		},
	}
)

func sendSignal(_ *cli.Context) error {
	repo, err := rounds.OpenRepository("stand/repo")
	if err != nil {
		return errors.Wrap(err, "cannot open repository")
	}

	mpk, err := repo.GetMPK()
	if err != nil {
		return errors.Wrap(err, "cannot retrieve MPK")
	}

	if recipientParty <= 0 || recipientParty > mpk.DDH.Params.L {
		return errors.Errorf("expected recipient in range [1; %d]", mpk.DDH.Params.L)
	}

	n, previousCiphertext, err := repo.GetLastRound()
	if err != nil {
		return errors.Wrap(err, "cannot retrieve last round")
	}

	plaintext := gofe.NewConstantVector(mpk.DDH.Params.L, big.NewInt(0))
	plaintext[recipientParty-1] = big.NewInt(1)

	ciphertext, err := gofe2.Encrypt(mpk, plaintext)
	if err != nil {
		return errors.Wrap(err, "can't encrypt a signal")
	}

	if n > 0 {
		err = ciphertext.Mul(previousCiphertext)
		if err != nil {
			return errors.Wrap(err, "calculating ciphertext * previousCiphertext")
		}
	}

	err = repo.PublishRound(n+1, &ciphertext)
	if err != nil {
		return errors.Wrap(err, "publish encrypted signal error")
	}

	fmt.Printf("You successfully sent encrypted signal to party %d in round %d!\n", recipientParty, n+1)
	return nil
}

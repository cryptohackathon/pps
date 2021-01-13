package subcommands

import (
	"fmt"
	"github.com/ZenGo-X/fe-hackaton-demo/internal/recipient"
	"github.com/ZenGo-X/fe-hackaton-demo/internal/rounds"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/ZenGo-X/fe-hackaton-demo/internal/gofe"
)

var (
	keygenParties int

	Keygen = cli.Command{
		Action: keygen,
		Name:   "keygen",
		Usage:  "Runs key generation & derivation",
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "parties",
				Usage:       "Total amount of recipients `N` (N >= 2)",
				Destination: &keygenParties,
				Required:    true,
			},
		},
	}
)

func keygen(_ *cli.Context) error {
	if keygenParties < 2 {
		return errors.New("expected at least 2 parties!")
	}
	mpk, sk, err := gofe.GenerateMasterKeys(keygenParties)
	if err != nil {
		return errors.Wrap(err, "keygen failed")
	}

	for j, skj := range sk {
		party := &recipient.Party{Secret: skj}
		err := party.SaveRecipient("stand/parties")
		if err != nil {
			return errors.Wrapf(err, "cannot save party %d", j+1)
		}
	}

	_, err = rounds.NewEmptyRepository("stand/repo", mpk)
	if err != nil {
		return errors.Wrap(err, "cannot create empty repository")
	}

	fmt.Println("Keygen completed!")

	return nil
}

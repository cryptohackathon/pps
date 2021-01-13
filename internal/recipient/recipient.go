package recipient

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/pkg/errors"

	"github.com/ZenGo-X/fe-hackaton-demo/internal/data"
)

// Party holding RecipientSecretKey which can be used
// to decrypt a signal
type Party struct {
	Secret data.RecipientSecretKey
}

// Saves party secret key at `{dir}/party_{i}.json`
//
// It's an error if this file already exist
func (p *Party) SaveRecipient(dir string) error {
	err := os.MkdirAll(dir, 0770)
	if err != nil {
		return errors.Wrap(err, "create dir")
	}
	filepath := path.Join(dir, fmt.Sprintf("party_%d.json", p.Secret.I+1))
	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return errors.Wrap(err, "create file")
	}

	err = json.NewEncoder(file).Encode(p)
	if err != nil {
		_ = file.Close()
		return errors.Wrap(err, "write/encode party secret key")
	}

	if err = file.Close(); err != nil {
		return errors.Wrap(err, "close file")
	}

	return nil
}

// Loads party secret key from `{path}/party_{i}.json`
func LoadRecipient(dir string, i int) (*Party, error) {
	filepath := path.Join(dir, fmt.Sprintf("party_%d.json", i))
	file, err := os.Open(filepath)
	if err != nil {
		return nil, errors.Wrap(err, "open file")
	}
	defer func() {
		_ = file.Close()
	}()

	var party Party
	err = json.NewDecoder(file).Decode(&party)
	if err != nil {
		return nil, errors.Wrap(err, "decode/read file")
	}

	return &party, nil
}

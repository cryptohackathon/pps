package rounds

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path"

	"github.com/ZenGo-X/fe-hackaton-demo/internal/data"
)

// Manage access to rounds
//
// In real worlds it would use blockchain, but in demo it
// just stores everything at filesystem.
type Repository struct {
	path string
}

// Creates new empty repository in directory `path`
//
// Will create `path` directory (if it isn't present) and file
// `{path}/round0.json` containing `mpk`
func NewEmptyRepository(dir string, mpk data.MPK) (*Repository, error) {
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return nil, errors.Wrap(err, "create dir")
	}

	round0 := path.Join(dir, "round_0.json")
	file, err := os.OpenFile(round0, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666)
	if err != nil {
		return nil, errors.Wrap(err, "create file round0.json")
	}

	err = json.NewEncoder(file).Encode(&mpk)
	if err != nil {
		_ = file.Close()
		return nil, errors.Wrap(err, "write/encode mpk")
	}

	if err = file.Close(); err != nil {
		return nil, errors.Wrap(err, "close file round0.json")
	}

	return &Repository{path: dir}, nil
}

// Tries to open existing repository
func OpenRepository(path string) (*Repository, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return nil, errors.Wrap(err, "cannot stat repo")
	}

	if !fileInfo.IsDir() {
		return nil, errors.New("repo is not a directory")
	}
	return &Repository{path: path}, nil
}

// Retrieves i-th round from repository
func (r *Repository) GetRound(i int) (*data.Ciphertext, error) {
	filename := path.Join(r.path, fmt.Sprintf("round_%d.json", i))
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.Wrap(err, "open round.json")
	}
	defer func() {
		_ = file.Close()
	}()

	var ciphertext data.Ciphertext
	err = json.NewDecoder(file).Decode(&ciphertext)
	if err != nil {
		return nil, errors.Wrap(err, "decode/read round")
	}

	return &ciphertext, nil
}

// Retrieves the last published round from repository
//
// If no rounds present, it'll return (0, nil, nil)
func (r *Repository) GetLastRound() (int, *data.Ciphertext, error) {
	n := 0
	var ciphertext data.Ciphertext

	for ; ; n++ {
		filename := path.Join(r.path, fmt.Sprintf("round_%d.json", n+1))
		file, err := os.Open(filename)
		if err != nil && os.IsNotExist(err) {
			break
		} else if err != nil {
			return 0, nil, errors.Wrap(err, "unexpected error while retrieving round")
		}

		err = json.NewDecoder(file).Decode(&ciphertext)
		_ = file.Close()
		if err != nil {
			return 0, nil, errors.Wrapf(err, "malformed round %d", n+1)
		}
	}
	if n == 0 {
		return 0, nil, nil
	}
	return n, &ciphertext, nil
}

// Publishes a new round into repository
//
// Creates file `{repository}/round_{n}.json`. It's an error if this file
// already exist.
func (r *Repository) PublishRound(n int, ciphertext *data.Ciphertext) error {
	filename := path.Join(r.path, fmt.Sprintf("round_%d.json", n))
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666)
	if err != nil {
		return errors.Wrap(err, "create file for round")
	}

	err = json.NewEncoder(file).Encode(ciphertext)
	if err != nil {
		_ = file.Close()
		return errors.Wrap(err, "encode/write to file")
	}
	if err = file.Close(); err != nil {
		return errors.Wrap(err, "close file")
	}

	return nil
}

// Retrieves master public key
func (r *Repository) GetMPK() (data.MPK, error) {
	filename := path.Join(r.path, "round_0.json")
	file, err := os.Open(filename)
	if err != nil {
		return data.MPK{}, errors.Wrap(err, "open round.json")
	}
	defer func() {
		_ = file.Close()
	}()

	var mpk data.MPK
	err = json.NewDecoder(file).Decode(&mpk)
	if err != nil {
		return data.MPK{}, errors.Wrap(err, "decode/read round")
	}

	return mpk, nil
}

package data

import (
	"github.com/fentec-project/gofe/innerprod/simple"
	"math/big"

	"github.com/pkg/errors"

	gofe "github.com/fentec-project/gofe/data"
)

type MPK struct {
	DDH    *simple.DDH
	Vector gofe.Vector
}

type Ciphertext struct {
	Vector gofe.Vector
}

// Performs ciphertext *= anotherCiphertext
//
// Might result in error, e.g. if adding ciphertext of different length.
// In this case, ciphertext is not modified.
func (c *Ciphertext) Mul(another *Ciphertext) error {
	v1 := ([]*big.Int)(c.Vector)
	v2 := ([]*big.Int)(another.Vector)

	if len(v1) != len(v2) {
		return errors.New("given ciphertexts have different lengths")
	}

	for i, x := range another.Vector {
		c.Vector[i] = new(big.Int).Mul(c.Vector[i], x)
	}

	return nil
}

type RecipientSecretKey struct {
	I          int
	DerivedKey *big.Int
}

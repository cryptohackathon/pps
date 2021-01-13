package gofe

import (
	"math/big"

	gofe "github.com/fentec-project/gofe/data"
	"github.com/fentec-project/gofe/innerprod/simple"
	"github.com/pkg/errors"

	"github.com/ZenGo-X/fe-hackaton-demo/internal/data"
)

func GenerateMasterKeys(parties int) (data.MPK, []data.RecipientSecretKey, error) {
	ddh, err := simple.NewDDH(parties, 512, big.NewInt(1024))
	if err != nil {
		return data.MPK{}, nil, errors.Wrap(err, "scheme could not be properly configured")
	}
	return GenerateMasterKeysDDH(ddh)
}

func GenerateMasterKeysDDH(ddh *simple.DDH) (data.MPK, []data.RecipientSecretKey, error) {
	msk, mpk, err := ddh.GenerateMasterKeys()
	if err != nil {
		return data.MPK{}, nil, errors.Wrap(err, "generate master key")
	}

	secretKeys := make([]data.RecipientSecretKey, 0)
	for j := 0; j < ddh.Params.L; j++ {
		y := gofe.NewConstantVector(ddh.Params.L, big.NewInt(0))
		y[j] = big.NewInt(1)

		sk, err := ddh.DeriveKey(msk, y)
		if err != nil {
			return data.MPK{}, nil, errors.Wrapf(err, "generate sk for party %d", j+1)
		}
		secretKeys = append(secretKeys, data.RecipientSecretKey{I: j, DerivedKey: sk})
	}

	return data.MPK{DDH: ddh, Vector: mpk}, secretKeys, nil
}

func Encrypt(mpk data.MPK, vector gofe.Vector) (data.Ciphertext, error) {
	ciphertext, err := mpk.DDH.Encrypt(gofe.NewVector(vector), mpk.Vector)
	if err != nil {
		return data.Ciphertext{}, err
	}
	return data.Ciphertext{Vector: ciphertext}, nil
}

func Decrypt(mpk data.MPK, sk data.RecipientSecretKey, ciphertext *data.Ciphertext) (*big.Int, error) {
	y := gofe.NewConstantVector(mpk.DDH.Params.L, big.NewInt(0))
	y[sk.I] = big.NewInt(1)
	return mpk.DDH.Decrypt(ciphertext.Vector, sk.DerivedKey, y)
}

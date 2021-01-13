package gofe

import (
	gofe "github.com/fentec-project/gofe/data"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestGenerateMasterKeys(t *testing.T) {
	f := func(n int) func(t *testing.T) {
		return func(t *testing.T) {
			mpk, sk, err := GenerateMasterKeys(n)
			assert.NoError(t, err, "generate master keys, n", n)
			assert.NotNil(t, mpk.DDH, "ddh is missing")
			assert.Equal(t, mpk.DDH.Params.L, n, "wrong length of input vectors")
			assert.NotNil(t, mpk.Vector, "mpk vector is missing")
			assert.Len(t, sk, n, "wrong number of derived keys")
		}
	}
	t.Run("n=2", f(2))
	t.Run("n=3", f(3))
	t.Run("n=4", f(4))
	t.Run("n=5", f(5))
}

func TestEncryptDecrypt(t *testing.T) {
	mpk, sk, err := GenerateMasterKeys(5)
	assert.NoError(t, err, "keygen failed")

	plaintext := gofe.NewConstantVector(5, big.NewInt(0))
	plaintext[2] = big.NewInt(1)

	ciphertext, err := Encrypt(mpk, plaintext)
	assert.NoError(t, err, "encrypt plaintext")

	v, err := Decrypt(mpk, sk[2], ciphertext)
	assert.NoError(t, err, "decrypt using sk2")
	assert.Equal(t, big.NewInt(1), v, "incorrectly decrypted v using sk2")

	for j := 0; j < 5; j++ {
		if j == 2 {
			continue
		}
		v, err = Decrypt(mpk, sk[j], ciphertext)
		assert.NoError(t, err, "decrypt using sk", j)
		assert.Equal(t, big.NewInt(0), v, "incorrectly decrypted v using sk", j)
	}

}

func TestCiphertextIsMultiplicative(t *testing.T) {
	mpk, sk, err := GenerateMasterKeys(2)
	assert.NoError(t, err, "keygen failed")

	// x1 = [1 0]
	x1 := gofe.NewConstantVector(2, big.NewInt(0))
	x1[0] = big.NewInt(1)
	// x2 = [0 1]
	x2 := gofe.NewConstantVector(2, big.NewInt(0))
	x2[1] = big.NewInt(1)

	// e1 = Encrypt(mpk, x1)
	e1, err := Encrypt(mpk, x1)
	assert.NoError(t, err, "encrypt x1")

	// e2 = e1 * Encrypt(mpk, x2)
	e2, err := Encrypt(mpk, x2)
	assert.NoError(t, err, "encrypt x1")
	err = e2.Mul(&e1)
	assert.NoError(t, err, "e1*e2")

	// check that Decrypt(mpk, sk0, e1) == 1
	v1Sk0, err := Decrypt(mpk, sk[0], e1)
	assert.NoError(t, err, "decrypt e1 using sk0")
	assert.Equal(t, big.NewInt(1), v1Sk0, "incorrectly decrypted e1 using sk0")

	// check that Decrypt(mpk, sk0, e2) == 1
	v2Sk0, err := Decrypt(mpk, sk[0], e2)
	assert.NoError(t, err, "decrypt e1*e2 using sk0")
	assert.Equal(t, big.NewInt(1), v2Sk0, "incorrectly decrypted e1*e2 using sk0")
}

// Same as TestCiphertextIsMultiplicative but larger
func TestCiphertextIsMultiplicative2(t *testing.T) {
	mpk, sk, err := GenerateMasterKeys(5)
	assert.NoError(t, err, "keygen failed")

	plaintext1 := gofe.NewConstantVector(5, big.NewInt(0))
	plaintext2 := gofe.NewConstantVector(5, big.NewInt(0))
	plaintext3 := gofe.NewConstantVector(5, big.NewInt(0))

	plaintext1[2] = big.NewInt(1)
	plaintext2[3] = big.NewInt(1)
	plaintext3[2] = big.NewInt(1)

	ciphertext1, err := Encrypt(mpk, plaintext1)
	assert.NoError(t, err, "encrypt plaintext1")

	ciphertext2, err := Encrypt(mpk, plaintext2)
	assert.NoError(t, err, "encrypt plaintext2")
	err = ciphertext2.Mul(&ciphertext1)
	assert.NoError(t, err, "ciphertext1*ciphertext2")

	ciphertext3, err := Encrypt(mpk, plaintext3)
	assert.NoError(t, err, "encrypt plaintext3")
	err = ciphertext3.Mul(&ciphertext2)
	assert.NoError(t, err, "ciphertext2*ciphertext3")

	v1Sk2, err := Decrypt(mpk, sk[2], ciphertext1)
	assert.NoError(t, err, "decrypt ciphertext1 using sk2")
	assert.Equal(t, big.NewInt(1), v1Sk2)

	for j := 0; j < 5; j++ {
		if j == 2 {
			continue
		}
		v1Skj, err := Decrypt(mpk, sk[j], ciphertext1)
		assert.NoError(t, err, "decrypt ciphertext1 using skj", j)
		assert.Equal(t, big.NewInt(0), v1Skj)
	}

	v2Sk2, err := Decrypt(mpk, sk[2], ciphertext2)
	assert.NoError(t, err, "decrypt ciphertext2 using sk2")
	assert.Equal(t, big.NewInt(1), v2Sk2)

	v2Sk3, err := Decrypt(mpk, sk[3], ciphertext2)
	assert.NoError(t, err, "decrypt ciphertext2 using sk3")
	assert.Equal(t, big.NewInt(1), v2Sk3)

	for j := 0; j < 5; j++ {
		if j == 2 || j == 3 {
			continue
		}
		v2Skj, err := Decrypt(mpk, sk[j], ciphertext2)
		assert.NoError(t, err, "decrypt ciphertext2 using skj", j)
		assert.Equal(t, big.NewInt(0), v2Skj)
	}

	v3Sk2, err := Decrypt(mpk, sk[2], ciphertext3)
	assert.NoError(t, err, "decrypt ciphertext3 using sk2")
	assert.Equal(t, big.NewInt(2), v3Sk2)

	v3Sk3, err := Decrypt(mpk, sk[3], ciphertext3)
	assert.NoError(t, err, "decrypt ciphertext3 using sk3")
	assert.Equal(t, big.NewInt(1), v3Sk3)

	for j := 0; j < 5; j++ {
		if j == 2 || j == 3 {
			continue
		}
		v3Skj, err := Decrypt(mpk, sk[j], ciphertext3)
		assert.NoError(t, err, "decrypt ciphertext3 using skj", j)
		assert.Equal(t, big.NewInt(0), v3Skj)
	}
}

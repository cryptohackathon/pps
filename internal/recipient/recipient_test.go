package recipient

import (
	"io/ioutil"
	"math/big"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ZenGo-X/fe-hackaton-demo/internal/data"
)

func TestMarshalling(t *testing.T) {
	// Create fake party for test
	secret := big.NewInt(1234)
	party := Party{Secret: data.RecipientSecretKey{
		I:          1,
		DerivedKey: secret,
	}}

	// Create temp dir
	dir, err := ioutil.TempDir("", "parties_secrets")
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = os.RemoveAll(dir)
	}()

	// Tests
	t.Run("Save", func(t *testing.T) {
		err := party.SaveRecipient(dir)
		assert.NoError(t, err, "save party")
	})

	t.Run("Load", func(t *testing.T) {
		party2, err := LoadRecipient(dir, 1)
		assert.NoError(t, err, "load party")
		assert.Equal(t, &party, party2)
	})
}

package rounds

import (
	"io/ioutil"
	"math/big"
	"os"
	"path"
	"testing"

	gofe "github.com/fentec-project/gofe/data"
	"github.com/fentec-project/gofe/innerprod/simple"
	"github.com/stretchr/testify/assert"

	"github.com/ZenGo-X/fe-hackaton-demo/internal/data"
)

func TestRepository(t *testing.T) {
	var r *Repository
	mpk := data.MPK{
		DDH: &simple.DDH{Params: &simple.DDHParams{
			L:     4,
			Bound: big.NewInt(1),
			G:     big.NewInt(2),
			P:     big.NewInt(3),
			Q:     big.NewInt(4),
		}},
		Vector: gofe.NewVector([]*big.Int{
			big.NewInt(1),
			big.NewInt(2),
			big.NewInt(3),
		}),
	}

	// Create temp dir
	dir, err := ioutil.TempDir("", "parties_secrets")
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = os.RemoveAll(dir)
	}()
	repo := path.Join(dir, "repo")

	// Tests
	repositoryCreated := t.Run("Create empty repository", func(t *testing.T) {
		r, err = NewEmptyRepository(repo, mpk)
		assert.NoError(t, err)
	})

	t.Run("Retrieve the same MPK", func(t *testing.T) {
		if !repositoryCreated {
			t.Skip()
		}
		mpk2, err := r.GetMPK()
		assert.NoError(t, err, "get mpk")
		assert.Equal(t, mpk, mpk2)
	})

	t.Run("Open repo should not produce error", func(t *testing.T) {
		if !repositoryCreated {
			t.Skip()
		}
		_, err := OpenRepository(repo)
		assert.NoError(t, err)
	})

	t.Run("GetLastRound returns n=0", func(t *testing.T) {
		if !repositoryCreated {
			t.Skip()
		}
		n, ciphertext, err := r.GetLastRound()
		assert.NoError(t, err)
		assert.Equal(t, 0, n, "wrong round number")
		assert.Nil(t, ciphertext)
	})

	ciphertext1 := data.Ciphertext{Vector: gofe.NewConstantVector(5, big.NewInt(10))}
	ciphertext2 := data.Ciphertext{Vector: gofe.NewConstantVector(5, big.NewInt(6))}
	roundsPublished := t.Run("Publish rounds", func(t *testing.T) {
		if !repositoryCreated {
			t.Skip()
		}
		err := r.PublishRound(1, &ciphertext1)
		assert.NoError(t, err, "publish round1")
		err = r.PublishRound(2, &ciphertext2)
		assert.NoError(t, err, "publish round2")
	})

	t.Run("Overwrite round is error", func(t *testing.T) {
		if !repositoryCreated || !roundsPublished {
			t.Skip()
		}
		err := r.PublishRound(1, &ciphertext2)
		assert.Error(t, err)
	})

	t.Run("Retrieve rounds", func(t *testing.T) {
		if !repositoryCreated || !roundsPublished {
			t.Skip()
		}
		c1, err := r.GetRound(1)
		assert.NoError(t, err, "get round1")
		assert.Equal(t, &ciphertext1, c1)
		c2, err := r.GetRound(2)
		assert.NoError(t, err, "get round2")
		assert.Equal(t, &ciphertext2, c2)
	})

	t.Run("GetLastRound returns ciphertext2", func(t *testing.T) {
		if !repositoryCreated || !roundsPublished {
			t.Skip()
		}
		n, c, err := r.GetLastRound()
		assert.NoError(t, err, "get last round")
		assert.Equal(t, &ciphertext2, c)
		assert.Equal(t, 2, n)
	})
}

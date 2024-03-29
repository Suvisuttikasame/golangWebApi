//go:build unit

package authentication

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T) {
	t.Run("create jwt token success should return token", func(t *testing.T) {
		b := Body{
			Id:       uuid.New(),
			Username: "test",
			Email:    "test@mail.com",
		}
		var s string

		j, err := NewJWTToken([]byte("secret"))

		assert.Nil(t, err)

		tk, err := j.CreateToken(b)
		assert.Nil(t, err)
		assert.IsType(t, s, tk)

	})

	t.Run("validate jwt token success should return true", func(t *testing.T) {
		b := Body{
			Id:       uuid.New(),
			Username: "test",
			Email:    "test@mail.com",
		}

		j, err := NewJWTToken([]byte("secret"))

		assert.Nil(t, err)
		tk, err := j.CreateToken(b)
		assert.Nil(t, err)

		r, err := j.Verification(tk)
		assert.Nil(t, err)
		t.Log(r)

	})

}

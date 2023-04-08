//go:build unit

package authentication

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestPasetoToken(t *testing.T) {
	t.Run("create paseto token should return paseto token", func(t *testing.T) {
		body := Body{
			Id:       uuid.New(),
			Username: "test",
			Email:    "test@mail.com",
		}
		var s string

		p, err := NewPasetoToken([]byte("asdfgjhtuasdfgjhtuasdfgjhtuasdfg"))
		assert.Nil(t, err)

		token, err := p.CreateToken(body)
		t.Log(token)
		assert.Nil(t, err)
		assert.IsType(t, s, token)

	})

	t.Run("verify paseto token should return true", func(t *testing.T) {
		body := Body{
			Id:       uuid.New(),
			Username: "test",
			Email:    "test@mail.com",
		}
		p, err := NewPasetoToken([]byte("asdfgjhtuasdfgjhtuasdfgjhtuasdfg"))
		assert.Nil(t, err)

		token, err := p.CreateToken(body)
		assert.Nil(t, err)

		r, err := p.Verification(token)

		assert.Nil(t, err)
		t.Log(r)

	})

}

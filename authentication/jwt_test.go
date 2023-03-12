package authentication

import (
	"goApp/util"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestJWT(t *testing.T) {
	t.Run("create jwt token success should return token", func(t *testing.T) {
		b := Body{
			id:         uuid.New(),
			username:   "test",
			email:      "test@mail.com",
			created_at: util.MockTime(),
			expired_at: util.MockTime().Add(time.Minute * 2),
		}
		var s string

		j, err := NewJWTToken([]byte("secret"))

		assert.Nil(t, err)

		tk, err := j.CreateToken(b)
		assert.Nil(t, err)
		assert.IsType(t, s, tk)

	})

	t.Run("validate jwt token success should return true", func(t *testing.T) {

	})

}
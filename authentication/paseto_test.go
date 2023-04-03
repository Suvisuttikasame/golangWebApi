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
		token := "v2.local.dm0h0B6enYmXne8La6U1oHJ1jcY2Vysen_yNnMz3KTrIRtKb9nVhhBekaPqfKxUomytr2Hpz9ryHeQVL5edTskm7dO86Z98Vj-pJiJBGSf4nhPtt7VPw_X48tQM9i1geUfB_pHdLgUSkOhmI0CM6fHxkldtJqstvNwwekBe4WN6wrnaEv3DBptseJ9Wx2TAH1HFT7UgtYXGhhWHF-W4jy_CYHS-GER_r36uSF1k_Qrtj6mieij2Bdd6Xe74HKxCwT1S_xjvdhvdQv4wySNmvHt4G6hA5hVsbmrfTcWb5.bnVsbA"
		// token := "v2.local.6cOSDF1bDm1hixLW-Mcus8w0bk-vABdaLRRdtVy4ktYs7-fUetnLXQkm7xAFcWeC_INopwzIuewVn5yfQLawzOBRCW5gUQL6SMyYsMRVwaxEz2LSVCLpLw6CifWoW4z55zd3FnpMJw5zA4CXPuS8zTDmTtZQgExaK9GdbC0kW0BoUQ3VAPe_ssB6NukJJSoPqEfj-YMDc2GGl2qEG4QhZin8LWyD5Nxacl8qB3_1L_Gxku5FB4nH3knTYM-c9XhjgRirEVx_bjGVEeZMBoHiBXS6xnIEqLF-2kqykW_N.bnVsbA"

		p, err := NewPasetoToken([]byte("asdfgjhtuasdfgjhtuasdfgjhtuasdfg"))
		assert.Nil(t, err)

		r := p.Verification(token)

		t.Log(r)

	})

}

package authentication

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTAuthen struct {
	token     string
	secretKey []byte
}

type MyClaim struct {
	Id       uuid.UUID
	Username string
	Email    string
	jwt.RegisteredClaims
}

func NewJWTToken(k []byte) (Authen, error) {
	var j JWTAuthen
	j.secretKey = k
	return &j, nil
}

func (j *JWTAuthen) CreateToken(b Body) (string, error) {
	claim := MyClaim{
		b.Id,
		b.Username,
		b.Email,
		jwt.RegisteredClaims{
			Issuer:    b.Username,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	//jwt create token
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	tokenString, err := tok.SignedString(j.secretKey)
	if err != nil {
		return "", err
	}
	j.token = tokenString

	return tokenString, nil
}

func (j *JWTAuthen) Verification(t string) bool {
	var c MyClaim
	token, err := jwt.ParseWithClaims(t, &c, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return j.secretKey, nil
	})
	if err != nil {
		return false
	}

	if _, ok := token.Claims.(*MyClaim); ok && token.Valid {
		return true
	} else {
		return false
	}

}

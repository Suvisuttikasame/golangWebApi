package authentication

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JWTAuthen struct {
	token     string
	secretKey []byte
}

func NewJWTToken(k []byte) (Authen, error) {
	var j JWTAuthen
	j.secretKey = k
	return &j, nil
}

func (j *JWTAuthen) CreateToken(b Body) (string, error) {
	//jwt create token
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":         b.id,
		"username":   b.username,
		"email":      b.email,
		"created_at": b.created_at,
		"expired_at": b.expired_at,
	})

	tokenString, err := tok.SignedString(j.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *JWTAuthen) Verification(t string) bool {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return j.secretKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
		return true
	} else {
		fmt.Println(err)
		return false
	}

}

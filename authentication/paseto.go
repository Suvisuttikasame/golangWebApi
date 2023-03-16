package authentication

import (
	"errors"
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoAuthen struct {
	token        *paseto.V2
	symmetricKey []byte
}

type PasetoPayload paseto.JSONToken

func (pl *PasetoPayload) Valid() error {
	if time.Now().After(pl.Expiration) {
		return errors.New("token has been expired!")
	}
	return nil
}

func NewPasetoToken(key []byte) (Authen, error) {
	if len(key) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)
	}
	return &PasetoAuthen{
		token:        paseto.NewV2(),
		symmetricKey: key,
	}, nil

}

func (p *PasetoAuthen) CreateToken(b Body) (string, error) {
	payload := PasetoPayload{
		Issuer:     b.Username,
		IssuedAt:   time.Now(),
		Expiration: time.Now().Add(2 * time.Hour),
	}
	// payload.Set("myMessage", "hello paseto token")

	token, err := p.token.Encrypt(p.symmetricKey, payload, nil)

	if err != nil {
		return "", err
	}
	return token, err
}

func (p *PasetoAuthen) Verification(t string) bool {
	var payload PasetoPayload

	err := p.token.Decrypt(t, p.symmetricKey, &payload, nil)
	fmt.Println(payload)
	fmt.Println(err)
	if err != nil {
		return false
	}

	err = payload.Valid()
	fmt.Println(err)
	if err != nil {
		return false
	}

	return true

}

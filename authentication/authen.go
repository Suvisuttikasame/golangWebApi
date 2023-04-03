package authentication

import (
	"github.com/google/uuid"
)

type Body struct {
	Id       uuid.UUID
	Username string
	Email    string
}

type AuthenPaseto interface {
	CreateToken(b Body) (string, error)

	Verification(t string) (*PasetoPayload, error)
}

type Authen interface {
	CreateToken(b Body) (string, error)

	Verification(t string) (*MyClaim, error)
}

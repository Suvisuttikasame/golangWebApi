package authentication

import (
	"github.com/google/uuid"
)

type Body struct {
	Id       uuid.UUID
	Username string
	Email    string
}

type Authen interface {
	CreateToken(b Body) (string, error)

	Verification(t string) bool
}

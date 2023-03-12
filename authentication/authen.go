package authentication

import (
	"github.com/google/uuid"
)

type Body struct {
	id       uuid.UUID
	username string
	email    string
}

type Authen interface {
	CreateToken(b Body) (string, error)

	Verification(t string) bool
}

package authentication

import (
	"time"

	"github.com/google/uuid"
)

type Body struct {
	id         uuid.UUID
	username   string
	email      string
	created_at time.Time
	expired_at time.Time
}

type Authen interface {
	CreateToken(b Body) (string, error)

	Verification(t string) bool
}

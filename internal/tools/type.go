package tools

import (
	"time"
	"github.com/gofrs/uuid/v5"
)

type Account struct {
	Email     string
	Password  string
	FirstName string
	LastName  string

	TakenDSA bool

	Year   int
	Online bool
}

func NewAccount(email, password, firstName, lastName string, takenDSA bool, year int) *Account {
	return &Account{
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,

		TakenDSA: takenDSA,

		Year:   year,
		Online: true,
	}
}

type Session struct {
	Token    uuid.UUID
	Username string
	Expires  time.Time
}

func NewSession(token uuid.UUID, user string, expires time.Time) *Session{
	return &Session{
		Token: token,
		Username: user,
		Expires: expires,
	}
}

func (s Session) isExpired() bool {
	return s.Expires.Before(time.Now())
}
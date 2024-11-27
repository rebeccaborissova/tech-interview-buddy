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

	Description string
}

func NewAccount(email, password, firstName, lastName string, takenDSA bool, year int, des string) *Account {
	return &Account{
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,

		TakenDSA: takenDSA,

		Year:   year,
		Online: true,

		Description: des,
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

func (s Session) IsExpired() bool {
	return s.Expires.Before(time.Now())
}
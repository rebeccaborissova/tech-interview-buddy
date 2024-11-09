package tools

type Account struct {
	Email     string
	Password  string
	FirstName string
	LastName  string

	TakenDSA    bool

	Year   int
	Online bool
}

func NewAccount(email, password, firstName, lastName string, takenDSA bool, year int) *Account {
	return &Account{
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,

		TakenDSA:    takenDSA,

		Year:   year,
		Online: true,
	}
}

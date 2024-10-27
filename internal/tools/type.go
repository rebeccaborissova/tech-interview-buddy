package tools

type Account struct {
	Email     string
	Password  string
	FirstName string
	LastName  string

	TakenDSA    bool
	KnowsPython bool
	KnowsCPP    bool

	Year   int
	Online bool
}

func NewAccount(email, password, firstName, lastName string, takenDSA, knowsPython, knowsCPP bool, year int) *Account {
	return &Account{
		Email:     email,
		Password:  password,
		FirstName: firstName,
		LastName:  lastName,

		TakenDSA:    takenDSA,
		KnowsPython: knowsPython,
		KnowsCPP:    knowsCPP,

		Year:   year,
		Online: true,
	}
}

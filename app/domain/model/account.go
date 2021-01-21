package model

type Account struct {
	id   string
	name string
}

func NewAccount(id, name string) *Account {
	return &Account{
		id:   id,
		name: name,
	}
}

func (a *Account) GetId() string {
	return a.id
}

func (a *Account) GetName() string {
	return a.name
}

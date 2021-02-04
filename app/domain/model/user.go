package model

type User struct {
	id    string
	email string
	password string
}

func NewUser(id, email, password string) *User {
	return &User{
		id:    id,
		email: email,
		password: password,
	}
}

func (u *User) GetId() string {
	return u.id
}

func (u *User) GetEmail() string {
	return u.email
}

func (u *User) GetPassword() string {
	return u.password
}

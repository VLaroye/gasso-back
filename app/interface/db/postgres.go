package db

import (
	"github.com/VLaroye/gasso-back/app/domain/model"
	"gorm.io/gorm"
)

type User struct {
	ID    string
	Email string
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (ur *userRepository) FindAll() ([]*model.User, error) {
	var users []*User
	result := ur.db.Find(&users)

	if result.Error != nil {
		return nil, result.Error
	}

	response := make([]*model.User, len(users))

	for i, user := range users {
		response[i] = model.NewUser(user.ID, user.Email)
	}

	return response, nil
}

func (ur *userRepository) FindByEmail(email string) (*model.User, error) {
	var user User
	result := ur.db.Where("email = ?", email).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return model.NewUser(user.ID, user.Email), nil
}

func (ur *userRepository) Save(user *model.User) error {
	result := ur.db.Create(&User{
		ID:    user.GetId(),
		Email: user.GetEmail(),
	})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

type Account struct {
	Id   string
	Name string
}

type accountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *accountRepository {
	return &accountRepository{
		db: db,
	}
}

// TODO: Implement "real" account repository functions
func (ar *accountRepository) FindAll() ([]*model.Account, error) {
	return nil, nil
}

func (ar *accountRepository) FindByName(name string) (*model.Account, error) {
	return nil, nil
}

func (ar *accountRepository) FindById(id string) (*model.Account, error) {
	return nil, nil
}

func (ar *accountRepository) Create(account *model.Account) error {
	return nil
}

func (ar *accountRepository) Update(account *model.Account) error {
	return nil
}

func (ar *accountRepository) Delete(id string) error {
	return nil
}

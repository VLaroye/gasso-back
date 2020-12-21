package db

import (
	"github.com/VLaroye/gasso-back/app/domain/model"
	"gorm.io/gorm"
)

type User struct {
	ID string `gorm:`
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
	return nil, nil
}

func (ur *userRepository) Save(user *model.User) error {
	return nil
}

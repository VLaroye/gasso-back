package db

import (
	"github.com/VLaroye/gasso-back/app/domain/model"
	"gorm.io/gorm"
)

type User struct {
	ID    string
	Email string
	Password string
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
		response[i] = model.NewUser(user.ID, user.Email, user.Password)
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

	return model.NewUser(user.ID, user.Email, user.Password), nil
}

func (ur *userRepository) Save(user *model.User) error {
	result := ur.db.Create(&User{
		ID:    user.GetId(),
		Email: user.GetEmail(),
		Password: user.GetPassword(),
	})

	if result.Error != nil {
		return result.Error
	}

	return nil
}

package user

import (
	"github.com/irvanrifai/mkp-backend-development-test/api-ticketing/models"
	"gorm.io/gorm"
)

type IUserRepository interface {
	Create(user *models.User) error
	FindByEmail(email string) (models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return user, err
}

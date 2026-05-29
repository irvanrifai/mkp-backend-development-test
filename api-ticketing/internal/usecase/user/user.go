package user

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/irvanrifai/mkp-backend-development-test/api-ticketing/internal/repository/user"
	"github.com/irvanrifai/mkp-backend-development-test/api-ticketing/models"
	"golang.org/x/crypto/bcrypt"
)

type IUserUsecase interface {
	Register(name, username, email, password string) (models.User, error)
	Login(email, password string) (string, error)
}

type userUsecase struct {
	repo      user.IUserRepository
	jwtSecret string
}

func NewUserUsecase(repo user.IUserRepository, secret string) IUserUsecase {
	return &userUsecase{repo, secret}
}

func (u *userUsecase) Register(name, username, email, password string) (models.User, error) {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := models.User{Name: name, Username: username, Email: email, Password: string(hashed)}
	err := u.repo.Create(&user)
	return user, err
}

func (u *userUsecase) Login(email, password string) (string, error) {
	user, err := u.repo.FindByEmail(email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		return "", errors.New("Wrong credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	return token.SignedString([]byte(u.jwtSecret))
}

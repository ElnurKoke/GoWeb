package service

import (
	"fmt"
	"foruum/models"
	"foruum/repository"
	"log"
	"time"

	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

const salt = "Sfasfasfasfas"

var empty models.User

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) CreateUser(user models.User) error {
	if !CheckLogin(user.Login) {
		return fmt.Errorf("UNCORRECT LOGIN")
	}
	if !CheckPassword(user.Password) {
		return fmt.Errorf("UNCORRECT PASSWORD")
	}
	if !CheckUsername(user.Username) {
		return fmt.Errorf("Uncorrect username")
	}
	user.Password = s.GeneratePassword(user.Password)

	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(login, password string) (models.User, error) {
	user, err := s.repo.GetUser(login)
	if err != nil {
		return empty, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password+salt)); err != nil {
		return empty, err
	}
	token, err := uuid.NewV4()
	if err != nil {
		log.Print(err)
		return empty, err
	}

	user.Token = token.String()
	user.TokenDuration = time.Now().Add(12 * time.Hour)
	if err := s.repo.SaveTokens(user.Login, user.TokenDuration, user.Token); err != nil {
		return empty, nil
	}
	return user, nil
}

func (s *AuthService) GeneratePassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password+salt), bcrypt.DefaultCost)
	return string(bytes)
}

func (s *AuthService) GetUserByToken(token string) (models.User, error) {
	return s.repo.GetUserByToken(token)
}

func (s *AuthService) DeleteToken(token string) error {
	return s.repo.DeleteToken(token)
}

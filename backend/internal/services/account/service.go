package account

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"orderly/internal/di"
	"orderly/internal/domain"
	"orderly/pkg/jwt"
	"time"
)

type Service struct {
	Repository di.IAccountRepository
}
type ServiceDeps struct {
	Repository di.IAccountRepository
}

func NewService(deps *ServiceDeps) *Service {
	return &Service{
		Repository: deps.Repository,
	}
}

func (service *Service) Register(email, password, name string) (int, error) {
	existsUser := service.Repository.FindByEmail(email)
	if existsUser != nil {
		return -1, errors.New("the user already exists")
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return -1, errors.New("password hashing error")
	}
	user := &domain.User{
		Role:     "Customer",
		Password: string(hashPassword),
		Name:     name,
		Email:    email,
	}
	id, err := service.Repository.Create(user)
	if err != nil {
		return -1, nil
	}
	return id, nil
}

func (service *Service) Login(email, password string) (int, string, error) {
	user := service.Repository.FindByEmail(email)
	if user == nil {
		return -1, "", errors.New("invalid email or password")
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err == nil {
		return -1, "", errors.New("invalid email or password")
	}
	return user.Id, user.Role, nil
}

func (service *Service) IssueTokens(secret string, data jwt.Data) (string, string, error) {
	j := jwt.NewJWT(secret)
	accessToken, err := j.Create(data, time.Now().Add(time.Hour*2).Add(time.Minute*10))
	if err != nil {
		return "", "", err
	}
	refreshToken, err := j.Create(data, time.Now().AddDate(0, 0, 2).Add(time.Hour*2))
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
	return "", "", nil
}

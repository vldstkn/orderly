package account

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"orderly/internal/di"
	"orderly/internal/domain"
	api "orderly/internal/services/api/dto"
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
		Role:     string(domain.Customer),
		Password: string(hashPassword),
		Name:     name,
		Email:    email,
	}
	id, err := service.Repository.Create(user)

	return id, err
}

func (service *Service) Login(email, password string) (int, string, error) {
	user := service.Repository.FindByEmail(email)
	if user == nil {
		return -1, "", errors.New("invalid email or password")
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return -1, "", errors.New("invalid email or password")
	}
	return user.Id, user.Role, nil
}

func (service *Service) GetProfileById(id int) (*api.GetProfileResponse, error) {
	user := service.Repository.FindById(id)
	if user == nil {
		return nil, errors.New("the user does not exist")
	}
	return &api.GetProfileResponse{
		Id:        user.Id,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
	}, nil
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
}

func (service *Service) ChangeRole(id int, role string) error {
	isValid := isValidRole(role)
	if !isValid {
		return errors.New(http.StatusText(http.StatusBadRequest))
	}
	err := service.Repository.ChangeRoleById(id, role)
	return err
}

func isValidRole(role string) bool {
	switch domain.UserRole(role) {
	case domain.Customer, domain.Provider:
		return true
	default:
		return false
	}
}

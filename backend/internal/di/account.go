package di

import (
	"orderly/internal/domain"
	api "orderly/internal/services/api/dto"
	"orderly/pkg/jwt"
)

type IAccountService interface {
	Register(email, password, name string) (int, error)
	Login(email, password string) (int, string, error)
	IssueTokens(secret string, data jwt.Data) (string, string, error)
	GetProfileById(id int) (*api.GetProfileResponse, error)
	ChangeRole(id int, role string) error
}

type IAccountRepository interface {
	FindById(id int) *domain.User
	FindByEmail(email string) *domain.User
	Create(user *domain.User) (int, error)
	ChangeRoleById(id int, role string) error
}

package di

import (
	"orderly/internal/domain"
	"orderly/pkg/jwt"
)

type IAccountService interface {
	Register(email, password, name string) (int, error)
	Login(email, password string) (int, string, error)
	IssueTokens(secret string, data jwt.Data) (string, string, error)
}

type IAccountRepository interface {
	FindById(id int) *domain.User
	FindByEmail(email string) *domain.User
	Create(user *domain.User) (int, error)
}

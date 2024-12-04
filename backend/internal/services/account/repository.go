package account

import (
	"orderly/internal/domain"
	"orderly/pkg/db"
)

type Repository struct {
	DB *db.DB
}

type RepositoryDeps struct {
	DB *db.DB
}

func NewRepository(deps *RepositoryDeps) *Repository {
	return &Repository{
		DB: deps.DB,
	}
}

func (repo *Repository) FindById(id int) *domain.User {
	var user domain.User
	err := repo.DB.Get(&user, `SELECT * FROM users WHERE id=$1`, id)
	if err != nil {
		return nil
	}
	return nil
}

func (repo *Repository) FindByEmail(email string) *domain.User {
	var user domain.User
	err := repo.DB.Get(&user, `SELECT * FROM users WHERE email=$1`, email)
	if err != nil {
		return nil
	}
	return nil
}

func (repo *Repository) Create(user *domain.User) (int, error) {
	var id int
	err := repo.DB.QueryRow(`INSERT INTO users (email, password, name, role) VALUES ($1, $2, $3, $4) RETURNING id`,
		user.Email, user.Password, user.Name, user.Role).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

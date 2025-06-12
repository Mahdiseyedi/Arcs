package user

import (
	"arcs/internal/clients/db"
	"arcs/internal/models"
	"context"
)

type Repository struct {
	db *db.Database
}

func NewUserRepository(
	db *db.Database,
) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) CreateUser(ctx context.Context, id string, balance float64) error {
	return r.db.DB.WithContext(ctx).Create(&models.User{
		ID:      id,
		Balance: balance,
	}).Error
}

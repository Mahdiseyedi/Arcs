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

func (r *Repository) CreateUser(ctx context.Context, id string, balance int64) error {
	user := models.User{
		ID:      id,
		Balance: balance,
	}

	if err := r.db.DB.WithContext(ctx).Create(&user).Error; err != nil {
		return err
	}

	return nil
}

package order

import (
	"arcs/internal/clients/db"
	"arcs/internal/models"
	"context"
)

type Repository struct {
	db *db.Database
}

func NewOrderRepository(
	db *db.Database,
) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Submit(ctx context.Context, order models.Order) error {
	return r.db.DB.WithContext(ctx).Create(&order).Error
}

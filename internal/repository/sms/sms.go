package sms

import (
	"arcs/internal/clients/db"
	"arcs/internal/models"
	consts "arcs/internal/utils/const"
	"context"
)

type Repository struct {
	db *db.Database
}

func NewSMSRepository(
	db *db.Database,
) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) UpdateStatus(ctx context.Context, smsID, status string) error {
	return r.db.DB.WithContext(ctx).
		Model(&models.SMS{}).
		Where("id = ?", smsID).
		Update("status", status).Error
}

func (r *Repository) ListPending(ctx context.Context) ([]models.SMS, error) {
	var smss []models.SMS
	if err := r.db.DB.WithContext(ctx).
		Where("status = ?", consts.PendingStatus).
		Find(&smss).Error; err != nil {
		return nil, err
	}

	return smss, nil
}

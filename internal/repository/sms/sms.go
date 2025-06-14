package sms

import (
	"arcs/internal/clients/db"
	"arcs/internal/configs"
	"arcs/internal/models"
	consts "arcs/internal/utils/const"
	"context"
)

type Repository struct {
	cfg configs.Config
	db  *db.Database
}

func NewSMSRepository(
	cfg configs.Config,
	db *db.Database,
) *Repository {
	return &Repository{
		cfg: cfg,
		db:  db,
	}
}

func (r *Repository) Create(ctx context.Context, sms models.SMS) error {
	return r.db.DB.WithContext(ctx).Create(&sms).Error
}

func (r *Repository) CreateSMSBatch(ctx context.Context, smss []models.SMS) error {
	return r.db.DB.WithContext(ctx).CreateInBatches(&smss, r.cfg.Basic.SMSBatchSize).Error
}

func (r *Repository) Update(ctx context.Context, smss []models.SMS) error {
	if len(smss) == 0 {
		return nil
	}

	var ids []string
	for _, sms := range smss {
		ids = append(ids, sms.ID)
	}

	return r.db.DB.WithContext(ctx).
		Model(&models.SMS{}).
		Where("id IN ?", ids).
		Update("status", consts.PublishedStatus).Error
}

func (r *Repository) ListPending(ctx context.Context) ([]models.SMS, error) {
	var smss []models.SMS
	if err := r.db.DB.WithContext(ctx).
		Where("status = ?", consts.PendingStatus).
		Preload("Order").
		Find(&smss).Error; err != nil {
		return nil, err
	}
	
	return smss, nil
}

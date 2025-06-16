package sms

import (
	"arcs/internal/clients/db"
	"arcs/internal/configs"
	"arcs/internal/models"
	consts "arcs/internal/utils/const"
	"context"
	"time"
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

func (r *Repository) ListPending(ctx context.Context, createdAfter time.Time, batchSize int) ([]models.SMS, error) {
	var smss []models.SMS
	if err := r.db.DB.WithContext(ctx).
		Where("status = ? AND created_at > ?", consts.PendingStatus, createdAfter).
		Limit(batchSize).
		Preload("Order").
		Find(&smss).Error; err != nil {
		return nil, err
	}

	return smss, nil
}

func (r *Repository) GetUserSMS(ctx context.Context, userID string, filters models.SMSFilter) ([]models.SMS, int64, error) {
	var smsList []models.SMS
	var total int64

	q := r.db.DB.WithContext(ctx).
		Joins("JOIN orders ON orders.id = sms.order_id").
		Where("orders.user_id = ?", userID)

	if filters.Status != "" {
		q = q.Where("sms.status = ?", filters.Status)
	}
	if filters.StartDate != nil {
		q = q.Where("sms.created_at >= ?", filters.StartDate)
	}
	if filters.EndDate != nil {
		q = q.Where("sms.created_at <= ?", filters.EndDate)
	}

	//calculate count before pagination
	if err := q.Model(&models.SMS{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (filters.Page - 1) * filters.PageSize
	if err := q.Limit(filters.PageSize).
		Offset(offset).
		Order("sms.created_at DESC").
		Find(&smsList).Error; err != nil {
		return nil, 0, err
	}

	return smsList, total, nil
}

func (r *Repository) MarkDelivered(ctx context.Context, smsID string) error {
	return r.db.DB.WithContext(ctx).
		Model(&models.SMS{}).
		Where("id = ?", smsID).
		Update("status", consts.DeliveredStatus).Error
}

func (r *Repository) MarkFailed(ctx context.Context, smsID string) error {
	return r.db.DB.WithContext(ctx).
		Model(&models.SMS{}).
		Where("id = ?", smsID).
		Update("status", consts.FailedStatus).Error
}

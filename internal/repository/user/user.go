package user

import (
	"arcs/internal/clients/db"
	"arcs/internal/models"
	"arcs/internal/utils/errmsg"
	"context"
	"errors"
	"gorm.io/gorm"
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

func (r *Repository) GetUserBalance(ctx context.Context, uid string) (int64, error) {
	var user models.User
	err := r.db.DB.WithContext(ctx).
		Model(&models.User{}).
		Select("balance").
		Where("id = ?", uid).
		Take(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errmsg.UserNotFound
		}
		return 0, err
	}

	return user.Balance, err
}

func (r *Repository) DecreaseBalance(ctx context.Context, uid string, amount int64) error {
	res := r.db.DB.WithContext(ctx).Exec(
		`UPDATE users SET balance = balance - ? WHERE id = ? AND balance >= ?`,
		amount, uid, amount)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errmsg.UserNotFound
	}

	return nil
}

func (r *Repository) IncreaseBalance(ctx context.Context, uid string, amount int64) error {
	res := r.db.DB.WithContext(ctx).Exec(
		`UPDATE users SET balance = balance + ? WHERE id = ?`,
		amount, uid)
	if res.Error != nil {
		return res.Error
	}

	if res.RowsAffected == 0 {
		return errmsg.UserNotFound
	}

	return nil
}

package user

import (
	"arcs/internal/repository/balance"
	"arcs/internal/repository/user"
	"context"
	"github.com/google/uuid"
)

type Svc struct {
	userRepo    *user.Repository
	balanceRepo *balance.Repository
}

func NewUserSvc(
	userRepo *user.Repository,
	balance *balance.Repository,
) *Svc {
	return &Svc{
		userRepo:    userRepo,
		balanceRepo: balance,
	}
}

func (s *Svc) CreateUser(ctx context.Context, balance int64) error {
	uid := uuid.New().String()

	if _, err := s.userRepo.CreateUser(ctx, uid, balance); err != nil {
		return err
	}

	if err := s.balanceRepo.Set(ctx, uid, balance); err != nil {
		return err
	}

	return nil
}

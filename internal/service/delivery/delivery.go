package delivery

import (
	"arcs/internal/configs"
	"arcs/internal/models"
	"arcs/internal/repository/sms"
	"context"
	"log"
	"math/rand"
)

type Svc struct {
	cfg     configs.Config
	smsRepo *sms.Repository
}

func NewDeliveryService(
	cfg configs.Config,
	smsRepo *sms.Repository,
) *Svc {
	return &Svc{
		cfg:     cfg,
		smsRepo: smsRepo,
	}
}

func (s *Svc) SendSMS(ctx context.Context, sms models.SMS) error {
	suc := rand.Float32() < float32(s.cfg.Delivery.SuccessRate)/100

	if suc {
		log.Printf("[DELIVERY] SMS delivered: [%v]", sms)
		return s.smsRepo.MarkDelivered(ctx, sms.ID)
	} else {
		log.Printf("[DELIVERY] SMS failed: [%v]", sms)
		return s.smsRepo.MarkFailed(ctx, sms.ID)
	}
}

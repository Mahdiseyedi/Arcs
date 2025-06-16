package delivery

import (
	"arcs/internal/configs"
	"arcs/internal/models"
	"arcs/internal/service/buffer"
	consts "arcs/internal/utils/const"
	"context"
	"log"
	"math/rand"
)

type Svc struct {
	cfg     configs.Config
	flusher *buffer.StatusFlusher
}

func NewDeliveryService(
	cfg configs.Config,
	flusher *buffer.StatusFlusher,
) *Svc {
	return &Svc{
		cfg:     cfg,
		flusher: flusher,
	}
}

func (s *Svc) SendSMS(ctx context.Context, sms models.SMS) error {
	suc := rand.Float32() < float32(s.cfg.Delivery.SuccessRate)/100

	if suc {
		log.Printf("[DELIVERY] SMS delivered: [%v]-[%v]", sms.Destination, sms.Order.Content)
		s.flusher.Add(models.StatusUpdate{
			ID:     sms.ID,
			Status: consts.DeliveredStatus,
		})
	} else {
		log.Printf("[DELIVERY] SMS failed: [%v]-[%v]", sms.Destination, sms.Order.Content)
		s.flusher.Add(models.StatusUpdate{
			ID:     sms.ID,
			Status: consts.FailedStatus,
		})
	}

	return nil
}

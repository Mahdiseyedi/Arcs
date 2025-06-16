package mock

import (
	"arcs/internal/configs"
	"arcs/internal/models"
	"math/rand"
)

type MockDelivery struct {
	cfg configs.Config
}

func NewMockDelivery(cfg configs.Config) *MockDelivery {
	return &MockDelivery{cfg: cfg}
}

func (d *MockDelivery) SendSMS(sms models.SMS) bool {
	return rand.Float32() < float32(d.cfg.Delivery.SuccessRate)/100
}

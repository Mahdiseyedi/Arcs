package worker

import (
	"arcs/internal/models"
	"arcs/internal/service/delivery"
	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
)

type SMSHandler struct {
	deliverySvc *delivery.Svc
}

func NewSMSHandler(deliverySvc *delivery.Svc) *SMSHandler {
	return &SMSHandler{deliverySvc: deliverySvc}
}

func (s *SMSHandler) Handle(ctx context.Context) nats.MsgHandler {
	return func(msg *nats.Msg) {
		var sms models.SMS
		if err := json.Unmarshal(msg.Data, &sms); err != nil {
			log.Printf("falied to marshal: [%v]", err)
			return
		}

		if err := s.deliverySvc.SendSMS(ctx, sms); err != nil {
			log.Printf("delivery falied: [%v]", err)
			return
		}

		msg.Ack()
	}
}

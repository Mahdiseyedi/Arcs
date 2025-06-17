package worker

import (
	pb "arcs/internal/models/proto"
	"arcs/internal/service/delivery"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type SMSHandler struct {
	deliverySvc *delivery.Svc
}

func NewSMSHandler(deliverySvc *delivery.Svc) *SMSHandler {
	return &SMSHandler{deliverySvc: deliverySvc}
}

func (s *SMSHandler) Handle() nats.MsgHandler {
	return func(msg *nats.Msg) {
		var sms pb.SMS

		if err := proto.Unmarshal(msg.Data, &sms); err != nil {
			msg.Nak()
			return
		}
		
		if err := s.deliverySvc.SendSMS(&sms); err != nil {
			msg.Nak()
			return
		}

		msg.Ack()
	}
}

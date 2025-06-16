package configs

import (
	"arcs/internal/configs/basic"
	"arcs/internal/configs/db"
	"arcs/internal/configs/delivery"
	"arcs/internal/configs/nats"
	"arcs/internal/configs/order"
	"arcs/internal/configs/redis"
)

type Config struct {
	Basic    basic.Basic       `koanf:"basic"`
	DB       db.DB             `koanf:"database"`
	Redis    redis.Redis       `koanf:"redis"`
	Nats     nats.Nats         `koanf:"nats"`
	Order    order.Order       `koanf:"order"`
	Delivery delivery.Delivery `koanf:"delivery"`

	//add more section here
}

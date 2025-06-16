package configs

import (
	"arcs/internal/configs/basic"
	"arcs/internal/configs/consumer"
	"arcs/internal/configs/db"
	"arcs/internal/configs/delivery"
	"arcs/internal/configs/order"
	"arcs/internal/configs/producer"
	"arcs/internal/configs/redis"
)

type Config struct {
	Basic    basic.Basic       `koanf:"basic"`
	DB       db.DB             `koanf:"database"`
	Redis    redis.Redis       `koanf:"redis"`
	Producer producer.Producer `koanf:"producer"`
	Consumer consumer.Consumer `koanf:"consumer"`
	Order    order.Order       `koanf:"order"`
	Delivery delivery.Delivery `koanf:"delivery"`

	//add more section here
}

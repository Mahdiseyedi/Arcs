package configs

import (
	"arcs/internal/configs/basic"
	"arcs/internal/configs/db"
	"arcs/internal/configs/order"
	"arcs/internal/configs/redis"
)

type Config struct {
	Basic basic.Basic `koanf:"basic"`
	DB    db.DB       `koanf:"database"`
	Redis redis.Redis `koanf:"redis"`
	Order order.Order `koanf:"order"`

	//add more section here
}

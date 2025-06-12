package configs

import (
	"arcs/internal/configs/basic"
	"arcs/internal/configs/db"
	"arcs/internal/configs/redis"
)

type Config struct {
	Basic basic.Basic `koanf:"basic"`
	DB    db.DB       `koanf:"database"`
	Redis redis.Redis `koanf:"redis"`

	//add more section here
}

package redis

type Redis struct {
	Addr     string `koanf:"addr"`
	Password string `koanf:"password"`
	DB       int    `koanf:"db"`
}

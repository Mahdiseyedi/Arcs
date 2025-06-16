package delivery

type Delivery struct {
	SuccessRate         int `koanf:"success-rate"`
	BufferFlushInterval int `koanf:"buffer-flush-interval"`
}

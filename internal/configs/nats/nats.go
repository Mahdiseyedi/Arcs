package nats

type Nats struct {
	URL           string   `koanf:"url"`
	ClientTimeout int      `koanf:"client-timeout"`
	ReconnectWait int      `koanf:"reconnect-wait"`
	MaxReconnects int      `koanf:"max-reconnects"`
	RetryTimeOut  int      `koanf:"retry-timeout"`
	Stream        string   `koanf:"stream"`
	Queue         string   `koanf:"queue"`
	Subjects      []string `koanf:"subjects"`
}

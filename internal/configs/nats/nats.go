package nats

type Nats struct {
	URL           string   `koanf:"url"`
	ClientTimeout int      `koanf:"client-timeout"`
	RetryTimeOut  int      `koanf:"retry-timeout"`
	Stream        string   `koanf:"stream"`
	Queue         string   `koanf:"queue"`
	Subjects      []string `koanf:"subjects"`
}

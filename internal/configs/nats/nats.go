package nats

type Nats struct {
	URL      string   `koanf:"url"`
	Stream   string   `koanf:"stream"`
	Subjects []string `koanf:"subjects"`
}

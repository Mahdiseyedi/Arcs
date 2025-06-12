package nats

type Nats struct {
	URL      string   `koanf:"url"`
	Stream   string   `koanf:"stream"`
	Queue    string   `koanf:"queue"`
	Subjects []string `koanf:"subjects"`
}

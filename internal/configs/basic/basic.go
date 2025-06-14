package basic

type Basic struct {
	Port         string `koanf:"port"`
	Environment  string `koanf:"environment"`
	SMSBatchSize int    `koanf:"batch-size"`
}

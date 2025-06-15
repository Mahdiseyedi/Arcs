package basic

type Basic struct {
	Port                    string `koanf:"port"`
	Environment             string `koanf:"environment"`
	SMSBatchSize            int    `koanf:"insert-batch-size"`
	PendingProcessBatchSize int    `koanf:"pending-process-batch-size"`
}

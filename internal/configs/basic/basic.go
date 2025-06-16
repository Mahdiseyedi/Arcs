package basic

type Basic struct {
	TimeZone                string `koanf:"time-zone"`
	Port                    string `koanf:"port"`
	Environment             string `koanf:"environment"`
	SMSBatchSize            int    `koanf:"insert-batch-size"`
	PendingProcessBatchSize int    `koanf:"pending-process-batch-size"`
	RepublishLockDuration   int    `koanf:"republish-lock-duration"`
}

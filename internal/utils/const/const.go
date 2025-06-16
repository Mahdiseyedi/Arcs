package consts

// order status
const (
	PendingStatus   = "pending"
	PublishedStatus = "published"
	FailedStatus    = "failed"
	DeliveredStatus = "delivered"
)

// locks
const (
	RepublishLock = "lock:recover-unpublished-sms"
)

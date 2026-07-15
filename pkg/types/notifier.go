package types

// Notifier is the interface that all notification services have in common
type Notifier interface {
	StartNotification()
	SendNotification(Report)
	// Flush blocks until all queued notifications have been sent, without
	// closing the notifier (unlike Close). Safe to call between cycles.
	Flush()
	AddLogHook()
	GetNames() []string
	GetURLs() []string
	Close()
}

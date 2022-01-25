// package lognsq allows to send log messages to NSQ messaging services.
package lognsq

// LogNSQ provides methods for interacting with nsq services.
type LogNSQ interface {
	// Write sends data to nsqd service, and optionally, prints logs to STDERR.
	Write([]byte) (int, error)
	// Stop breaks connection to nsqd service. It has to run when LogNSQ is no
	// longer used.
	Stop()
}

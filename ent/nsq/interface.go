package nsq

// NSQ provides methods for interacting with NSQ API
type NSQ interface {
	Write([]byte) (int, error)
	Stop()
}

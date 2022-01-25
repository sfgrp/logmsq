package config

// Config provides data necessary to start NSQ session.
// NSQ is a distributed messaging service.
type Config struct {
	// PrintLogs enables printing of logs to STDERR.
	StderrLogs bool

	// Topic sets a namespace for messages.
	Topic string

	// NsqdURL is the URL to http service of an nsqd daemon.
	NSQdAddr string

	// Regex -- when set, only lines corresponding the pattern are printed out.
	Regex string
}

// Option type allows to send options to Config
type Option func(*Config)

// OptStderrLogs sets flag for printing logs to STDERR.
func OptStderrLogs(b bool) Option {
	return func(cfg *Config) {
		cfg.StderrLogs = b
	}
}

// OptTopic sets NSQ's topic to send logs to.
func OptTopic(s string) Option {
	return func(cfg *Config) {
		cfg.Topic = s
	}
}

// OptNSQdAddr
func OptNSQdAddr(s string) Option {
	return func(cfg *Config) {
		cfg.NSQdAddr = s
	}
}

package config

import (
	"log"
	"regexp"
)

// Config provides data necessary to start NSQ session.
// NSQ is a distributed messaging service.
type Config struct {
	// PrintLogs enables printing of logs to STDERR.
	StderrLogs bool

	// Topic sets a namespace for messages.
	Topic string

	// NsqdURL is the URL to http service of an nsqd daemon.
	Address string

	// Regex -- when set, only lines corresponding the pattern are printed out.
	Regex *regexp.Regexp
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

// OptAddress sets the address to a TCP nsqd service (e.g. 127.0.0.1:4150).
func OptAddress(s string) Option {
	return func(cfg *Config) {
		cfg.Address = s
	}
}

func OptRegex(s string) Option {
	regex, err := regexp.Compile(s)
	if err != nil {
		log.Fatal(err)
	}

	return func(cfg *Config) {
		cfg.Regex = regex
	}
}

// New creates a new Config according to given options.
func New(opts ...Option) Config {
	cfg := Config{}

	for i := range opts {
		opts[i](&cfg)
	}

	if cfg.Address == "" {
		log.Fatal("the address to a TCP nsqd service must be set (e.g. 127.0.0.1:4150)")
	}

	if cfg.Topic == "" {
		log.Fatal("the topic for nsqd service must be set (e.g. `mylogs`)")
	}

	return cfg
}

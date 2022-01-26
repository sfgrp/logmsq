package config

import (
	"log"
	"regexp"
)

// Config provides data necessary to start NSQ session.
// NSQ is a distributed messaging service.
type Config struct {
	// Topic sets a namespace for messages.
	Topic string

	// NsqdURL is the URL to http service of an nsqd daemon.
	Address string

	// Regex -- when set, only lines corresponding the pattern are kept for
	// nsqd server. The settings do not influence STDERR output. If both
	// Regex and Contains are given, thir effect is additive (both of them
	// must match.
	Regex *regexp.Regexp

	// Contains -- when set, lines that contain (or not contain) the pattern
	// will be kept. The setting does not influence StderrLogs output.
	// If a pattern should not match, it has to start with '!'.
	//
	// "api" -- log must contain "api" pattern.
	// "!api" -- log must NOT contain "api" pattern.
	Contains string

	// StderrLogs enables printing of logs to STDERR. This output is not
	// filtered and contains all log lines.
	StderrLogs bool

	// Sends results of filtering to STDOUT.
	Debug bool
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

// OptRegex sets a regular expression for filtering logs.
func OptRegex(s string) Option {
	regex, err := regexp.Compile(s)
	if err != nil {
		log.Fatal(err)
	}

	return func(cfg *Config) {
		cfg.Regex = regex
	}
}

// OptContains sets the Contains field.
func OptContains(s string) Option {
	return func(cfg *Config) {
		cfg.Contains = s
	}
}

// OptDebug sets Debug field
func OptDebug(b bool) Option {
	return func(cfg *Config) {
		cfg.Debug = b
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

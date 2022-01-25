package lognsq

import "github.com/sfgrp/lognsq/ent/nsq"

type lognsq struct {
	// NSQ provides implementations for interacting with nsq services.
	nsq.NSQ
}

// New creates new LogNSQ object for sending logs to nsq services.
func New(n nsq.NSQ) LogNSQ {
	return &lognsq{NSQ: n}
}

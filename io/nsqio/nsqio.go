// package nsqio provides redirection of messages to a nsqd TCP service.
// It is used, for example, for aggregating logs from several server instances.
package nsqio

import (
	"errors"
	"os"
	"regexp"

	nsq "github.com/nsqio/go-nsq"
	"github.com/sfgrp/lognsq/config"
)

// nsqio creates a "producer" to the nsqd service. The producer is able to
// publish messages to nsqd. It cannot consume the messages.
// If the same config is used for many different apps, all of them are able
// to contribute their logs to the same namespace using nsqd.
type nsqio struct {
	cfg   config.Config
	regex *regexp.Regexp
	*nsq.Producer
}

// New Creates a new nsqio instance. If creation of "producer" failed, it
// returns an error.
func New(cfg config.Config) (l *nsqio, err error) {
	var regex *regexp.Regexp
	var prod *nsq.Producer
	if cfg.Topic == "" {
		err = errors.New("config for nsqio cannot have an empty Topic field")
		return nil, err
	}

	if cfg.Regex != "" {
		regex, err = regexp.Compile(cfg.Regex)
		if err != nil {
			return nil, err
		}
	}

	nsqCfg := nsq.NewConfig()
	prod, err = nsq.NewProducer(cfg.NSQdAddr, nsqCfg)
	if err != nil {
		return nil, err
	}

	l = &nsqio{
		cfg:      cfg,
		regex:    regex,
		Producer: prod,
	}
	return l, err
}

// Write takes a slice of bytes and publishes it to STDERR as well as to
// nsqd service. It uses Topic given in the config.
func (l *nsqio) Write(bs []byte) (n int, err error) {
	if l.cfg.StderrLogs {
		n, err = os.Stderr.Write(bs)
	}
	if err == nil && l.regexOK(bs) {
		err = l.Publish(l.cfg.Topic, bs)
	}
	return n, err
}

func (n *nsqio) regexOK(bs []byte) bool {
	if n.regex == nil {
		return true
	}
	return n.regex.Match(bs)
}

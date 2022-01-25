// package nsqio provides redirection of messages to a nsqd TCP service.
// It is used, for example, for aggregating logs from several server instances.
package nsqio

import (
	"errors"
	"os"

	nsq "github.com/nsqio/go-nsq"
	"github.com/sfgrp/lognsq/config"
)

// nsqio creates a "producer" to the nsqd service. The producer is able to
// publish messages to nsqd. It cannot consume the messages.
// If the same config is used for many different apps, all of them are able
// to contribute their logs to the same namespace using nsqd.
type nsqio struct {
	cfg config.Config
	*nsq.Producer
}

// New Creates a new nsqio instance. If creation of "producer" failed, it
// returns an error.
func New(cfg config.Config) (n *nsqio, err error) {
	var prod *nsq.Producer
	if cfg.Topic == "" {
		err = errors.New("config for nsqio cannot have an empty Topic field")
		return nil, err
	}

	nsqCfg := nsq.NewConfig()
	prod, err = nsq.NewProducer(cfg.Address, nsqCfg)
	if err != nil {
		return nil, err
	}

	n = &nsqio{
		cfg:      cfg,
		Producer: prod,
	}
	return n, err
}

// Write takes a slice of bytes and publishes it to STDERR as well as to
// nsqd service. It uses Topic given in the config.
func (n *nsqio) Write(bs []byte) (num int, err error) {
	if !n.regexOK(bs) {
		return 0, nil
	}

	if n.cfg.StderrLogs {

		num, err = os.Stderr.Write(append(bs, byte('\n')))
	}
	if err == nil && n.regexOK(bs) {
		err = n.Publish(n.cfg.Topic, bs)
	}
	return num, err
}

func (n *nsqio) regexOK(bs []byte) bool {
	if n.cfg.Regex == nil {
		return true
	}
	return n.cfg.Regex.Match(bs)
}

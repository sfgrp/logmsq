// package nsqio provides redirection of messages to a nsqd TCP service.
// It is used, for example, for aggregating logs from several server instances.
package nsqio

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	gonsq "github.com/nsqio/go-nsq"
	"github.com/sfgrp/lognsq/config"
)

// nsqio creates a "producer" to the nsqd service. The producer is able to
// publish messages to nsqd. It cannot consume the messages.
// If the same config is used for many different apps, all of them are able
// to contribute their logs to the same namespace using nsqd.
type nsqio struct {
	cfg config.Config
	*gonsq.Producer
}

// New Creates a new nsqio instance. If creation of "producer" failed, it
// returns an error.
func New(cfg config.Config) (n *nsqio, err error) {
	var prod *gonsq.Producer
	if cfg.Topic == "" {
		err = errors.New("config for nsqio cannot have an empty Topic field")
		return nil, err
	}
	nsqCfg := gonsq.NewConfig()
	prod, err = gonsq.NewProducer(cfg.Address, nsqCfg)
	if cfg.Debug {
		prod.SetLoggerLevel(gonsq.LogLevelDebug)
	} else {
		nullLogger := log.New(ioutil.Discard, "", log.LstdFlags)
		prod.SetLogger(nullLogger, gonsq.LogLevelInfo)
	}

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
	if len(bs) == 0 {
		return 0, nil
	}

	if n.cfg.StderrLogs {
		num, err = os.Stderr.Write(append(bs, byte('\n')))
	}
	if err == nil && n.regexOK(bs) && n.containsOK(bs) {
		if n.cfg.Debug {
			fmt.Println(string(bs))
		}
		num = len(bs)
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

func (n *nsqio) containsOK(bs []byte) bool {
	if len(n.cfg.Contains) == 0 {
		return true
	}
	var negate bool
	pattern := []byte(n.cfg.Contains)
	if pattern[0] == '!' {
		negate = true
		if len(pattern) == 1 {
			return true
		}

		pattern = pattern[1:]

	}
	contains := bytes.Contains(bs, pattern)
	if negate {
		return !contains
	}

	return contains
}

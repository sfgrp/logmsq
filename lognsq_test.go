// package lognsq_test tests functionality of LogNSQ.
// To have tests running nsqld service has to be running on `localhost`.
// It is also useful to monitor the progress by running nsqadmin web-interface.
//
// in one terminal run
// nsql
// in another terminal run
// nsqadmin --nsqd-http-address='localhost:4151'
package lognsq_test

import (
	"reflect"
	"testing"

	"github.com/sfgrp/lognsq"
	"github.com/sfgrp/lognsq/config"
	"github.com/sfgrp/lognsq/io/nsqio"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	opts := []config.Option{
		config.OptTopic("test"),
		config.OptAddress("localhost:4150"),
	}
	cfg := config.New(opts...)
	n, err := nsqio.New(cfg)
	assert.Nil(t, err)

	l := lognsq.New(n)
	lType := reflect.TypeOf(l).String()
	assert.Equal(t, lType, "*lognsq.lognsq")
	l.Stop()
}

func TestWriteRegex(t *testing.T) {
	opts := []config.Option{
		config.OptTopic("test"),
		config.OptAddress("localhost:4150"),
		config.OptRegex(`^{"test":`),
	}
	cfg := config.New(opts...)
	n, err := nsqio.New(cfg)
	assert.Nil(t, err)
	l := lognsq.New(n)
	defer l.Stop()

	tests := []struct {
		msg, log string
		err      bool
		num      int
	}{
		{"nomatch", `test {"test":`, false, 0},
		{"match", `{"test": "value"}`, false, 17},
	}

	for _, v := range tests {
		num, err := l.Write([]byte(v.log))
		assert.Equal(t, v.num, num)
		assert.Equal(t, v.err, err != nil)
	}
}

func TestWriteContains(t *testing.T) {
	opts := []config.Option{
		config.OptTopic("test"),
		config.OptAddress("localhost:4150"),
		config.OptContains("!api"),
	}
	cfg := config.New(opts...)
	n, err := nsqio.New(cfg)
	assert.Nil(t, err)
	l := lognsq.New(n)
	defer l.Stop()

	tests := []struct {
		msg, log string
		err      bool
		num      int
	}{
		{"nomatch", `test {"test":`, false, 13},
		{"match", `{"test": "api"}`, false, 0},
	}

	for _, v := range tests {
		num, err := l.Write([]byte(v.log))
		assert.Equal(t, v.num, num)
		assert.Equal(t, v.err, err != nil)
	}
}

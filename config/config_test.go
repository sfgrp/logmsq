package config_test

import (
	"regexp"
	"testing"

	"github.com/sfgrp/lognsq/config"
	"github.com/stretchr/testify/assert"
)

var minConfig = config.Config{
	Topic:   "test",
	Address: "localhost:4150",
}

var maxConfig = config.Config{
	Topic:      "test",
	Address:    "localhost:4150",
	Regex:      regexp.MustCompile(`http:\/\/`),
	Contains:   "!api",
	StderrLogs: true,
	Debug:      true,
}

func TestMin(t *testing.T) {
	opts := []config.Option{
		config.OptTopic("test"),
		config.OptAddress("localhost:4150"),
	}
	cfg := config.New(opts...)
	assert.Equal(t, cfg, minConfig)
}

func TestCustom(t *testing.T) {
	opts := []config.Option{
		config.OptTopic("test"),
		config.OptAddress("localhost:4150"),
		config.OptRegex(`http:\/\/`),
		config.OptContains("!api"),
		config.OptStderrLogs(true),
		config.OptDebug(true),
	}
	cfg := config.New(opts...)
	assert.Equal(t, cfg, maxConfig)
}

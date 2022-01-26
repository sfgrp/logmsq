/*
Copyright © 2022 Dmitry Mozzherin <dmozzherin@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"bufio"
	"log"
	"os"

	"github.com/sfgrp/lognsq"
	"github.com/sfgrp/lognsq/config"
	"github.com/sfgrp/lognsq/io/nsqio"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

// cfgData purpose is to achieve automatic import of data from the
// environment variables, if they are given.
type cfgData struct {
	Topic      string
	Address    string
	Regex      string
	Contains   string
	StderrLogs bool
	Debug      bool
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lognsq",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, _ []string) {
		cfgd := getConf()
		topicFlag(cmd, cfgd)
		addressFlag(cmd, cfgd)
		regexFlag(cmd, cfgd)
		containsFlag(cmd, cfgd)
		printFlag(cmd, cfgd)
		debugFlag(cmd, cfgd)
		opts := getOpts(cfgd)
		cfg := config.New(opts...)
		n, err := nsqio.New(cfg)
		if err != nil {
			log.Fatal(err)
		}

		l := lognsq.New(n)
		defer l.Stop()

		processStdin(l)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringP("topic", "t", "", "a topic to send to nsqd service (required).")
	rootCmd.Flags().StringP("nsqd-tcp-address", "a", "", "the address of an nsqd service (`127.0.0.1:4150`).")
	rootCmd.Flags().StringP("regex-filter", "r", "", "rejects all log messages that do not match the regex.")
	rootCmd.Flags().StringP("contains-filter", "c", "", "log must contai (or not) a pattern (`api`, `!api`).")
	rootCmd.Flags().BoolP("print-log", "p", false, "print logs to STDERR as well.")
	rootCmd.Flags().BoolP("debug", "d", false, "prints filtered logs to STDOUT.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	// Set environment variables for configuration
	_ = viper.BindEnv("StderrLogs", "LOGNSQ_STDERR_LOGS")
	_ = viper.BindEnv("Topic", "LOG_NSQ_TOPIC")
	_ = viper.BindEnv("NSQdAddr", "LOG_NSQ_ADDR")
	_ = viper.BindEnv("Regex", "LOG_NSQ_REGEX")
	_ = viper.BindEnv("Contains", "LOG_NSQ_CONTAINS")

	viper.AutomaticEnv() // read in environment variables that match
}

func getConf() *cfgData {
	cfg := &cfgData{}
	err := viper.Unmarshal(cfg)

	if err != nil {
		log.Fatalf("Cannot deserialize config data: %s.", err)
	}

	return cfg
}

// getOpts imports data from the configuration file. Some of the settings can
// be overriden by command line flags.
func getOpts(cfg *cfgData) []config.Option {
	var opts []config.Option
	if cfg.Regex != "" {
		opts = append(opts, config.OptRegex(cfg.Regex))
	}
	if cfg.Topic != "" {
		opts = append(opts, config.OptTopic(cfg.Topic))
	}
	if cfg.Address != "" {
		opts = append(opts, config.OptAddress(cfg.Address))
	}
	if cfg.Contains != "" {
		opts = append(opts, config.OptContains(cfg.Contains))
	}
	if cfg.StderrLogs {
		opts = append(opts, config.OptStderrLogs(true))
	}
	if cfg.Debug {
		opts = append(opts, config.OptDebug(true))
	}
	return opts
}

func processStdin(l lognsq.LogNSQ) {
	defer l.Stop()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		l.Write(scanner.Bytes())
	}

	if scanner.Err() != nil {
		log.Fatal(scanner.Err())
	}
}

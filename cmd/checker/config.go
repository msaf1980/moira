package main

import (
	"github.com/moira-alert/moira"
	"github.com/moira-alert/moira/checker"
	"github.com/moira-alert/moira/cmd"
	"github.com/xiam/to"
)

type config struct {
	Redis     cmd.RedisConfig     `yaml:"redis"`
	Logger    cmd.LoggerConfig    `yaml:"log"`
	Checker   checkerConfig       `yaml:"checker"`
	Telemetry cmd.TelemetryConfig `yaml:"telemetry"`
	Remote    cmd.RemoteConfig    `yaml:"remote"`
}

type triggerLogConfig struct {
	ID    string `yaml:"id"`
	Level string `yaml:"level"`
}

type triggersLogConfig struct {
	TriggersToLevel []triggerLogConfig `yaml:"triggers"`
}

type checkerConfig struct {
	// Period for every trigger to perform forced check on
	NoDataCheckInterval string `yaml:"nodata_check_interval"`
	// Period for every trigger to cancel forced check (earlier than 'NoDataCheckInterval') if no metrics were received
	StopCheckingInterval string `yaml:"stop_checking_interval"`
	// Min period to perform triggers re-check. Note: Reducing of this value leads to increasing of CPU and memory usage values
	CheckInterval string `yaml:"check_interval"`
	// Max period to perform lazy triggers re-check. Note: lazy triggers are triggers which has no subscription for it. Moira will check its state less frequently.
	// Delay for check lazy trigger is random between LazyTriggersCheckInterval/2 and LazyTriggersCheckInterval.
	LazyTriggersCheckInterval string `yaml:"lazy_triggers_check_interval"`
	// Max concurrent checkers to run. Equals to the number of processor cores found on Moira host by default or when variable is defined as 0.
	MaxParallelChecks int `yaml:"max_parallel_checks"`
	// Max concurrent remote checkers to run. Equals to the number of processor cores found on Moira host by default or when variable is defined as 0.
	MaxParallelRemoteChecks int `yaml:"max_parallel_remote_checks"`
	// Specify log level by entities
	SetLogLevel triggersLogConfig `yaml:"set_log_level"`
}

func (config *checkerConfig) getSettings(logger moira.Logger) *checker.Config {
	logTriggersToLevel := make(map[string]string)
	for _, v := range config.SetLogLevel.TriggersToLevel {
		logTriggersToLevel[v.ID] = v.Level
	}
	logger.Infof("Found dynamic log rules in config for %d triggers", len(logTriggersToLevel))

	return &checker.Config{
		CheckInterval:               to.Duration(config.CheckInterval),
		LazyTriggersCheckInterval:   to.Duration(config.LazyTriggersCheckInterval),
		NoDataCheckInterval:         to.Duration(config.NoDataCheckInterval),
		StopCheckingIntervalSeconds: int64(to.Duration(config.StopCheckingInterval).Seconds()),
		MaxParallelChecks:           config.MaxParallelChecks,
		MaxParallelRemoteChecks:     config.MaxParallelRemoteChecks,
		LogTriggersToLevel:          logTriggersToLevel,
	}
}

func getDefault() config {
	return config{
		Redis: cmd.RedisConfig{
			Addrs:       "localhost:6379",
			MetricsTTL:  "1h",
			DialTimeout: "500ms",
		},
		Logger: cmd.LoggerConfig{
			LogFile:         "stdout",
			LogLevel:        "info",
			LogPrettyFormat: false,
		},
		Checker: checkerConfig{
			NoDataCheckInterval:       "60s",
			CheckInterval:             "5s",
			LazyTriggersCheckInterval: "10m",
			StopCheckingInterval:      "30s",
			MaxParallelChecks:         0,
			MaxParallelRemoteChecks:   0,
		},
		Telemetry: cmd.TelemetryConfig{
			Listen: ":8092",
			Graphite: cmd.GraphiteConfig{
				Enabled:      false,
				RuntimeStats: false,
				URI:          "localhost:2003",
				Prefix:       "DevOps.Moira",
				Interval:     "60s",
			},
			Pprof: cmd.ProfilerConfig{Enabled: false},
		},
		Remote: cmd.RemoteConfig{
			CheckInterval: "60s",
			Timeout:       "60s",
			MetricsTTL:    "7d",
		},
	}
}

package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/moira-alert/moira"
	"github.com/moira-alert/moira/api/handler"
	"github.com/moira-alert/moira/cmd"
	"github.com/moira-alert/moira/database/redis"
	"github.com/moira-alert/moira/index"
	logging "github.com/moira-alert/moira/logging/zerolog_adapter"
	metricSource "github.com/moira-alert/moira/metric_source"
	"github.com/moira-alert/moira/metric_source/local"
	"github.com/moira-alert/moira/metric_source/prometheus"
	"github.com/moira-alert/moira/metric_source/remote"
	_ "go.uber.org/automaxprocs"
)

const serviceName = "api"

var (
	configFileName         = flag.String("config", "/etc/moira/api.yml", "Path to configuration file")
	printVersion           = flag.Bool("version", false, "Print version and exit")
	printDefaultConfigFlag = flag.Bool("default-config", false, "Print default config and exit")
)

// Moira api bin version
var (
	MoiraVersion = "unknown"
	GitCommit    = "unknown"
	GoVersion    = "unknown"
)

func main() {
	flag.Parse()
	if *printVersion {
		fmt.Println("Moira Api")
		fmt.Println("Version:", MoiraVersion)
		fmt.Println("Git Commit:", GitCommit)
		fmt.Println("Go Version:", GoVersion)
		os.Exit(0)
	}

	applicationConfig := getDefault()
	if *printDefaultConfigFlag {
		cmd.PrintConfig(applicationConfig)
		os.Exit(0)
	}

	err := cmd.ReadConfig(*configFileName, &applicationConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can not read settings: %s\n", err.Error())
		os.Exit(1)
	}

	apiConfig := applicationConfig.API.getSettings(applicationConfig.Redis.MetricsTTL, applicationConfig.Remote.MetricsTTL)

	logger, err := logging.ConfigureLog(applicationConfig.Logger.LogFile, applicationConfig.Logger.LogLevel, serviceName, applicationConfig.Logger.LogPrettyFormat)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Can not configure log: %s\n", err.Error())
		os.Exit(1)
	}
	defer logger.Info().
		String("moira_version", MoiraVersion).
		Msg("Moira API stopped")

	telemetry, err := cmd.ConfigureTelemetry(logger, applicationConfig.Telemetry, serviceName)
	if err != nil {
		logger.Fatal().
			Error(err).
			Msg("Can not start telemetry")
	}
	defer telemetry.Stop()

	databaseSettings := applicationConfig.Redis.GetSettings()
	notificationHistorySettings := applicationConfig.NotificationHistory.GetSettings()
	database := redis.NewDatabase(logger, databaseSettings, notificationHistorySettings, redis.API)

	// Start Index right before HTTP listener. Fail if index cannot start
	searchIndex := index.NewSearchIndex(logger, database, telemetry.Metrics)
	if searchIndex == nil {
		logger.Fatal().Msg("Failed to create search index")
	}

	err = searchIndex.Start()
	if err != nil {
		logger.Fatal().
			Error(err).
			Msg("Failed to start search index")
	}
	defer searchIndex.Stop() //nolint

	stats := newTriggerStats(logger, database, telemetry.Metrics)
	stats.Start()
	defer stats.Stop() //nolint

	if !searchIndex.IsReady() {
		logger.Fatal().Msg("Search index is not ready, exit")
	}

	// Start listener only after index is ready
	listener, err := net.Listen("tcp", apiConfig.Listen)
	if err != nil {
		logger.Fatal().
			Error(err).
			Msg("Failed to start listening")
	}

	logger.Info().
		String("listen_address", apiConfig.Listen).
		Msg("Start listening")

	remoteConfig := applicationConfig.Remote.GetRemoteSourceSettings()
	prometheusConfig := applicationConfig.Prometheus.GetPrometheusSourceSettings()

	localSource := local.Create(database)
	remoteSource := remote.Create(remoteConfig)
	prometheusSource, err := prometheus.Create(prometheusConfig, logger)
	if err != nil {
		logger.Fatal().
			Error(err).
			Msg("Failed to initialize prometheus metric source")
	}

	metricSourceProvider := metricSource.CreateMetricSourceProvider(
		localSource,
		remoteSource,
		prometheusSource,
	)

	webConfigContent, err := applicationConfig.Web.getSettings(remoteConfig.Enabled || prometheusConfig.Enabled)
	if err != nil {
		logger.Fatal().
			Error(err).
			Msg("Failed to get web applicationConfig content ")
	}

	httpHandler := handler.NewHandler(database, logger, searchIndex, apiConfig, metricSourceProvider, webConfigContent)
	server := &http.Server{
		Handler: httpHandler,
	}

	go func() {
		server.Serve(listener) //nolint
	}()
	defer Stop(logger, server)

	logger.Info().
		String("moira_version", MoiraVersion).
		Msg("Moira Api Started")
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)

	signal := fmt.Sprint(<-ch)
	logger.Info().
		String("signal", signal).
		Msg("Moira API shutting down.")
}

// Stop Moira API HTTP server
func Stop(logger moira.Logger, server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) //nolint
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error().
			Error(err).
			Msg("Can't stop Moira API correctly")
	}
}

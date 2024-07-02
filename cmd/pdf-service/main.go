package main

import (
	"fmt"
	"os"

	"github.com/hop-/goconfig"
	"github.com/hop-/golog"
	"github.com/hop-/pdf-service/internal/app"
	"github.com/hop-/pdf-service/internal/kafka"
)

func getHttpOptions() (bool, *uint16, bool, *string, *string) {
	enabled, err := goconfig.Get[bool]("http.enabled")
	if err != nil {
		golog.Fatalf("Failed to get configuration %s", err.Error())
	}

	port, err := goconfig.Get[uint16]("http.port")
	if err != nil {
		golog.Fatalf("Failed to get configuration %s", err.Error())
	}
	secure, err := goconfig.Get[bool]("http.secure.enabled")
	if err != nil {
		golog.Fatalf("Failed to get configuration %s", err.Error())
	}
	key, err := goconfig.Get[string]("http.secure.keyFile")
	if err != nil {
		golog.Fatalf("Failed to get configuration %s", err.Error())
	}
	cert, err := goconfig.Get[string]("http.secure.certFile")
	if err != nil {
		golog.Fatalf("Failed to get configuration %s", err.Error())
	}

	return *enabled, port, *secure, key, cert
}

func getKafkaOptions() (bool, *string, *string, bool, *string, *string) {
	enabled, err := goconfig.Get[bool]("kafka.enabled")
	if err != nil {
		golog.Fatalf("Failed to get configuration %s", err.Error())
	}

	kafkaHost, err := goconfig.Get[string]("kafka.host")
	if err != nil {
		golog.Fatalf("Failed to get configuration %s", err.Error())
	}
	kafkaConsumerGroupId, err := goconfig.Get[string]("kafka.group.id")
	if err != nil {
		golog.Fatalf("Failed to get configuration %s", err.Error())
	}
	createConsumerTopics, err := goconfig.Get[bool]("kafka.createConsumerTopics")
	if err != nil {
		golog.Fatalf("Failed to get configuration %s", err.Error())
	}
	requestsTopic, err := goconfig.Get[string]("kafka.topic.requests")
	if err != nil {
		golog.Fatalf("Failed to get configuration %s", err.Error())
	}
	responsesTopic, err := goconfig.Get[string]("kafka.topic.responses")
	if err != nil {
		golog.Fatalf("Failed to get configuration %s", err.Error())
	}

	return *enabled, kafkaHost, kafkaConsumerGroupId, *createConsumerTopics, requestsTopic, responsesTopic
}

func createKafkaConsumerTopics(host string, topic string) {
	kafkaOptions := kafka.UtilsOptions{
		Host: host,
	}

	kafkaUtils, err := kafka.NewUtils(&kafkaOptions)
	if err != nil {
		golog.Fatalf("Faild to connect to the kafka %s", err.Error())
	}
	err = kafkaUtils.CreateTopics([]string{topic}, 50) // partition number is hardcoded
	if err != nil {
		golog.Warningf("Failed to create kafka topics %s", err.Error())
	}
}

func main() {
	// Load config
	if err := goconfig.Load(); err != nil {
		fmt.Printf("Failed to load configs %s\n", err.Error())
		os.Exit(1)
	}

	logMode, err := goconfig.Get[string]("log.mode")
	if err != nil {
		mode := "INFO"
		fmt.Printf("Failed to get log mode default is %s\n", mode)
		logMode = &mode
	}
	// Init Logging
	golog.Init(*logMode)

	// Get configs
	name, err := goconfig.Get[string]("name")
	if err != nil {
		golog.Fatalf("Failed to get configuration %s", err.Error())
	}
	golog.Infof("Starting %s", *name)

	engineType, err := goconfig.Get[string]("engine")
	if err != nil {
		golog.Fatalf("Failed to get configuration %s", err.Error())
	}
	concurrency, err := goconfig.Get[uint8]("concurrency")
	if err != nil {
		golog.Fatalf("Failed to get configuration %s", err.Error())
	}

	// Create app options
	opts := []app.OptionModifier{
		app.WithEngine(*engineType),
		app.WithConcurrency(*concurrency),
	}

	httpIsEnabled,
		httpPort,
		httpsIsEnabled,
		keyFile,
		certFile := getHttpOptions()

	kafkaIsEnabled,
		kafkaHost,
		kafkaConsumerGroupId,
		createConsumerTopics,
		requestsTopic,
		responsesTopic := getKafkaOptions()

	if !kafkaIsEnabled && !httpIsEnabled {
		golog.Fatalf("At least one of the services should be enabled")
	}

	// Add kafka if enabled
	if kafkaIsEnabled {
		opts = append(opts, app.WithKafka(
			*kafkaHost,
			*kafkaConsumerGroupId,
			*requestsTopic,
			*responsesTopic,
		))
	}

	// Add http if enabled
	if httpIsEnabled {
		opts = append(opts, app.WithHttp(
			*httpPort,
			httpsIsEnabled,
			*certFile,
			*keyFile,
		))
	}

	// Create App
	app := app.New(opts...)

	// Create consumer topics if preferred
	if createConsumerTopics {
		createKafkaConsumerTopics(*kafkaHost, *requestsTopic)
	}

	// Run app
	app.Start()
}

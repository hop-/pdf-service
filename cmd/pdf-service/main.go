package main

import (
	"fmt"
	"os"

	"github.com/hop-/goconfig"
	"github.com/hop-/golog"
	"github.com/hop-/pdf-service/internal/app"
	"github.com/hop-/pdf-service/internal/kafka"
)

func getKafkaOptions() (bool, *string, *string, bool, *string, *string) {
	enabled, err := goconfig.Get[bool]("kafka.enabled")
	if err != nil {
		golog.Fatalf("Failed to get configuration %s", err.Error())
	}

	if !*enabled {
		return false, nil, nil, false, nil, nil
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

	return true, kafkaHost, kafkaConsumerGroupId, *createConsumerTopics, requestsTopic, responsesTopic
}

func main() {
	// Load config
	if err := goconfig.Load(); err != nil {
		fmt.Printf("Failed to load configs %s\n", err.Error())
		os.Exit(1)
	}

	logMode, err := goconfig.Get[string]("log.mode")
	if err != nil {
		mode := "WARNING"
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

	kafkaIsEnabled,
		kafkaHost,
		kafkaConsumerGroupId,
		createConsumerTopics,
		requestsTopic,
		responsesTopic := getKafkaOptions()

	// Create App
	app := app.NewApp(app.Options{
		EngineType:     *engineType,
		Concurrency:    *concurrency,
		KafkaIsEnabled: kafkaIsEnabled,
		Kafka: app.KafkaOptions{
			Host:            *kafkaHost,
			ConsumerGroupId: *kafkaConsumerGroupId,
			RequestsTopic:   *requestsTopic,
			ResponsesTopic:  *responsesTopic,
		},
	})

	if createConsumerTopics {
		kafkaOptions := kafka.UtilsOptions{
			Host: *kafkaHost,
		}

		kafkaUtils, err := kafka.NewUtils(&kafkaOptions)
		if err != nil {
			golog.Fatalf("Faild to connect to the kafka %s", err.Error())
		}
		err = kafkaUtils.CreateTopics([]string{*requestsTopic}, 50) // partition number is hardcoded
		if err != nil {
			golog.Warningf("Failed to create kafka topics %s", err.Error())
		}
	}

	// Run app
	app.Start()
}

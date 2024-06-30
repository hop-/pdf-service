package kafka

import (
	"context"

	confluentkafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type UtilsOptions struct {
	Host string
}

type Utils struct {
	admin *confluentkafka.AdminClient
}

func NewUtils(options *UtilsOptions) (*Utils, error) {
	config := confluentkafka.ConfigMap{
		"bootstrap.servers": options.Host,
	}

	u, err := confluentkafka.NewAdminClient(&config)
	if err != nil {
		return nil, err
	}

	return &Utils{u}, nil
}

func (u *Utils) CreateTopics(topics []string, partitionNumber int) error {
	var topicSpecifications []confluentkafka.TopicSpecification
	for _, topic := range topics {
		topicSpecifications = append(topicSpecifications, confluentkafka.TopicSpecification{
			Topic:             topic,
			NumPartitions:     partitionNumber,
			ReplicationFactor: 1,
		})
	}

	// Creating topics in background
	_, err := u.admin.CreateTopics(context.Background(), topicSpecifications, nil)

	return err
}

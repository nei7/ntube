package kafka_service

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaConfig struct {
	Host  string `mapstructure:"KAFKA_HOST"`
	Topic string `mapstructure:"KAFKA_TOPIC"`
}

func NewKafkaProducer(conf KafkaConfig) (*kafka.Producer, error) {
	config := kafka.ConfigMap{
		"bootstrap.servers": conf.Host,
	}

	client, err := kafka.NewProducer(&config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewKafkaConsumer(conf KafkaConfig, groupID string) (*kafka.Consumer, error) {
	config := kafka.ConfigMap{
		"bootstrap.servers":  conf.Host,
		"group.id":           groupID,
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	}

	client, err := kafka.NewConsumer(&config)
	if err != nil {
		return nil, err
	}

	if err := client.Subscribe(conf.Topic, nil); err != nil {
		return nil, err
	}

	return client, nil
}

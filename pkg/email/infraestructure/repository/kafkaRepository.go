package repository

import (
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jeffleon/email-service/internal/config"
	"github.com/jeffleon/email-service/pkg/email/domain"
	"github.com/sirupsen/logrus"
)

type kafkaRepository struct {
	consumer *kafka.Consumer
}

func NewKafkaRepository(consumer *kafka.Consumer) domain.KafkaRepository {
	return &kafkaRepository{
		consumer,
	}
}

func (k *kafkaRepository) TopicConsume() error {
	err := k.consumer.SubscribeTopics([]string{config.EnvConfigs.KafkaUserTopic}, nil)
	logrus.Infof("kafka consumer listen topic %s", config.EnvConfigs.KafkaUserTopic)
	if err != nil {
		return err
	}

	// A signal handler or similar could be used to set this to false to break the loop.
	go func() {
		run := true
		for run {
			msg, err := k.consumer.ReadMessage(time.Second)
			if err == nil {
				fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
			} else if !err.(kafka.Error).IsTimeout() {
				// The client will automatically try to recover from all errors.
				// Timeout is not considered an error because it is raised by
				// ReadMessage in absence of messages.
				fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			}
		}
		k.consumer.Close()
	}()

	return nil
}

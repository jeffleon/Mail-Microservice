package repository

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/jeffleon/email-service/internal/config"
	"github.com/jeffleon/email-service/pkg/email/domain"
	"github.com/sirupsen/logrus"
)

type kafkaRepository struct {
	consumer *kafka.Consumer
	sendMail domain.SendMailer
}

func NewKafkaRepository(consumer *kafka.Consumer, sendMail domain.SendMailer) domain.KafkaRepository {
	return &kafkaRepository{
		consumer,
		sendMail,
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
				fmt.Printf("Message on %s: %s %s\n", msg.TopicPartition, string(msg.Key), string(msg.Value))
				err := k.Listen(msg)
				if err != nil {
					logrus.Errorf("Consumer error sending email: %v (%v)\n", err, msg)
				}
			} else if !err.(kafka.Error).IsTimeout() {
				fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			}
		}
		k.consumer.Close()
	}()

	return nil
}

func (k *kafkaRepository) Listen(msg *kafka.Message) error {

	if string(msg.Key) == "created_user" {
		var user domain.User
		err := json.Unmarshal(msg.Value, &user)
		if err != nil {
			return err
		}

		err = k.sendMail(domain.Message{
			From:     config.EnvConfigs.EmailFromAddress,
			FromName: config.EnvConfigs.EmailFromName,
			To:       user.Email,
			Subject:  fmt.Sprintf("Welcome from kafka %s", user.FirstName),
			Message:  fmt.Sprintf("Welcome from kafka %s", user.FirstName),
		})
		if err != nil {
			return err
		}
		logrus.Infof("Sending from Kafka succedded to %s", user.Email)
		return nil
	}
	return nil
}

package domain

type KafkaRepository interface {
	TopicConsume() error
}

package di

import (
	"strings"

	"github.com/segmentio/kafka-go"
)

func GetKafkaReader(url string, topic string) *kafka.Reader {
	brokers := strings.Split(url, ",")

	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}

func GetKafkaWriter(kafkaURL string, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(kafkaURL),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

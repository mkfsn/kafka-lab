package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/mkfsn/kafka-lab/internal/di"
	"github.com/segmentio/kafka-go"
	"golang.org/x/exp/slog"
)

func main() {
	kafkaURL := "localhost:9092"
	topic := "my-events"

	logger := slog.New(slog.NewTextHandler(os.Stdout))

	logger.Info("program start",
		slog.String("kafkaURL", kafkaURL),
		slog.String("topic", topic),
	)

	writer := di.GetKafkaWriter(kafkaURL, topic)
	defer writer.Close()

	reader := di.GetKafkaReader(kafkaURL, topic)
	defer reader.Close()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func(ctx context.Context, logger *slog.Logger) {
		logger.Info("start consuming ... !!")

		for {
			m, err := reader.ReadMessage(ctx)
			if err != nil {
				logger.Error("read message", err)
				break
			}

			logger.Info("read message",
				slog.String("topic", m.Topic),
				slog.Int("partition", m.Partition),
				slog.Int64("offset", m.Offset),
				slog.String("value", string(m.Value)),
			)
		}
	}(ctx, logger.WithGroup("consumer"))

	go func(ctx context.Context, logger *slog.Logger) {
		logger.Info("start producing ... !!")

		for i := 0; ; i++ {
			key := fmt.Sprintf("Key-%d", i)
			msg := kafka.Message{
				Key:   []byte(key),
				Value: []byte(fmt.Sprint(uuid.New())),
			}

			err := writer.WriteMessages(ctx, msg)
			if err != nil {
				logger.Error("write messages", err)
			} else {
				logger.Info("produced", slog.String("key", key))
			}

			time.Sleep(1 * time.Second)
		}
	}(ctx, logger.WithGroup("producer"))

	logger.Info("press Ctrl+C to stop program")
	<-ctx.Done()
	logger.Info("graceful shutdown")
}

package kafka

import (
	"time"

	kaf "github.com/segmentio/kafka-go"
)

type KafkaGoClient struct {
	Consumer *kaf.Reader
	Producer *kaf.Writer
}

func NewGoClient(
	openConsumer bool,
	openProducer bool,
	brokers []string,
	consumerTopic string,
	consumerGroup string,
	producerTopic string,
	producerGroup string,
	readOldest bool,
) *KafkaGoClient {
	goClient := new(KafkaGoClient)

	if openConsumer {
		goClient.Consumer = kaf.NewReader(kaf.ReaderConfig{
			Brokers:  brokers,
			Topic:    consumerTopic,
			GroupID:  consumerGroup,
			MaxBytes: 10e6, // 10MB,
			MaxWait:  60 * 2 * time.Second,
		})
		if readOldest {
			goClient.Consumer.SetOffset(kaf.FirstOffset) // EARLIEST = -2
		}
	}

	if openProducer {
		goClient.Producer = kaf.NewWriter(kaf.WriterConfig{
			Brokers: brokers,
			Topic:   producerTopic,
		})
	}

	return goClient
}

func (g *KafkaGoClient) Close() {
	if g.Consumer != nil {
		g.Consumer.Close()
	}
	if g.Producer != nil {
		g.Producer.Close()
	}
}

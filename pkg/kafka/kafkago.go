package kafka

import (
	"log"
	"strings"
	"time"

	kaf "github.com/segmentio/kafka-go"
)

// KafkaGoClient struct.
type KafkaClient struct {
	Consumer *kaf.Reader
	Producer *kaf.Writer
}

// NewKafkaClient creates a new KafkaClient struct.
func NewKafkaClient(
	read bool,
	write bool,
	brokers []string,
	consumerTopic string,
	consumerGroup string,
	producerTopic string,
	producerGroup string,
	readOldest bool,
) *KafkaClient {
	client := new(KafkaClient)

	if read {
		client.Consumer = NewConsumer(brokers, consumerTopic, consumerGroup, readOldest)
		if readOldest {
			_ = client.Consumer.SetOffset(kaf.FirstOffset) // EARLIEST = -2
		}
	}

	if write {
		client.Producer = NewProducer(brokers, producerTopic)
	}

	return client
}

// NewConsumer creates a new kafka consumer.
func NewConsumer(brokers []string, topic string, group string, oldest bool) *kaf.Reader {
	return kaf.NewReader(kaf.ReaderConfig{
		Brokers:     brokers,
		GroupTopics: strings.Split(topic, ","),
		GroupID:     group,
		MinBytes:    1,
		MaxBytes:    1e5,
		MaxWait:     3 * time.Second,
	})
}

// NewProducer creates a new kafka producer.
func NewProducer(brokers []string, topic string) *kaf.Writer {
	return &kaf.Writer{
		Addr: kaf.TCP(brokers...),
		// Topic:                  topic,
		AllowAutoTopicCreation: true,
		BatchSize:              1,
		BatchTimeout:           2 * time.Second,
		Balancer:               &kaf.LeastBytes{},
		RequiredAcks:           1,
	}
}

// Close closes active comnusmers/producers.
func (client *KafkaClient) Close() {
	if client.Consumer != nil {
		if err := client.Consumer.Close(); err != nil {
			log.Fatal(err)
		}
	}
	if client.Producer != nil {
		if err := client.Producer.Close(); err != nil {
			log.Fatal(err)
		}
	}
}

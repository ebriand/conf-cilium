package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/Shopify/sarama"
)

var (
	brokers  string
	producer sarama.SyncProducer
	version  = "2.1.0"
	group    = "2"
)

func newProducer() sarama.SyncProducer {

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 10
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(strings.Split(brokers, ","), config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	return producer
}

type Consumer func(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error

func (consumer Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (consumer Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	return consumer(session, claim)
}

func syncEntityFromKafka(topic string, consumer *Consumer) {
	config := sarama.NewConfig()
	config.Version, _ = sarama.ParseKafkaVersion(version)
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	ctx := context.Background()
	client, err := sarama.NewConsumerGroup(strings.Split(brokers, ","), group, config)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			err := client.Consume(ctx, []string{topic}, consumer)
			if err != nil {
				panic(err)
			}
		}
	}()
}

func syncEventsFromKafka() {
	topic := "events"
	eventConsumer := Consumer(func(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
		for message := range claim.Messages() {
			log.Printf("New message: %v\n", message)
			var e Event
			err := json.Unmarshal(message.Value, &e)
			if err != nil {
				log.Printf("Unable to parse msg: %v\n", message.Value)
			} else {
				log.Printf("Adding event: %v\n", e)
				events = append(events, e)
			}
		}
		return nil
	})
	syncEntityFromKafka(topic, &eventConsumer)
}

func addEvent(e Event) {
	eJSON, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: "events",
		Value: sarama.StringEncoder(eJSON),
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		log.Printf("%s written\n", eJSON)
	}
}

func init() {
	brokers = os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		panic("env variable KAFKA_BROKERS is missing")
	}
}

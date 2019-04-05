package main

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/ebriand/conf-cilium/types"
)

var (
	brokers  string
	producer sarama.SyncProducer
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

func addHero(h types.Hero, i types.Identity) {
	hJSON, err := json.Marshal(h)
	if err != nil {
		panic(err)
	}
	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: "heroes",
		Value: sarama.StringEncoder(hJSON),
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		log.Printf("%s written\n", hJSON)
	}

	iJSON, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: "identities",
		Value: sarama.StringEncoder(iJSON),
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
	} else {
		log.Printf("%s written\n", iJSON)
	}
}

func init() {
	brokers = os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		panic("env variable KAFKA_BROKERS is missing")
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Shopify/sarama"
	"github.com/ebriand/conf-cilium/types"
	"github.com/google/uuid"
)

var heroes = []types.Hero{
	{"batman", uuid.Must(uuid.Parse("38bdf3f3-3a4d-4786-86bc-91f20860d804"))},
	{"superman", uuid.Must(uuid.Parse("ae9ac595-da8f-4089-8ba1-c8bcb5dd6b01"))},
}

var identities = []types.Identity{
	{uuid.Must(uuid.Parse("38bdf3f3-3a4d-4786-86bc-91f20860d804")), "Bruce Wayne"},
	{uuid.Must(uuid.Parse("38bdf3f3-3a4d-4786-86bc-91f20860d804")), "Kalel"},
}

func newDataCollector(brokerList []string) sarama.SyncProducer {

	// For the data collector, we are looking for strong consistency semantics.
	// Because we don't change the flush settings, sarama will try to produce messages
	// as fast as possible to keep latency low.
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true

	// On the broker side, you may want to change the following settings to get
	// stronger consistency guarantees:
	// - For your broker, set `unclean.leader.election.enable` to false
	// - For the topic, you could increase `min.insync.replicas`.

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	return producer
}

func createIdentities(producer sarama.SyncProducer) {
	for _, h := range identities {
		jsonMsg, err := json.Marshal(h)
		if err != nil {
			panic(err)
		}
		_, _, err = producer.SendMessage(&sarama.ProducerMessage{
			Topic: "identities",
			Value: sarama.StringEncoder(jsonMsg),
		})
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("%s written\n", jsonMsg)
		}
	}
}

func createHeroes(producer sarama.SyncProducer) {
	for _, h := range heroes {
		jsonMsg, err := json.Marshal(h)
		if err != nil {
			panic(err)
		}
		_, _, err = producer.SendMessage(&sarama.ProducerMessage{
			Topic: "heroes",
			Value: sarama.StringEncoder(jsonMsg),
		})
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("%s written\n", jsonMsg)
		}
	}
}

func main() {
	producer := newDataCollector([]string{"localhost:9092"})

	fmt.Println("Simple kafka producer")
	fmt.Println("---------------------")

	createHeroes(producer)
	createIdentities(producer)
}

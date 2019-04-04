package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/ebriand/conf-cilium/types"
	"github.com/gorilla/mux"
)

var heroes []types.Hero

var identities []types.Identity

var (
	version = "2.1.0"
	brokers = "localhost:9092"
	group   = "1"
	topics  = "heroes"
)

func heroesToNames(heroes []types.Hero) []string {
	var names = []string{}
	for _, h := range heroes {
		names = append(names, h.Name)
	}
	return names
}

func getHeroes() []Hero {
	return heroes
}

func getHeroByName(name string) (*Hero, error) {
	for _, h := range getHeroes() {
		if name == h.Name {
			return &h, nil
		}
	}
	return nil, fmt.Errorf("hero %s not found", name)
}

func getIdentityByID(id int) (*types.Identity, error) {
	for _, i := range identities {
		if id == i.ID {
			return &i, nil
		}
	}
	return nil, fmt.Errorf("identity %d not found", id)
}

func HeroesHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(heroesToNames(getHeroes()))
}

func HeroHandler(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	hero, err := getHeroByName(name)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(hero)
}

func IdentitiesHandler(w http.ResponseWriter, r *http.Request) {
	idString := mux.Vars(r)["identityID"]
	id, err := strconv.Atoi(idString)
	if err != nil {
		fmt.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	identity, err := getIdentityByID(id)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(identity)
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func ReadyHandler(w http.ResponseWriter, r *http.Request) {
	//TODO: check kafka connection
	w.WriteHeader(http.StatusNoContent)
}

// Consumer represents a Sarama consumer group consumer
type Consumer struct{}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (consumer *Consumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (consumer *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (consumer *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		fmt.Printf("New message: %v\n", message)
		var h Hero
		err := json.Unmarshal(message.Value, &h)
		if err != nil {
			fmt.Printf("Unable to parse msg: %v\n", message.Value)
		} else {
			fmt.Printf("Adding hero: %v", h)
			heroes = append(heroes, h)
		}
	}

	return nil
}

func syncHeroesFromKafka() {

	topic := "heroes"
	config := sarama.NewConfig()
	config.Version, _ = sarama.ParseKafkaVersion(version)
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer := Consumer{}

	ctx := context.Background()
	client, err := sarama.NewConsumerGroup(strings.Split(brokers, ","), group, config)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			err := client.Consume(ctx, []string{topic}, &consumer)
			if err != nil {
				panic(err)
			}
		}
	}()
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/heroes", HeroesHandler).Methods("GET")
	r.HandleFunc("/heroes/{name}", HeroHandler).Methods("GET")
	r.HandleFunc("/identities/{identityID:[0-9]+}", IdentitiesHandler).Methods("GET")
	r.HandleFunc("/health", HealthHandler).Methods("GET")
	r.HandleFunc("/ready", ReadyHandler).Methods("GET")
	http.Handle("/", r)

	go syncHeroesFromKafka()

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		panic(err)
	}
}

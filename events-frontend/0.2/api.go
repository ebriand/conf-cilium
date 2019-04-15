package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/ebriand/conf-cilium/types"
	"github.com/google/uuid"
)

var apiURL string

func api(resource string) string {
	return apiURL + "/" + resource
}

func getHeroesName() ([]string, error) {
	url := api("heroes")
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var heroes []string
	err = json.Unmarshal(body, &heroes)
	if err != nil {
		return nil, err
	}
	return heroes, nil
}

func getHero(name string) (*types.Hero, error) {
	url := api("heroes/" + name)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var hero types.Hero
	json.Unmarshal(body, &hero)
	return &hero, nil
}

func getIdentity(id uuid.UUID) (*types.Identity, error) {
	var identity types.Identity
	for _, i := range getIdentitiesFromKafkaSync() {
		if i.ID == id {
			identity = i
			break
		}
	}
	if &identity == nil {
		identity = types.Identity{}
	}
	return &identity, nil
}

func init() {
	apiURL = os.Getenv("API_URL")
	if apiURL == "" {
		panic("env variable API_URL is missing")
	}
}

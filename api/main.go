package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ebriand/conf-cilium/types"
	"github.com/gorilla/mux"
)

var heroes = []types.Hero{
	{Name: "batman", SecretIdentityID: 1},
	{Name: "superman", SecretIdentityID: 2},
}

var identities = []types.Identity{
	{ID: 1, RealName: "Bruce Wayne"},
	{ID: 2, RealName: "Kalel"},
}

func heroesToNames(heroes []types.Hero) []string {
	var names = []string{}
	for _, h := range heroes {
		names = append(names, h.Name)
	}
	return names
}

func getHeroByName(name string) (*types.Hero, error) {
	for _, h := range heroes {
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
	json.NewEncoder(w).Encode(heroesToNames(heroes))
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

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/heroes", HeroesHandler).Methods("GET")
	r.HandleFunc("/heroes/{name}", HeroHandler).Methods("GET")
	r.HandleFunc("/identities/{identityID:[0-9]+}", IdentitiesHandler).Methods("GET")
	r.HandleFunc("/health", HealthHandler).Methods("GET")
	r.HandleFunc("/ready", ReadyHandler).Methods("GET")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":80", nil))
}

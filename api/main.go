package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Hero struct {
	Name             string `json:"name"`
	SecretIdentityID int    `json:"secretIdentityID"`
}

type Identity struct {
	ID       int    `json:"id"`
	RealName string `json:"realName"`
}

var heroes = []Hero{
	{"batman", 1},
	{"superman", 2},
}

var identities = []Identity{
	{1, "Bruce Wayne"},
	{2, "Kalel"},
}

func heroesToNames(heroes []Hero) []string {
	var names = []string{}
	for _, h := range heroes {
		names = append(names, h.Name)
	}
	return names
}

func getHeroByName(name string) (*Hero, error) {
	for _, h := range heroes {
		if name == h.Name {
			return &h, nil
		}
	}
	return nil, fmt.Errorf("hero %s not found", name)
}

func getIdentityByID(id int) (*Identity, error) {
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
	http.ListenAndServe(":80", nil)
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/ebriand/conf-cilium/types"
	"github.com/gorilla/mux"
)

func api(resource string) (string, error) {
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		return "", fmt.Errorf("API_URL env variable is missing.")
	}
	return apiURL + "/" + resource, nil
}

func HeroesHandler(w http.ResponseWriter, r *http.Request) {
	url, err := api("heroes")
	if err != nil {
		fmt.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp, err := http.Get(url)
	if err != nil {
		fmt.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	var heroes []types.Hero
	json.Unmarshal(body, &heroes)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(heroes)
}

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func ReadyHandler(w http.ResponseWriter, r *http.Request) {
	url, err := api("heroes")
	if err != nil {
		fmt.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = http.Get(url)
	if err != nil {
		fmt.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/heroes", HeroesHandler).Methods("GET")
	//r.HandleFunc("/heroes/{name}", HeroHandler).Methods("GET")
	r.HandleFunc("/health", HealthHandler).Methods("GET")
	r.HandleFunc("/ready", ReadyHandler).Methods("GET")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":80", nil))
}

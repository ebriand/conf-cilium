package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/ebriand/conf-cilium/types"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

var (
	events     []Event
	identities []types.Identity

	indexTemplate  *template.Template
	detailTemplate *template.Template
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	heroesName, err := getHeroesName()
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	indexTemplate.Execute(w, indexData{events, heroesName})
}

func addEventHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	e := Event{ID: uuid.New(), Name: r.FormValue("name"), Heroes: r.Form["heroes"]}
	addEvent(e)

	w.Header().Add("location", "/")
	w.WriteHeader(http.StatusFound)
}

func getEvent(eventIDString string) (*Event, error) {
	eventID, err := uuid.Parse(eventIDString)
	if err != nil {
		return nil, err
	}

	for _, e := range events {
		if e.ID == eventID {
			return &e, nil
		}
	}
	return nil, fmt.Errorf("event %s not found", eventIDString)
}

func eventDetailHandler(w http.ResponseWriter, r *http.Request) {
	eventIDString := mux.Vars(r)["eventID"]

	event, err := getEvent(eventIDString)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var heroes []Hero

	for _, name := range event.Heroes {
		h, err := getHero(name)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		i, err := getIdentity(h.SecretIdentityID)
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		heroes = append(heroes, Hero{h.Name, i.RealName})
	}
	w.WriteHeader(http.StatusOK)
	detailTemplate.Execute(w, detailData{event.Name, heroes})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	_, err := getHeroesName()
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func init() {
	indexTemplate = template.Must(template.ParseFiles("templates/index.html"))
	detailTemplate = template.Must(template.ParseFiles("templates/detail.html"))
}

func main() {
	producer = newProducer()

	syncEventsFromKafka()

	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/events", addEventHandler).Methods("POST")
	r.HandleFunc("/events/{eventID}", eventDetailHandler).Methods("GET")
	r.HandleFunc("/health", healthHandler).Methods("GET")
	r.HandleFunc("/ready", readyHandler).Methods("GET")
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":80", nil))
}

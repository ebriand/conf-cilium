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
	indexTemplate  *template.Template
	heroesTemplate *template.Template
	detailTemplate *template.Template
)

type heroesData struct {
	Heroes []types.Hero
}

type detailData struct {
	Hero     *types.Hero
	Identity *types.Identity
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	indexTemplate.Execute(w, nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("location", "/heroes")
	w.WriteHeader(http.StatusFound)
}

func heroesHandler(w http.ResponseWriter, r *http.Request) {
	heroes, err := getHeroesName()
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)

	heroesTemplate.Execute(w, heroesData{heroes})
}

func heroAddHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var h types.Hero
	var i types.Identity
	name := r.FormValue("name")

	i = types.Identity{ID: uuid.New(), RealName: r.FormValue("identity")}
	h = types.Hero{Name: name, SecretIdentityID: i.ID}
	imgURL, err := getHeroImage(name)
	if err != nil {
		log.Printf("Image not found for %s. Err is: %v", name, err)
		h.ImageURL = "/assets/pow-1601674_1280.png"
	} else {
		h.ImageURL = imgURL
	}

	addHero(h, i)

	w.Header().Add("location", "/heroes")
	w.WriteHeader(http.StatusFound)
}

func heroDetailHandler(w http.ResponseWriter, r *http.Request) {
	name := mux.Vars(r)["name"]
	hero, err := getHero(name)
	if err != nil {
		fmt.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	identity, err := getIdentity(hero.SecretIdentityID)

	w.WriteHeader(http.StatusOK)
	detailTemplate.Execute(w, detailData{hero, identity})
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	_, err := getHeroesName()
	if err != nil {
		fmt.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func init() {
	indexTemplate = template.Must(template.ParseFiles("templates/index.html"))
	heroesTemplate = template.Must(template.ParseFiles("templates/heroes.html"))
	detailTemplate = template.Must(template.ParseFiles("templates/detail.html"))
}

func main() {
	producer = newProducer()

	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/login", loginHandler).Methods("POST")
	r.HandleFunc("/heroes", heroesHandler).Methods("GET")
	r.HandleFunc("/heroes", heroAddHandler).Methods("POST")
	r.HandleFunc("/heroes/{name}", heroDetailHandler).Methods("GET")
	r.HandleFunc("/health", healthHandler).Methods("GET")
	r.HandleFunc("/ready", readyHandler).Methods("GET")
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(":80", nil))
}

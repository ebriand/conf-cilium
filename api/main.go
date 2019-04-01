package main

import (
  "fmt"
  "net/http"
  "github.com/gorilla/mux"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    fmt.Fprint(w, "Hello")
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/", HomeHandler)
    http.Handle("/", r)
    http.ListenAndServe(":80", nil)
}

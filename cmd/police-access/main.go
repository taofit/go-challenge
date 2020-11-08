package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/taofit/go-challenge/internal"
)

func handleOfficers(router *mux.Router) {
	router.HandleFunc("/officers", internal.GetOfficers).Methods("GET")
	router.HandleFunc("/officers/{id}", internal.GetOfficer).Methods("GET")
	router.HandleFunc("/officers/{id}", internal.UpdateOfficer).Methods("PUT")
	router.HandleFunc("/officers", internal.CreateOfficer).Methods("POST")
	router.HandleFunc("/officers/{id}", internal.DeleteOfficer).Methods("DELETE")
}

func handleBikeThefts(router *mux.Router) {
	router.HandleFunc("/bike-thefts", internal.CreateCase).Methods("POST")
	router.HandleFunc("/bike-thefts", internal.GetCases).Methods("GET")
	router.HandleFunc("/bike-thefts/{id}", internal.GetCase).Methods("GET")
	router.HandleFunc("/bike-thefts/{id}", internal.UpdateCase).Methods("PUT")
}

func main() {
	fmt.Println("see the light")
	router := mux.NewRouter()
	handleOfficers(router)
	handleBikeThefts(router)
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err.Error())
	}
}

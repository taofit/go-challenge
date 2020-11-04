package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/taofit/go-challenge/internal"
)

func main() {
	fmt.Println("see the light")
	router := mux.NewRouter()

	router.HandleFunc("/officers", internal.GetOfficers).Methods("GET")
	router.HandleFunc("/officers/{id}", internal.GetOfficer).Methods("GET")
	router.HandleFunc("/officers/{id}", internal.UpdateOfficer).Methods("PUT")
	router.HandleFunc("/officers", internal.CreateOfficer).Methods("POST")
	router.HandleFunc("/officers/{id}", internal.DeleteOfficer).Methods("DELETE")
	http.ListenAndServe(":8080", router)
}

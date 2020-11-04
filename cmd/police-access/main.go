package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Officer struct {
	ID   string `json:"id"`
	NAME string `json:"name"`
}

var officers []Officer

func getOfficers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(officers)
}

func getOfficer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, item := range officers {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			break
		}
		return
	}
	json.NewEncoder(w).Encode(&Officer{})
}

func updateOfficer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	fmt.Print(params)
	for index, item := range officers {
		if item.ID == params["id"] {
			officers = append(officers[:index], officers[index+1:]...)

			var officer Officer
			_ = json.NewDecoder(r.Body).Decode(&officer)
			officer.ID = params["id"]
			officers = append(officers, officer)
			json.NewEncoder(w).Encode(&officer)

			return
		}
	}

	json.NewEncoder(w).Encode(officers)
}

func createOfficer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var officer Officer
	_ = json.NewDecoder(r.Body).Decode(&officer)

	officer.ID = strconv.Itoa(rand.Intn(1000000))
	officers = append(officers, officer)
	json.NewEncoder(w).Encode(&officer)
}

func deleteOfficer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, item := range officers {
		if item.ID == params["id"] {
			officers = append(officers[:index], officers[index+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(&officers)
}

func main() {
	fmt.Println("see the light")
	router := mux.NewRouter()
	officers = append(officers, Officer{ID: "1", NAME: "My first post"})
	officers = append(officers, Officer{ID: "2", NAME: "more officer will see the post coming about the stolen bikes in malmo"})

	router.HandleFunc("/officers", getOfficers).Methods("GET")
	router.HandleFunc("/officers/{id}", getOfficer).Methods("GET")
	router.HandleFunc("/officers/{id}", updateOfficer).Methods("PUT")
	router.HandleFunc("/officers", createOfficer).Methods("POST")
	router.HandleFunc("/officers/{id}", deleteOfficer).Methods("DELETE")
	http.ListenAndServe(":8080", router)
}

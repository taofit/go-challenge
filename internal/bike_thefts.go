package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type TheftCase struct {
	ID          int       `json:"id"`
	TITLE       string    `json:"title"`
	BRAND       string    `json:"brand"`
	CITY        string    `json:"city"`
	DESCRIPTION string    `json:"description"`
	REPORTED    time.Time `json:"reported"`
	UPDATED     time.Time `json:"updated"`
	SOLVED      bool      `json:"solved"`
	OFFICER     Officer
	// image string ``
}

func CreateCase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var theftCase TheftCase
	_ = json.NewDecoder(r.Body).Decode(&theftCase)

	if theftCase.TITLE == "" || theftCase.BRAND == "" || theftCase.CITY == "" || theftCase.DESCRIPTION == "" {
		respondWithJSON(w, http.StatusBadRequest, "Some fields are missing please enter then again")
		return
	}

	db := dbConn()
	insert, err := db.Prepare("INSERT INTO bike_thefts(title, brand, city, description) VALUES(?,?,?,?)")
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
	}

	_, err = insert.Exec(theftCase.TITLE, theftCase.BRAND, theftCase.CITY, theftCase.DESCRIPTION)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
	}
	fmt.Println("Bike theft '" + theftCase.TITLE + "' is created")
	defer db.Close()
}

func GetCases(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()
	selResult, err := db.Query(`SELECT bt.id, bt.title, bt.brand, bt.city, bt.description, bt.reported, bt.updated, bt.solved, IFNULL(o.id, 0), IFNULL(o.name, '')
								FROM bike_thefts bt
								LEFT JOIN officers o
								ON o.id = bt.officer
								ORDER BY bt.id DESC`)
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
	}

	theftCases := []TheftCase{}
	theftCase := TheftCase{}
	for selResult.Next() {
		var officerId int
		var officerName string
		err = selResult.Scan(&theftCase.ID, &theftCase.TITLE, &theftCase.BRAND,
			&theftCase.CITY, &theftCase.DESCRIPTION, &theftCase.REPORTED, &theftCase.UPDATED, &theftCase.SOLVED, &officerId, &officerName)
		if err != nil {
			respondWithJSON(w, http.StatusBadRequest, err.Error())
		}
		theftCase.OFFICER.ID = officerId
		theftCase.OFFICER.NAME = officerName
		theftCases = append(theftCases, theftCase)
	}

	defer db.Close()
	json.NewEncoder(w).Encode(&theftCases)
}

func GetCase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	db := dbConn()
	var officerId int
	var officerName string
	selResult, err := db.Query(`SELECT bt.id, bt.title, bt.brand, bt.city, bt.description, bt.reported, bt.updated, bt.solved, IFNULL(o.id, 0), IFNULL(o.name, '')
								FROM bike_thefts bt
								LEFT JOIN officers o
								ON o.id = bt.officer
								WHERE bt.id=?`, id)
	if err != nil {
		panic(err.Error())
	}

	theftCase := TheftCase{}
	for selResult.Next() {
		err = selResult.Scan(&theftCase.ID, &theftCase.TITLE, &theftCase.BRAND,
			&theftCase.CITY, &theftCase.DESCRIPTION, &theftCase.REPORTED, &theftCase.UPDATED, &theftCase.SOLVED, &officerId, &officerName)
		if err != nil {
			panic(err.Error())
		}
	}
	defer db.Close()
	theftCase.OFFICER.ID = officerId
	theftCase.OFFICER.NAME = officerName
	json.NewEncoder(w).Encode(&theftCase)
}

func UpdateCase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}
	var theftCase TheftCase
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&theftCase); err != nil {
		respondWithJSON(w, http.StatusBadRequest, err.Error())
		return
	}

	db := dbConn()
	updateResult, err := db.Prepare("UPDATE bike_thefts SET solved=? WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	_, err = updateResult.Exec(theftCase.SOLVED, id)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	resolved := "unresolved"
	if theftCase.SOLVED {
		resolved = "resolved"
	}
	message := "UPDATE bike theft ID:" + strconv.Itoa(theftCase.ID) + " to " + resolved
	fmt.Println(message)
	respondWithJSON(w, http.StatusBadRequest, message)
}

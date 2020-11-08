package internal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type TheftCase struct {
	ID          int            `json:"id"`
	TITLE       string         `json:"title"`
	BRAND       sql.NullString `json:"brand"`
	CITY        sql.NullString `json:"city"`
	DESCRIPTION sql.NullString `json:"description"`
	REPORTED    time.Time      `json:"reported"`
	UPDATED     time.Time      `json:"updated"`
	SOLVED      bool           `json:"solved"`
	OFFICER     Officer        //`json:"officer"`
	// image string ``
}

func CreateCase(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var theftCase TheftCase
	_ = json.NewDecoder(r.Body).Decode(&theftCase)

	json.NewEncoder(w).Encode(&theftCase)

	db := dbConn()
	insert, err := db.Prepare("INSERT INTO bike_thefts(title, brand, city, description) VALUES(?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}

	_, err = insert.Exec(theftCase.TITLE, theftCase.BRAND, theftCase.CITY, theftCase.DESCRIPTION)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Bile theft '" + theftCase.TITLE + "' is created")
	defer db.Close()
}

func GetCases(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()
	selResult, err := db.Query(`SELECT bt.id, bt.title, bt.brand, bt.city, bt.description, bt.reported, IFNULL(o.id, 0) AS officer_id, IFNULL(o.name, '') AS officer_name
								FROM bike_thefts bt
								LEFT JOIN officers o
								ON o.id = bt.officer
								ORDER BY bt.id DESC`)
	if err != nil {
		panic(err.Error())
	}

	theftCases := []TheftCase{}
	theftCase := TheftCase{}
	for selResult.Next() {
		var officerId int
		var officerName string
		err = selResult.Scan(&theftCase.ID, &theftCase.TITLE, &theftCase.BRAND,
			&theftCase.CITY, &theftCase.DESCRIPTION, &theftCase.REPORTED, &officerId, &officerName)
		if err != nil {
			panic(err.Error())
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
		respondWithJSON(w, http.StatusBadRequest, "Invalid theft case ID")
		return
	}

	db := dbConn()
	selResult, err := db.Query("SELECT id, title, brand, city, description, reported, officer FROM bike_thefts WHERE id=?", id)
	if err != nil {
		panic(err.Error())
	}

	theftCase := TheftCase{}
	for selResult.Next() {
		err = selResult.Scan(&theftCase.ID, &theftCase.TITLE, &theftCase.BRAND,
			&theftCase.CITY, &theftCase.DESCRIPTION, &theftCase.REPORTED, &theftCase.OFFICER)
		if err != nil {
			panic(err.Error())
		}
	}
	defer db.Close()
	json.NewEncoder(w).Encode(&theftCase)
}

package internal

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Officer struct {
	ID   int    `json:"id"`
	NAME string `json:"name"`
}
type Message struct {
	CONTENT string
}

var officers []Officer

func dbConn() (db *sql.DB) {
	db, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/golang")
	if err != nil {
		panic(err.Error())
	}

	return db
}

func respondWithJSON(w http.ResponseWriter, code int, message interface{}) {
	response, _ := json.Marshal(message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func GetOfficers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()
	selResult, err := db.Query("SELECT * FROM officers ORDER BY id DESC")
	if err != nil {
		panic(err.Error())
	}

	officer := Officer{}
	for selResult.Next() {
		var id int
		var name string
		err = selResult.Scan(&id, &name)
		if err != nil {
			panic(err.Error())
		}
		officer.ID = id
		officer.NAME = name
		officers = append(officers, officer)
	}

	json.NewEncoder(w).Encode(&officers)
}

func GetOfficer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, "Invalid officer ID")
		return
	}

	db := dbConn()
	selResult, err := db.Query("SELECT * FROM officers WHERE id=?", id)
	if err != nil {
		panic(err.Error())
	}
	officer := Officer{}
	for selResult.Next() {
		var id int
		var name string
		err = selResult.Scan(&id, &name)
		if err != nil {
			panic(err.Error())
		}
		officer.ID = id
		officer.NAME = name
	}

	json.NewEncoder(w).Encode(&officer)
}

func UpdateOfficer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, "Invalid officer ID")
		return
	}
	var officer Officer
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&officer); err != nil {
		respondWithJSON(w, http.StatusBadRequest, "Invalid request officer")
		return
	}
	officer.ID = id

	db := dbConn()
	updateResult, err := db.Prepare("UPDATE officers SET name=? WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	updateResult.Exec(officer.NAME, officer.ID)
	fmt.Println("UPDATE: Name: " + officer.NAME + " for ID:" + strconv.Itoa(officer.ID))
	json.NewEncoder(w).Encode(&officer)
}

func CreateOfficer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var officer Officer
	_ = json.NewDecoder(r.Body).Decode(&officer)

	//officer.ID = strconv.Itoa(rand.Intn(1000000))
	officers = append(officers, officer)
	json.NewEncoder(w).Encode(&officer)

	db := dbConn()
	insert, err := db.Prepare("INSERT INTO officers(name) VALUES(?)")
	if err != nil {
		panic(err.Error())
	}
	_, err = insert.Exec(officer.NAME)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Officer '" + officer.NAME + "' is deleted")
	defer db.Close()
}

func DeleteOfficer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		respondWithJSON(w, http.StatusBadRequest, "Invalue officer ID")
		return
	}

	db := dbConn()
	delResult, err := db.Prepare("DELETE FROM officers WHERE id=?")
	if err != nil {
		panic(err.Error())
	}
	delResult.Exec(id)
	messageContent := "officer with ID: " + params["id"] + " is deleted"
	fmt.Println(messageContent)

	message := Message{}
	message.CONTENT = messageContent
	json.NewEncoder(w).Encode(&message)
}
package internal

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func AssignCase() {
	var availableOfficerIds = getIdsOfAvailableOfficer()
	lenOfAvilableOfficer := len(availableOfficerIds)
	lenOfAvailableCase := getNumOfAvailableCase()

	updateNumOfCase := lenOfAvailableCase
	if updateNumOfCase > lenOfAvilableOfficer {
		updateNumOfCase = lenOfAvilableOfficer
	}

	for _, id := range availableOfficerIds[:updateNumOfCase] {
		updateBikeTheft(id)
	}
}

func getIdsOfAvailableOfficer() []int {
	db := dbConn()
	rows, err := db.Query(`SELECT o.id FROM officers o
		LEFT JOIN bike_thefts bt
		ON o.id = bt.officer AND bt.solved = 0
		WHERE bt.id IS NULL`)

	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	var officerIds []int
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			panic(err.Error())
		}
		officerIds = append(officerIds, id)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	return officerIds
}

func getNumOfAvailableCase() int {
	db := dbConn()
	availableCaseNum := 0

	err := db.QueryRow(
		`SELECT COUNT(id) FROM bike_thefts
		WHERE solved = 0 AND officer IS NULL`).Scan(&availableCaseNum)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	return availableCaseNum
}

func updateBikeTheft(officerId int) {
	db := dbConn()
	updateResult, err := db.Prepare("UPDATE bike_thefts SET officer=? WHERE solved = 0 AND officer IS NULL LIMIT 1")
	if err != nil {
		panic(err.Error())
	}
	updateResult.Exec(officerId)
	log.Println("UPDATE: bike_theft table with offficer id")
	defer db.Close()
}

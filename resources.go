package main

import (
	"log"
	"os"
	"timesheet/db"

	"github.com/jackc/pgx/v4/pgxpool"
)

var commandDB *pgxpool.Pool

func initResourcesOrFail() {
	if !initCommandDatabase() {
		log.Printf("Exiting as initializing command database failed")
		os.Exit(1)
	}

	initServices()
}

func initCommandDatabase() bool {

	pool, err := db.ConnectPool(config.CommandDatabase.URL)
	if err != nil {
		log.Printf("Unable to connect to Command database: %s\n", err.Error())
		return false
	}

	//run test query to see if everything is alright.
	if !db.IsPoolHealthy(pool) {
		log.Println("Pool is unhealthy")
		return false
	}

	commandDB = pool

	log.Printf("Command database initialized")
	return true
}

func initServices() {
	log.Println("Initialising services")
}

package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

// Connect is function which is used to ensure the connection is successful with the DB
func Connect(url string) (*pgx.Conn, error) {
	return pgx.Connect(context.Background(), url)
}

func ConnectPool(url string) (*pgxpool.Pool, error) {
	return pgxpool.Connect(context.Background(), url)
}

func Close(conn *pgx.Conn) {
	if conn != nil {
		conn.Close(context.Background())
		log.Printf("Postgres connection closed")
	} else {
		log.Printf("Postgres connection already closed")
	}
}

func ClosePool(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
		log.Printf("Postgres connectionpool closed")
	} else {
		log.Printf("Postgres connectionpool already closed")
	}
}

func IsHealthy(conn *pgx.Conn) bool {

	retval := 0
	var err error

	err = conn.QueryRow(context.Background(), "select 1").Scan(&retval)
	if err != nil {
		log.Printf("Failed running a test query. Perhaps database connection is not healthy: %v\n", err)
		return false
	}

	if retval != 1 {
		log.Printf("Test query returned [%d]. Expected [%d]. Perhaps database connection is not healthy", retval, 1)
		return false
	}

	return true
}

func IsPoolHealthy(pool *pgxpool.Pool) bool {

	retval := 0
	var err error

	err = pool.QueryRow(context.Background(), "select 1").Scan(&retval)
	if err != nil {
		log.Printf("Failed running a test query. Perhaps database connection is not healthy: %v\n", err)
		return false
	}

	if retval != 1 {
		log.Printf("Test query returned [%d]. Expected [%d]. Perhaps database connection is not healthy", retval, 1)
		return false
	}

	return true
}

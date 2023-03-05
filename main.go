package main

import (
	"database/sql"
	"goApp/api"
	db "goApp/db/sqlc"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	psqlInfo := "host=localhost port=5432 user=postgres password=postgres dbname=simple_bank sslmode=disable"
	conn, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	addr := "localhost:3000"
	err = server.Start(addr)
	if err != nil {
		log.Fatal(err)
	}

}

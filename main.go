package main

import (
	"database/sql"
	"goApp/api"
	db "goApp/db/sqlc"
	"goApp/util"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	conn, err := sql.Open("postgres", config.PostgresInfo)
	if err != nil {
		log.Fatal(err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store, config)

	err = server.Start(config.Addr)
	if err != nil {
		log.Fatal(err)
	}

}

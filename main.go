package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/minhdang2803/simple_bank/api"
	db "github.com/minhdang2803/simple_bank/db/sqlc"
	"github.com/minhdang2803/simple_bank/utils"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load configuration", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to DB", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
}

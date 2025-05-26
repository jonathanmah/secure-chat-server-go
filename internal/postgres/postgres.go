package postgres

import (
	"chatapp/internal/config"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	connStr := config.App.PG.PgConnString()
	var err error
	DB, err = sql.Open(config.App.PG.DriverName, connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("Could not ping DB: ", err)
	}
	log.Println("Connected to DB successfully")
}

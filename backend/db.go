package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// Db es una variable global que representa la conexión a la base de datos
var Db *sql.DB

func InitDB(dataSourceName string) {
	var err error
	Db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}
	if err := Db.Ping(); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to the database")
}

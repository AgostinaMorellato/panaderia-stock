package main

import (
	"database/sql"
	"fmt"

	//"os"
	//"path/filepath"
	//"strings"

	_ "github.com/go-sql-driver/mysql"
)

// Db es una variable global que representa la conexión a la base de datos
var Db *sql.DB

// InitDB inicializa la conexion a la base de datos
func InitDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	return db, nil
}

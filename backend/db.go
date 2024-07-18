package main

/*
import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// Db es una variable global que representa la conexión a la base de datos
var Db *sql.DB

// InitDB initializes the database connection
func InitDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Verify the connection to the database
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	return db, nil
}
func executeSQLFile(db *sql.DB, filename string) error {
	filepath := filepath.Join("..", "db", filename)
	// Leer el contenido del archivo SQL
	sqlBytes, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	// Convertir el contenido del archivo a string
	sql := string(sqlBytes)

	// Remover cualquier caracter de retorno de carro adicional que pueda causar problemas
	sql = strings.ReplaceAll(sql, "\r", "")

	// Dividir el contenido del archivo en sentencias individuales
	sqlStatements := strings.Split(sql, ";")

	// Ejecutar cada sentencia individualmente
	for _, statement := range sqlStatements {
		statement = strings.TrimSpace(statement)
		if statement != "" {
			_, err = db.Exec(statement)
			if err != nil {
				return fmt.Errorf("error ejecutando la sentencia %q: %v", statement, err)
			}
		}
	}

	return nil
}*/

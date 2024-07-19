package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var db *sql.DB // Declaración de la variable global db

// Función para inicializar la conexión a la base de datos
func InitDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Verificar la conexión a la base de datos
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	return db, nil
}

// Función para ejecutar un archivo SQL
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
}

type Insumo struct {
	ID       int    `json:"id"`
	Nombre   string `json:"nombre"`
	Cantidad int    `json:"cantidad"`
	Unidad   string `json:"unidad"`
}

func main() {
	// DataSourceName formato: username:password@protocolo(dirección)/nombredb
	dataSourceName := "ly5zxxzt8321adlw:s9k9g8o3jihkteld@tcp(k9xdebw4k3zynl4u.cbetxkdyhwsb.us-east-1.rds.amazonaws.com:3306)/fywxazl2ckiibk60?multiStatements=true"

	var err error
	db, err = InitDB(dataSourceName) // Inicializa la variable global db
	if err != nil {
		log.Fatalf("Error al inicializar la base de datos: %v", err)
	}
	defer db.Close()

	log.Println("Conexión exitosa a la base de datos!")

	err = executeSQLFile(db, "init.sql")
	if err != nil {
		log.Fatalf("Error al ejecutar el archivo SQL: %v\n", err)
	}

	log.Println("Archivo SQL ejecutado correctamente")

	// Crear el router
	router := mux.NewRouter()
	router.HandleFunc("/api/stock", getStock).Methods("GET")
	router.HandleFunc("/api/stock", addInsumo).Methods("POST")
	router.HandleFunc("/api/stock/{id}", deleteInsumo).Methods("DELETE")
	router.HandleFunc("/api/stock/{id}", updateInsumo).Methods("PUT")

	allowedOrigins := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})

	// Inicia el servidor con los manejadores de CORS
	http.ListenAndServe(":8080", handlers.CORS(allowedOrigins, allowedMethods, allowedHeaders)(router))

}

func getStock(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		http.Error(w, "Base de datos no inicializada", http.StatusInternalServerError)
		return
	}

	rows, err := db.Query("SELECT id, nombre, cantidad, unidad FROM stock")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var insumos []Insumo
	for rows.Next() {
		var insumo Insumo
		err := rows.Scan(&insumo.ID, &insumo.Nombre, &insumo.Cantidad, &insumo.Unidad)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		insumos = append(insumos, insumo)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(insumos)
}

func addInsumo(w http.ResponseWriter, r *http.Request) {
	var insumo Insumo
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&insumo); err != nil {
		http.Error(w, "Error al decodificar la solicitud", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Validar los campos del insumo
	if insumo.Nombre == "" || insumo.Cantidad <= 0 {
		http.Error(w, "Revisar que nombre, cantidad y unidad estén completos y sean válidos", http.StatusBadRequest)
		return
	}

	// Insertar el insumo en la base de datos
	result, err := db.Exec("INSERT INTO stock (nombre, cantidad, unidad) VALUES (?, ?, ?)", insumo.Nombre, insumo.Cantidad, insumo.Unidad)
	if err != nil {
		http.Error(w, "Error al agregar el insumo en la base de datos", http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Error al obtener el ID del insumo insertado", http.StatusInternalServerError)
		return
	}
	insumo.ID = int(id)

	// Devolver el insumo con el ID asignado en formato JSON y código de estado 201
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(insumo)
}

func deleteInsumo(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de insumo inválido", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM stock WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Error al eliminar el insumo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func updateInsumo(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de insumo inválido", http.StatusBadRequest)
		return
	}

	var insumo Insumo
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&insumo); err != nil {
		http.Error(w, "Error al decodificar la solicitud", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	// Validar los campos del insumo
	if insumo.Cantidad <= 0 {
		http.Error(w, "Cantidad y unidad deben ser valores válidos y mayores a cero", http.StatusBadRequest)
		return
	}

	// Actualizar la cantidad del insumo en la base de datos
	_, err = db.Exec("UPDATE stock SET cantidad = ? WHERE id = ?", insumo.Cantidad, id)
	if err != nil {
		http.Error(w, "Error al actualizar la cantidad del insumo", http.StatusInternalServerError)
		return
	}

	insumo.ID = id

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(insumo)
}

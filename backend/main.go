package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

type Insumo struct {
	ID       int    `json:"id"`
	Nombre   string `json:"nombre"`
	Cantidad int    `json:"cantidad"`
}

func main() {
	// Inicializar la base de datos
	dataSourceName := "root:rootagos@tcp(localhost:3306)/panaderia_stock"
	InitDB(dataSourceName)

	// Crear el router
	router := mux.NewRouter()
	router.HandleFunc("/api/stock", getStock).Methods("GET")
	router.HandleFunc("/api/stock", addInsumo).Methods("POST")
	router.HandleFunc("/api/stock/{id}", deleteInsumo).Methods("DELETE")
	router.HandleFunc("/api/cocinar/{producto}", cocinar).Methods("POST")

	// Iniciar el servidor
	log.Println("Servidor ejecutándose en :3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}

func getStock(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, nombre, cantidad FROM insumos")
	if err != nil {
		http.Error(w, "Error al consultar el stock", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var insumos []Insumo
	for rows.Next() {
		var insumo Insumo
		if err := rows.Scan(&insumo.ID, &insumo.Nombre, &insumo.Cantidad); err != nil {
			http.Error(w, "Error al escanear filas", http.StatusInternalServerError)
			return
		}
		insumos = append(insumos, insumo)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(insumos)
}

func addInsumo(w http.ResponseWriter, r *http.Request) {
	var insumo Insumo
	if err := json.NewDecoder(r.Body).Decode(&insumo); err != nil {
		http.Error(w, "Error al decodificar la solicitud", http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO insumos (nombre, cantidad) VALUES (?, ?)", insumo.Nombre, insumo.Cantidad)
	if err != nil {
		http.Error(w, "Error al agregar el insumo", http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	insumo.ID = int(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(insumo)
}

func deleteInsumo(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID de insumo inválido", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM insumos WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Error al eliminar el insumo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func cocinar(w http.ResponseWriter, r *http.Request) {
	// Implementar lógica para cocinar
}

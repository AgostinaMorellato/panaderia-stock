package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"log"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var db *sql.DB // Declaración de la variable global db

// Función para inicializar la conexión a la base de datos
/*func InitDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	// Verificar la conexión a la base de datos
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	return db, nil
}*/

var PORT = getPort()

type Insumo struct {
	ID       int    `json:"id"`
	Nombre   string `json:"nombre"`
	Cantidad int    `json:"cantidad"`
	Unidad   string `json:"unidad"`
}

func getPort() string {
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "8080" // Puerto predeterminado si no se especifica
	}
	return PORT
}

func main() {
	// DataSourceName formato: username:password@protocolo(dirección)/nombredb
	dataSourceName := "wzdmwrg5qn734yj0:cfugoznnbbov4lr9@(rnr56s6e2uk326pj.cbetxkdyhwsb.us-east-1.rds.amazonaws.com:3306)/t41qgcrzm28aij2i?multiStatements=true"

	var err error
	db, err = InitDB(dataSourceName) // Inicializa la variable global db
	if err != nil {
		log.Fatalf("Error al inicializar la base de datos: %v", err)
	}
	defer db.Close()

	log.Println("Conexión exitosa a la base de datos!")

	// Crear el router
	router := mux.NewRouter()
	router.HandleFunc("/api/stock", getStock).Methods("GET")
	router.HandleFunc("/api/stock", addInsumo).Methods("POST")
	router.HandleFunc("/api/stock/{id}", deleteInsumo).Methods("DELETE")
	router.HandleFunc("/api/stock/{id}", updateInsumo).Methods("PUT")

	/*allowedOrigins := handlers.AllowedOrigins([]string{"https://panaderia-stock-frontend-app-6df615a13979.herokuapp.com"})
	allowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})*/

	// Inicia el servidor con los manejadores de CORS
	corsHandler := cors.Default().Handler(router)

	fmt.Printf("Server is running on :%s\n", PORT)
	log.Fatal(http.ListenAndServe(":"+PORT, corsHandler))

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

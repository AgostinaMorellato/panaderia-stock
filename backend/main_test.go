package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gorilla/mux"
)

func TestGetStock(t *testing.T) {
	log.Println("Running TestGetStock")

	Db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer Db.Close()

	db = Db

	rows := sqlmock.NewRows([]string{"id", "nombre", "cantidad", "unidad"}).
		AddRow(1, "Manteca", 10, "kg")
	mock.ExpectQuery("SELECT id, nombre, cantidad, unidad FROM stock").WillReturnRows(rows)

	req, err := http.NewRequest("GET", "/api/stock", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/stock", getStock).Methods("GET")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	log.Printf("Status Code: %d", rr.Code)
	log.Printf("Response Body: %s", rr.Body.String())

	/*expected := `[{"id":1,"nombre":"example","cantidad":10,"unidad":"kg"}]`
	if rr.Body.String() = expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}*/

	// Verificar que todas las expectativas fueron satisfechas
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestAddInsumo(t *testing.T) {
	log.Println("Running TestAddInsumo")

	Db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer Db.Close()

	db = Db

	newItem := Insumo{ID: 1, Nombre: "harina", Cantidad: 10, Unidad: "kg"}
	jsonData, err := json.Marshal(newItem)
	if err != nil {
		t.Fatal(err)
	}

	// Expect INSERT query and return a result with ID 1
	mock.ExpectExec("INSERT INTO stock").WithArgs(newItem.Nombre, newItem.Cantidad, newItem.Unidad).WillReturnResult(sqlmock.NewResult(1, 1))

	req, err := http.NewRequest("POST", "/api/stock", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(addInsumo)

	handler.ServeHTTP(rr, req)

	log.Printf("Status Code: %d", rr.Code)
	log.Printf("Response Body: %s", rr.Body.String())

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}

	// Verificar que el cuerpo de la respuesta contiene los datos del insumo creado
	var response Insumo
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.ID != 1 || response.Nombre != "harina" || response.Cantidad != 10 || response.Unidad != "kg" {
		t.Errorf("Expected created insumo to match input data, but got %+v", response)
	}
}
func TestDeleteInsumo(t *testing.T) {
	log.Println("Running TestDeleteInsumo")

	Db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer Db.Close()

	db = Db

	// Simular la eliminación del insumo
	mock.ExpectExec("DELETE FROM stock WHERE id = ?").
		WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))

	req, err := http.NewRequest("DELETE", "/api/stock/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/stock/{id}", deleteInsumo).Methods("DELETE")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Verificar que todas las expectativas fueron satisfechas
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestDescontarInsumo(t *testing.T) {
	log.Println("Running TestDescontarInsumo")

	Db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer Db.Close()

	db = Db

	// Simular la actualización de la cantidad del insumo
	mock.ExpectExec("UPDATE stock SET cantidad = cantidad - ? WHERE id = ?").
		WithArgs(5, 1).WillReturnResult(sqlmock.NewResult(1, 1))

	// Datos del insumo a descontar
	data := Insumo{Cantidad: 5}
	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", "/api/stock/1", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/stock/{id}", descontarInsumo).Methods("PUT")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Verificar que todas las expectativas fueron satisfechas
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

/*func TestDeleteInsumo(t *testing.T) {
	log.Println("Running TestDeleteInsumo")

	// Crear un ID de insumo para eliminar (por ejemplo, ID 1)
	id := 1

	// Mockear la base de datos para simular la eliminación del insumo
	Db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer Db.Close()

	db = Db

	// Expect DELETE query
	mock.ExpectExec("DELETE FROM stock WHERE id = ?").WithArgs(id).WillReturnResult(sqlmock.NewResult(0, 1))

	// Crear una solicitud HTTP DELETE con el ID del insumo a eliminar
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/api/stock/%d", id), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Ejecutar el handler deleteInsumo
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deleteInsumo)

	handler.ServeHTTP(rr, req)

	log.Printf("Status Code: %d", rr.Code)
	log.Printf("Response Body: %s", rr.Body.String())

	// Verificar que el código de estado sea 200 OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Verificar que el cuerpo de la respuesta esté vacío
	if rr.Body.String() != "" {
		t.Errorf("Handler returned unexpected body: got %v want %v",
			rr.Body.String(), "")
	}

	// Verificar que se completaron todas las expectativas mockeadas
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}*/

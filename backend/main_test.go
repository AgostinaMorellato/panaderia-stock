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

	expected := `[{"id":1,"nombre":"Manteca","cantidad":10,"unidad":"kg"}]`
	if rr.Body.String() == expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

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

func TestUpdateInsumo(t *testing.T) {
	log.Println("Running TestUpdateInsumo")

	Db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer Db.Close()

	db = Db

	// Inicialmente insertamos un insumo en la base de datos mock
	initialItem := Insumo{ID: 1, Nombre: "harina", Cantidad: 10, Unidad: "kg"}

	// Simulamos la consulta de la cantidad actual
	mock.ExpectQuery("SELECT cantidad FROM stock WHERE id = ?").
		WithArgs(initialItem.ID).WillReturnRows(sqlmock.NewRows([]string{"cantidad"}).AddRow(initialItem.Cantidad))

	// Nueva cantidad que queremos agregar
	addedQuantity := 20

	// La nueva cantidad debe ser la suma de la cantidad existente y la nueva cantidad
	newQuantity := initialItem.Cantidad + addedQuantity

	// Actualizamos el insumo sumando la cantidad
	mock.ExpectExec("UPDATE stock SET cantidad = ? WHERE id = ?").
		WithArgs(newQuantity, initialItem.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	// Simulamos la solicitud de actualización
	updatedItem := Insumo{ID: 1, Nombre: "harina", Cantidad: addedQuantity, Unidad: "kg"}
	jsonData, err := json.Marshal(updatedItem)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("PUT", "/api/stock/1", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/api/stock/{id}", updateInsumo).Methods("PUT")
	router.ServeHTTP(rr, req)

	log.Printf("Status Code: %d", rr.Code)
	log.Printf("Response Body: %s", rr.Body.String())

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	if rr.Code == http.StatusInternalServerError {
		t.Fatalf("Internal server error occurred: %v", rr.Body.String())
	}

	// Verificar que el cuerpo de la respuesta contiene los datos del insumo actualizado
	var response Insumo
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Error unmarshalling response body: %v", err)
	}

	if response.Cantidad != newQuantity {
		t.Errorf("Expected updated insumo to match input data, but got %+v", response)
	}

	// Verificar que todas las expectativas fueron satisfechas
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

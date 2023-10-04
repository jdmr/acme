package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/gorilla/mux"
)

var db *sql.DB

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var err error
	db, err = sql.Open("pgx", "postgres://acme:acme@localhost:5432/acme?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/customers", getCustomers).Methods("GET")
	r.HandleFunc("/api/v1/customers", createCustomer).Methods("POST")
	r.HandleFunc("/api/v1/customers/{customerID}", getCustomer).Methods("GET")
	r.HandleFunc("/api/v1/customers/{customerID}", deleteCustomer).Methods("DELETE")
	r.HandleFunc("/api/v1/customers/{customerID}", updateCustomer).Methods("PUT")

	r.HandleFunc("/api/v1/products", getProducts).Methods("GET")
	r.HandleFunc("/api/v1/products", createProduct).Methods("POST")
	r.HandleFunc("/api/v1/products/{productID}", getProduct).Methods("GET")
	r.HandleFunc("/api/v1/products/{productID}", deleteProduct).Methods("DELETE")
	r.HandleFunc("/api/v1/products/{productID}", updateProduct).Methods("PUT")

	r.HandleFunc("/api/v1/invoices", getInvoices).Methods("GET")
	r.HandleFunc("/api/v1/invoices", createInvoice).Methods("POST")
	r.HandleFunc("/api/v1/invoices/{invoiceID}", getInvoice).Methods("GET")
	r.HandleFunc("/api/v1/invoices/{invoiceID}", deleteInvoice).Methods("DELETE")
	r.HandleFunc("/api/v1/invoices/{invoiceID}", updateInvoice).Methods("PUT")

	log.Println("Server started on: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

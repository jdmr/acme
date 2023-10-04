package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type item struct {
	ID       string  `json:"id"`
	Quantity string  `json:"quantity"`
	Product  product `json:"product"`
}

type invoice struct {
	ID       string    `json:"id"`
	Customer customer  `json:"customer"`
	Date     time.Time `json:"date"`
	Items    []item    `json:"items"`
}

func getInvoices(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query(`
		select i.id, i.date, c.id, c.name
		from invoice i
		join customer c on c.id = i.customer_id
		order by i.date desc
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	invoices := []*invoice{}
	for rows.Next() {
		i := &invoice{}
		err := rows.Scan(&i.ID, &i.Date, &i.Customer.ID, &i.Customer.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rows2, err := db.Query(`
			select ii.id, ii.quantity, p.id, p.name, p.price
			from invoice_item ii
			join product p on p.id = ii.product_id
			where ii.invoice_id = $1
		`, i.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows2.Close()

		for rows2.Next() {
			it := item{}
			err := rows2.Scan(&it.ID, &it.Quantity, &it.Product.ID, &it.Product.Name, &it.Product.Price)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			i.Items = append(i.Items, it)
		}

		invoices = append(invoices, i)
	}

	result, err := json.Marshal(invoices)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func createInvoice(w http.ResponseWriter, r *http.Request) {
	var i invoice
	err := json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	i.ID = uuid.New().String()
	i.Date = time.Now()
	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = tx.Exec(`
		insert into invoice (id, customer_id, date)
		values ($1, $2, $3)
	`, i.ID, i.Customer.ID, i.Date)
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, it := range i.Items {
		_, err = tx.Exec(`
			insert into invoice_item (id, invoice_id, product_id, quantity)
			values ($1, $2, $3, $4)
		`, it.ID, i.ID, it.Product.ID, it.Quantity)
		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	tx.Commit()

	result, err := json.Marshal(i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(result)
}

func getInvoice(w http.ResponseWriter, r *http.Request) {
	var i invoice
	err := db.QueryRow(`
		select i.id, i.date, c.id, c.name
		from invoice i
		join customer c on c.id = i.customer_id
		where i.id = $1
	`, mux.Vars(r)["invoiceID"]).Scan(&i.ID, &i.Date, &i.Customer.ID, &i.Customer.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := db.Query(`
		select ii.id, ii.quantity, p.id, p.name, p.price
		from invoice_item ii
		join product p on p.id = ii.product_id
		where ii.invoice_id = $1
	`, i.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		it := item{}
		err := rows.Scan(&it.ID, &it.Quantity, &it.Product.ID, &it.Product.Name, &it.Product.Price)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		i.Items = append(i.Items, it)
	}

	result, err := json.Marshal(i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func deleteInvoice(w http.ResponseWriter, r *http.Request) {
	_, err := db.Exec(`
		delete from invoice
		where id = $1
	`, mux.Vars(r)["invoiceID"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func updateInvoice(w http.ResponseWriter, r *http.Request) {
	var i invoice
	err := json.NewDecoder(r.Body).Decode(&i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = tx.Exec(`
		update invoice
		set customer_id = $1
		where id = $2
	`, i.Customer.ID, mux.Vars(r)["invoiceID"])
	if err != nil {
		tx.Rollback()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, it := range i.Items {
		_, err = tx.Exec(`
			update invoice_item
			set product_id = $1, quantity = $2
			where id = $3
		`, it.Product.ID, it.Quantity, it.ID)
		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	tx.Commit()

	result, err := json.Marshal(i)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

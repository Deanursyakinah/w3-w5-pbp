package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	m "week5/model"

	"github.com/gorilla/mux"
)

func GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT * FROM transactions"

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return
	}
	var transaction m.Transaction
	var transactions []m.Transaction
	for rows.Next() {
		if err := rows.Scan(&transaction.ID, &transaction.UserID, &transaction.ProductID, &transaction.Quantity); err != nil {
			log.Println(err)
			return
		} else {
			transactions = append(transactions, transaction)
		}
	}
	if len(transactions) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		var response m.ProductsResponse
		response.Status = 404
		response.Message = "Data not found"
		response.Data = nil
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var response m.TransactionsResponse
	response.Status = 200
	response.Message = "Succes"
	response.Data = transactions
	json.NewEncoder(w).Encode(response)
}

func UpdateTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	userID := r.FormValue("userID")
	productID := r.FormValue("productID")
	quantity := r.FormValue("quantity")
	transactionID := mux.Vars(r)["id"]

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM transactions WHERE id = ?", transactionID).Scan(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count == 0 {
		http.Error(w, "transaction not found", http.StatusNotFound)
		return
	}

	query := "UPDATE transactions SET userID = ?, productID = ?, quantity = ? WHERE id = ?"

	_, err = db.Exec(query, userID, productID, quantity, transactionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "data transaction berhasil di update!")
}

func DeleteTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	productID := mux.Vars(r)["id"]

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM transactions WHERE id = ?", productID).Scan(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count == 0 {
		http.Error(w, "transaction not found", http.StatusNotFound)
		return
	}

	query := "DELETE FROM transactions WHERE id = ?"

	_, err = db.Exec(query, productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "data transaction berhasil dihapus!")
}
func GetDetailUserTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT t.ID, u.ID, u.name, u.age, u.address, p.ID, p.name, p.price, t.quantity FROM transactions t INNER JOIN users u ON t.UserID = u.ID INNER JOIN products p ON t.ProductID = p.ID"

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return
	}

	var transactions []m.Transaction
	for rows.Next() {
		var transaction m.Transaction
		var user m.User
		var product m.Product

		if err := rows.Scan(&transaction.ID, &user.ID, &user.Name, &user.Age, &user.Address, &product.ID, &product.Name, &product.Price, &transaction.Quantity); err != nil {
			log.Println(err)
			return
		} else {
			transaction.UserID = user
			transaction.ProductID = product
			transactions = append(transactions, transaction)
		}
	}

	if len(transactions) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		var response m.TransactionsResponse
		response.Status = 404
		response.Message = "Data not found"
		response.Data = nil
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var response m.TransactionsResponse
	response.Status = 200
	response.Message = "Succes"
	response.Data = transactions
	json.NewEncoder(w).Encode(response)
}

func GetDetailUserTransactionbyID(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	userID := mux.Vars(r)["id"]
	query := "SELECT t.ID, u.ID, u.name, u.age, u.address, p.ID, p.name, p.price, t.quantity FROM transactions t INNER JOIN users u ON t.UserID = u.ID INNER JOIN products p ON t.ProductID = p.ID WHERE u.ID = ?"

	rows, err := db.Query(query, userID)
	if err != nil {
		log.Println(err)
		return
	}

	var transactions []m.Transaction
	for rows.Next() {
		var transaction m.Transaction
		var user m.User
		var product m.Product

		if err := rows.Scan(&transaction.ID, &user.ID, &user.Name, &user.Age, &user.Address, &product.ID, &product.Name, &product.Price, &transaction.Quantity); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			var response m.TransactionsResponse
			response.Status = 404
			response.Message = "Data not found"
			response.Data = nil
			json.NewEncoder(w).Encode(response)
			log.Println(err)
			return
		} else {
			transaction.UserID = user
			transaction.ProductID = product

			transactions = append(transactions, transaction)
		}
	}

	if len(transactions) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		var response m.TransactionsResponse
		response.Status = 404
		response.Message = "Data not found"
		response.Data = nil
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var response m.TransactionsResponse
	response.Status = 200
	response.Message = "Succes"
	response.Data = transactions
	json.NewEncoder(w).Encode(response)
}

func InsertTransaction(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	if err := r.ParseForm(); err != nil {
		fmt.Println("Error parsing form data:", err)
		return
	}

	id := r.Form.Get("id")
	userID := r.Form.Get("userID")
	productID := r.Form.Get("productID")
	quantity := r.Form.Get("quantity")

	if userID == "" || productID == "" || quantity == "" {
		http.Error(w, "data tidak lengkap atau kelebihan", http.StatusBadRequest)
		return
	}

	if len(r.Form) >= 4 {
		http.Error(w, "Unexpected fields in form data", http.StatusBadRequest)
		fmt.Println("Unexpected fields in form data")
		return
	}

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM products WHERE id = ?", productID).Scan(&count)
	if err != nil {
		fmt.Println("product gaada di dalam database", err)
		return
	}

	if count == 0 {
		_, err := db.Exec("INSERT INTO products (id, name, price) VALUES (?, '', 0)", productID)
		if err != nil {
			fmt.Println("error insert new produk", err)
			return
		}
	}

	_, err = db.Exec("INSERT INTO transactions (id, userID, productID, quantity) VALUES (?, ?, ?, ?)", id, userID, productID, quantity)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		var response m.TransactionsResponse
		response.Status = 404
		response.Message = "Data not found"
		json.NewEncoder(w).Encode(response)
		fmt.Println("error insert new transaction", err)
		return
	}

	fmt.Fprintf(w, "berhasil insert transaction")
	w.Header().Set("Content-Type", "application/json")
	var response m.TransactionsResponse
	response.Status = 200
	response.Message = "Succes"
	json.NewEncoder(w).Encode(response)
}

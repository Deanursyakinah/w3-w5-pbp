package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	m "week4/model"

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

package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	m "week5/model"

	"github.com/gorilla/mux"
)

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT * FROM products"
	name := r.URL.Query()["name"]
	price := r.URL.Query()["price"]

	if name != nil {
		fmt.Println(name[0])
		query += " WHERE name= '" + name[0] + "'"
	}

	if price != nil {
		if name[0] != "" {
			query += "AND"
		} else {
			query += "WHERE"
		}
		query += " price= '" + price[0] + "'"
	}

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return
	}

	var product m.Product
	var products []m.Product
	for rows.Next() {
		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			log.Println(err)
			return
		} else {
			products = append(products, product)
		}
	}
	if len(products) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		var response m.ProductsResponse
		response.Status = 404
		response.Message = "Data not found"
		response.Data = nil
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var response m.ProductsResponse
	response.Status = 200
	response.Message = "Succes"
	response.Data = products
	json.NewEncoder(w).Encode(response)
}

func InsertNewProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form data", http.StatusInternalServerError)
		fmt.Println("Error parsing form data:", err)
		return
	}

	id := r.Form.Get("id")
	name := r.Form.Get("name")
	price := r.Form.Get("price")

	if name == "" || price == "" {
		http.Error(w, "Incomplete data provided", http.StatusBadRequest)
		return
	}

	if len(r.Form) >= 3 {
		http.Error(w, "Unexpected fields in form data", http.StatusBadRequest)
		fmt.Println("Unexpected fields in form data")
		return
	}

	query := "INSERT INTO products (id, name, price) VALUES (?, ?, ?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		http.Error(w, "Error preparing SQL statement", http.StatusInternalServerError)
		fmt.Println("Error preparing SQL statement:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, name, price)
	if err != nil {
		http.Error(w, "Error executing SQL statement", http.StatusInternalServerError)
		fmt.Println("Error executing SQL statement:", err)
		return
	}

	fmt.Fprintf(w, "new products berhasil di insert!")
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	name := r.FormValue("name")
	price := r.FormValue("price")
	productID := mux.Vars(r)["id"]

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM products WHERE id = ?", productID).Scan(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count == 0 {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	query := "UPDATE products SET name = ?, price = ? WHERE id = ?"

	_, err = db.Exec(query, name, price, productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "data product berhasil di update!")
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	productID := mux.Vars(r)["id"]

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM products WHERE id = ?", productID).Scan(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count == 0 {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}

	query := "DELETE FROM products WHERE id = ?"

	_, err = db.Exec(query, productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "data product berhasil dihapus!")
}

func DeleteSingleProduct(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	productID := mux.Vars(r)["id"]
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM transactions WHERE productID = ?", productID).Scan(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count == 0 {
		http.Error(w, "transaction not found", http.StatusNotFound)
		return
	}

	_, err = db.Exec("DELETE FROM transactions WHERE productID = ?", productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("DELETE FROM products WHERE id = ?", productID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "produk berhasil dihapus")
	w.Header().Set("Content-Type", "application/json")
	var response m.TransactionsResponse
	response.Status = 200
	response.Message = "Succes"
	json.NewEncoder(w).Encode(response)
}

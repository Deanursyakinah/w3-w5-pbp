package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	m "week5/model"

	"github.com/gorilla/mux"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		log.Println(err)
		return
	}

	defer rows.Close()

	var user m.User
	var users []m.User
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Address); err != nil {
			log.Println(err)
			return
		} else {
			users = append(users, user)
		}
	}
	if len(users) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		var response m.UsersResponse
		response.Status = 404
		response.Message = "Data not found"
		response.Data = nil
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	var response m.UsersResponse
	response.Status = 200
	response.Message = "Succes"
	response.Data = users
	json.NewEncoder(w).Encode(response)
}

func InsertNewUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form data", http.StatusInternalServerError)
		fmt.Println("Error parsing form data:", err)
		return
	}

	id := r.Form.Get("id")
	name := r.Form.Get("name")
	age := r.Form.Get("age")
	address := r.Form.Get("address")
	if name == "" || age == "" || address == "" {
		http.Error(w, "Incomplete data provided", http.StatusBadRequest)
		return
	}

	if len(r.Form) >= 4 {
		http.Error(w, "Unexpected fields in form data", http.StatusBadRequest)
		fmt.Println("Unexpected fields in form data")
		return
	}

	ageInt, err := strconv.Atoi(age)
	if err != nil {
		http.Error(w, "Invalid age", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO users (id, name, age, address) VALUES (?, ?, ?, ?)"
	stmt, err := db.Prepare(query)
	if err != nil {
		http.Error(w, "Error preparing SQL statement", http.StatusInternalServerError)
		fmt.Println("Error preparing SQL statement:", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, name, ageInt, address)
	if err != nil {
		http.Error(w, "Error executing SQL statement", http.StatusInternalServerError)
		fmt.Println("Error executing SQL statement:", err)
		return
	}

	fmt.Fprintf(w, "new user berhasil di insert!")
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	name := r.FormValue("name")
	age := r.FormValue("age")
	address := r.FormValue("address")
	userID := mux.Vars(r)["id"]

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", userID).Scan(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	query := "UPDATE users SET name = ?, age = ?, address = ? WHERE id = ?"

	_, err = db.Exec(query, name, age, address, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "data user berhasil di update!")
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	userID := mux.Vars(r)["id"]

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", userID).Scan(&count)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	query := "DELETE FROM users WHERE id = ?"

	_, err = db.Exec(query, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "data user berhasil dihapus!")
}

func Login(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		http.Error(w, "Email dan password harus diisi", http.StatusBadRequest)
		return
	}

	var (
		dbPassword string
		userID     int
		name       string
		age        int
		address    string
	)
	query := "SELECT id, password, name, age, address FROM users WHERE email = ?"
	err := db.QueryRow(query, email).Scan(&userID, &dbPassword, &name, &age, &address)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "email atau password salah", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if password != dbPassword {
		http.Error(w, "password salah", http.StatusUnauthorized)
		return
	}

	platform := r.Header.Get("platform")
	fmt.Fprintf(w, "Success login from %s", platform)
}

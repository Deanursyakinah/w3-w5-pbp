package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	m "week3/model"
)

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT * FROM users"
	name := r.URL.Query()["name"]
	age := r.URL.Query()["age"]

	if name != nil {
		fmt.Println(name[0])
		query += " WHERE name= '" + name[0] + "'"
	}

	if age != nil {
		if name[0] != "" {
			query += "AND"
		} else {
			query += "WHERE"
		}
		query += " age= '" + age[0] + "'"
	}

	if age != nil {
		if name[0] != "" {
			query += "AND"
		} else {
			query += "WHERE"
		}
		query += " age= '" + age[0] + "'"
	}

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return
	}
	var user m.User
	var users []m.User
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Age, &user.Address, &user.UserType); err != nil {
			log.Println(err)
			return
		} else {
			users = append(users, user)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	var response m.UsersResponse
	response.Status = 200
	response.Message = "Succes"
	response.Data = users
	json.NewEncoder(w).Encode(response)
}

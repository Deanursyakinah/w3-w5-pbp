package main

import (
	"fmt"
	"log"
	"net/http"
	"week4/controller"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	//ini yang buat get
	router.HandleFunc("/users", controller.GetAllUsers).Methods("GET")
	router.HandleFunc("/products", controller.GetAllProducts).Methods("GET")
	router.HandleFunc("/transactions", controller.GetAllTransactions).Methods("GET")

	//ini yang buat post
	router.HandleFunc("/users", controller.InsertNewUser).Methods("POST")
	router.HandleFunc("/products", controller.InsertNewProduct).Methods("POST")

	//ini buat update
	router.HandleFunc("/users/{id}", controller.UpdateUser).Methods("PUT")
	router.HandleFunc("/products/{id}", controller.UpdateProduct).Methods("PUT")
	router.HandleFunc("/transactions/{id}", controller.UpdateTransaction).Methods("PUT")

	//ini buat delete
	router.HandleFunc("/users/{id}", controller.DeleteUser).Methods("DELETE")
	router.HandleFunc("/products/{id}", controller.DeleteProduct).Methods("DELETE")
	router.HandleFunc("/transactions/{id}", controller.DeleteTransaction).Methods("DELETE")

	http.Handle("/", router)
	fmt.Println("Connected to port 8888")
	log.Println("Connected to port 8888")
	log.Fatal(http.ListenAndServe(":8080", router))
}

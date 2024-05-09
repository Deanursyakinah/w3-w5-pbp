package main

import (
	"fmt"
	"log"
	"net/http"
	"week5/controller"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	//-------------------------------ini endpoint minggu 5--------------------------------------
	router.HandleFunc("/v2/users", controller.InsertUserGORM).Methods("POST")
	router.HandleFunc("/v2/users/{id}", controller.UpdateUserGORM).Methods("PUT")
	router.HandleFunc("/v2/users/{id}", controller.DeleteUserGORM).Methods("DELETE")
	router.HandleFunc("/v2/users", controller.SelectUserGORM).Methods("GET")
	router.HandleFunc("/v2/users/{id}", controller.RawGorm).Methods("GET")

	//-------------------------------ini endpoint minggu 4--------------------------------------
	//ini no 1 endpoint
	router.HandleFunc("/v1/transactions", controller.GetDetailUserTransaction).Methods("GET")

	//ini endpoint no 2
	router.HandleFunc("/v1/product/{id}", controller.DeleteSingleProduct).Methods("DELETE")

	//ini endpoint no 3
	router.HandleFunc("/v1/transactions", controller.InsertTransaction).Methods("POST")

	//ini endpoint no 4
	router.HandleFunc("/v1/users", controller.Login).Methods("POST")

	//-------------------------------ini endpoint minggu 3--------------------------------------
	router.HandleFunc("/v1/products", controller.GetAllProducts).Methods("GET")
	router.HandleFunc("/v1/transactions", controller.GetAllTransactions).Methods("GET")

	//ini yang buat post
	router.HandleFunc("/v1/users", controller.InsertNewUser).Methods("POST")
	router.HandleFunc("/products", controller.InsertNewProduct).Methods("POST")

	//ini buat update
	router.HandleFunc("/v1/users/{id}", controller.UpdateUser).Methods("PUT")
	router.HandleFunc("/v1/products/{id}", controller.UpdateProduct).Methods("PUT")
	router.HandleFunc("/v1/transactions/{id}", controller.UpdateTransaction).Methods("PUT")

	//ini buat delete
	router.HandleFunc("/v1/users/{id}", controller.DeleteUser).Methods("DELETE")
	router.HandleFunc("/v1/products/{id}", controller.DeleteProduct).Methods("DELETE")
	router.HandleFunc("/v1/transactions/{id}", controller.DeleteTransaction).Methods("DELETE")

	http.Handle("/", router)
	fmt.Println("Connected to port 8888")
	log.Println("Connected to port 8888")
	log.Fatal(http.ListenAndServe(":8888", router))
}

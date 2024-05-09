package models

//kalau disuruh pakai gorm modelnya tambahin ini :
// type Transaction struct {
// 	ID			int json:"id,omitempty"
// 	UserID  	int json:"user_id,omitempty" gorm:"column:UserID;ForeignKey:ID"
// 	ProductID 	int json:"product_id,omitempty" gorm:"column:ProductID; ForeignKey:ID"
// 	Quantity 	int json:"quantity,omitempty"
// }
// type User struct {
// 	ID      int    json:"id,omitempty" gorm:"column:ID"
// 	Name    string json:"name,omitempty"
// 	Age     int    json:"age,omitempty"
// 	Address string json:"address,omitempty"
// 	Email string json:"email,omitempty"
// 	Password string json:"password,omitempty"
// }

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Address  string `json:"address"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type UserResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    User   `json:"data"`
}
type UsersResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []User `json:"data"`
}

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type Transaction struct {
	ID       int     `json:"id"`
	User     User    `json:"user"`
	Product  Product `json:"product"`
	Quantity int     `json:"quantity"`
}

type TransactionResponseDetail struct {
	Transaction []Transaction `json:"transaction"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type SuccessResponseGorm struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	ID      int    `json:"id"`
}

type TransactionResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    Transaction `json:"data,omitempty"`
}

type TransactionDetailsResponse struct {
	Status  int                       `json:"status"`
	Message string                    `json:"message"`
	Data    TransactionResponseDetail `json:"data"`
}

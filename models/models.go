package models

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"nombre"`
	Slug string `json:"slug"`
}

type Client struct {
	Id    int
	Name  string
	Phone string
	Email string
}

type Clients []Client

type Categories []Category

type ClientHttp struct {
	Css     string
	Message string
	Data    Clients
}

type ClientHttp2 struct {
	Css     string
	Message string
	Data    Client
}

type User struct {
	Id       int
	Name     string
	Email    string
	Phone    string
	Password string
}

type Users []User

type PaypalOrderResponseModel struct {
	Id            string                   `json:"id"`
	Status        string                   `json:"status"`
	PaymentSource map[string]interface{}   `json:"payment_source"`
	Links         []map[string]interface{} `json:"links"`
}

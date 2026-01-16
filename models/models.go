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

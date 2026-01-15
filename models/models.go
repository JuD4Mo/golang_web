package models

type Category struct {
	Id   int    `json:"id"`
	Name string `json:"nombre"`
	Slug string `json:"slug"`
}

type Categories []Category

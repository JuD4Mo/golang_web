package routes

import (
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

func Home(response http.ResponseWriter, request *http.Request) {
	template, err := template.ParseFiles("templates/example/home.html")
	if err != nil {
		panic(err)
	}
	template.Execute(response, nil)
}

func Params(response http.ResponseWriter, request *http.Request) {
	template, err := template.ParseFiles("templates/example/params.html")
	vars := mux.Vars(request)
	data := map[string]string{
		"id":   vars["id"],
		"slug": vars["slug"],
	}
	if err != nil {
		panic(err)
	}
	template.Execute(response, data)
}

// func ParamsQueryString(response http.ResponseWriter, request *http.Request) {
// 	template, err := template.ParseFiles("templates/example/home.html")
// 	vars := mux.Vars(request)
// 	data := map[string]string{
// 		"id": vars["id"],
// 		"slug": vars["slug"],
// 	}
// 	if err != nil {
// 		panic(err)
// 	}
// 	template.Execute(response, nil)
// }

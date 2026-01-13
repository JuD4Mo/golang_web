package routes

import (
	"html/template"
	"net/http"

	"github.com/JuD4Mo/golang-web/utilities"
	"github.com/gorilla/mux"
)

func Home(response http.ResponseWriter, request *http.Request) {
	template := template.Must(template.ParseFiles("templates/example/home.html", utilities.Frontend))
	template.Execute(response, nil)
	// template, err := template.ParseFiles("templates/example/home.html", "templates/layout/frontend.html")
	// if err != nil {
	// 	panic(err)
	// }
	// template.Execute(response, nil)
}

func Page404(response http.ResponseWriter, request *http.Request) {
	template := template.Must(template.ParseFiles("templates/example/404.html", utilities.Frontend))
	template.Execute(response, nil)
}

func Params(response http.ResponseWriter, request *http.Request) {
	template, err := template.ParseFiles("templates/example/params.html", utilities.Frontend)
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

func ParamsQueryString(response http.ResponseWriter, request *http.Request) {
	template, err := template.ParseFiles("templates/example/params_querystring.html", utilities.Frontend)
	id := request.URL.Query().Get("id")
	msg := request.URL.Query().Get("msg")
	data := map[string]string{
		"id":  id,
		"msg": msg,
	}
	if err != nil {
		panic(err)
	}
	template.Execute(response, data)
}

type Data struct {
	Name      string
	Age       int
	Type      int
	Abilities []Ability
}

type Ability struct {
	Name string
}

func Structs(response http.ResponseWriter, request *http.Request) {
	template, err := template.ParseFiles("templates/example/structs.html", utilities.Frontend)
	if err != nil {
		panic(err)
	}
	ability1 := Ability{Name: "Programming"}
	ability2 := Ability{Name: "Nutrition"}
	ability3 := Ability{Name: "Hacking"}
	abilities := []Ability{ability1, ability2, ability3}
	template.Execute(response, Data{Name: "Juan", Age: 20, Type: 1, Abilities: abilities})
}

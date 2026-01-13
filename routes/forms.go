package routes

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/JuD4Mo/golang-web/utilities"
	"github.com/JuD4Mo/golang-web/validations"
)

func Forms_get(response http.ResponseWriter, request *http.Request) {
	template := template.Must(template.ParseFiles("templates/forms/form.html", utilities.Frontend))
	template.Execute(response, nil)
	// template, err := template.ParseFiles("templates/example/home.html", "templates/layout/frontend.html")
	// if err != nil {
	// 	panic(err)
	// }
	// template.Execute(response, nil)
}

func Forms_post(response http.ResponseWriter, request *http.Request) {
	msg := ""
	if len(request.FormValue("name")) == 0 {
		msg += "Field name is empty"
	}
	if len(request.FormValue("email")) == 0 {
		msg += " Field email is empty"
	}

	if validations.Regex_correo.FindStringSubmatch(request.FormValue("email")) == nil {
		msg += " Email is not valid"
	}

	if validations.ValidatePassword(request.FormValue("password")) == false {
		msg += " Password is invalid"
	}

	if msg != "" {
		fmt.Fprintln(response, msg)
		return
	}
}

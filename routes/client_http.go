package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/JuD4Mo/golang-web/models"
	"github.com/JuD4Mo/golang-web/utilities"
	"github.com/gorilla/mux"
)

var token string = "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpZCI6MzYsImlhdCI6MTc2ODQ0OTQ5NSwiZXhwIjoxNzcxMDQxNDk1fQ.X8fshI-vUvrbrfDsJvErVeTeCqga1M8qNBTWjbx5vsc"

func Client_http(response http.ResponseWriter, request *http.Request) {
	template := template.Must(template.ParseFiles("templates/client_http/home.html", utilities.Frontend))

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.api.tamila.cl/api/categorias", nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Authorization", token)
	reg, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer reg.Body.Close()

	fmt.Println(reg.Status)

	body, err := io.ReadAll(reg.Body)
	if err != nil {
		fmt.Println(err)
	}

	//Info to slice
	dataCat := models.Categories{}
	err = json.Unmarshal(body, &dataCat)
	if err != nil {
		fmt.Println(err)
	}

	//css_session, css_message := utilities.ReturnFlashMessage(response, request)
	data := map[string]models.Categories{
		"body": dataCat,
	}

	template.Execute(response, data)

	// template, err := template.ParseFiles("templates/example/home.html", "templates/layout/frontend.html")
	// if err != nil {
	// 	panic(err)
	// }
	// template.Execute(response, nil)
}

func Client_http_create(response http.ResponseWriter, request *http.Request) {
	template := template.Must(template.ParseFiles("templates/client_http/create.html", utilities.Frontend))

	css_session, css_message := utilities.ReturnFlashMessage(response, request)
	data := map[string]string{
		"css":     css_session,
		"message": css_message,
	}

	template.Execute(response, data)
}

func Client_http_create_post(response http.ResponseWriter, request *http.Request) {
	msg := ""
	if len(request.FormValue("name")) == 0 {
		msg += "Field name is empty. "
	}

	if msg != "" {
		utilities.CreateFlashMessage(response, request, "danger", msg)
		http.Redirect(response, request, "/client-http/create", http.StatusSeeOther)
	}

	data := map[string]string{"nombre": request.FormValue("name")}
	jsonValue, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://www.api.tamila.cl/api/categorias", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Authorization", token)
	reg, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer reg.Body.Close()

	utilities.CreateFlashMessage(response, request, "success", "data created")
	http.Redirect(response, request, "/client-http/create", http.StatusSeeOther)
}

func Client_http_edit(response http.ResponseWriter, request *http.Request) {

	template := template.Must(template.ParseFiles("templates/client_http/edit.html", utilities.Frontend))

	vars := mux.Vars(request)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.api.tamila.cl/api/categorias/"+vars["id"], nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Authorization", token)

	reg, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer reg.Body.Close()
	body, err := io.ReadAll(reg.Body)
	//fmt.Printf("%s", body)
	//convertimos el resultado a un slice
	datos := models.Category{}
	errJson := json.Unmarshal(body, &datos)
	if errJson != nil {

	}

	data := map[string]string{

		"id":     vars["id"],
		"nombre": datos.Name,
		"slug":   datos.Slug,
	}
	template.Execute(response, data)
}

func Client_http_edit_post(response http.ResponseWriter, request *http.Request) {
	mensaje := ""
	if len(request.FormValue("nombre")) == 0 {
		mensaje = mensaje + "El campo Nombre está vacío. "
	}
	if mensaje != "" {

	}

	vars := mux.Vars(request)
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.api.tamila.cl/api/categorias/"+vars["id"], nil)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Authorization", token)

	reg, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer reg.Body.Close()
	body, err := io.ReadAll(reg.Body)
	//fmt.Printf("%s", body)
	//convertimos el resultado a un slice
	datos := models.Category{}
	errJson := json.Unmarshal(body, &datos)
	if errJson != nil {

	}
	//edito el registro
	datosJson := map[string]string{"nombre": request.FormValue("name")}
	jsonValue, _ := json.Marshal(datosJson)
	req2, err2 := http.NewRequest("PUT", "https://www.api.tamila.cl/api/categorias/"+vars["id"], bytes.NewBuffer(jsonValue))
	req2.Header.Set("Authorization", token)
	if err2 != nil {
		fmt.Println(err2)
	}
	reg2, err3 := client.Do(req2)
	defer reg.Body.Close()
	if err3 != nil {

	}

	defer reg2.Body.Close()

	http.Redirect(response, request, "/client-http/edit/"+vars["id"], http.StatusSeeOther)
}

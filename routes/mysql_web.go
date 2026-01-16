package routes

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/JuD4Mo/golang-web/db"
	"github.com/JuD4Mo/golang-web/models"
	"github.com/JuD4Mo/golang-web/utilities"
	"github.com/JuD4Mo/golang-web/validations"
	"github.com/gorilla/mux"
)

func Mysql_list(response http.ResponseWriter, request *http.Request) {
	template := template.Must(template.ParseFiles("templates/mysql/home.html", utilities.Frontend))

	//Conecct to DB
	db.Connect()

	sql := "select * from clientes order by id desc;"
	clients := models.Clients{}
	dataBd, err := db.Db.Query(sql)
	if err != nil {
		fmt.Println(err)
	}

	defer db.CloseDB()

	for dataBd.Next() {
		x := models.Client{}
		dataBd.Scan(&x.Id, &x.Name, &x.Email, &x.Phone)
		clients = append(clients, x)
	}

	//return

	css_session, css_message := utilities.ReturnFlashMessage(response, request)
	data := models.ClientHttp{
		Css:     css_session,
		Message: css_message,
		Data:    clients,
	}

	template.Execute(response, data)
}

func Mysql_create(response http.ResponseWriter, request *http.Request) {
	template := template.Must(template.ParseFiles("templates/mysql/create.html", utilities.Frontend))
	css_sesion, css_mensaje := utilities.ReturnFlashMessage(response, request)
	data := map[string]string{

		"css":     css_sesion,
		"message": css_mensaje,
	}
	template.Execute(response, data)
}

func Mysql_create_post(response http.ResponseWriter, request *http.Request) {
	mensaje := ""
	if len(request.FormValue("name")) == 0 {
		mensaje = mensaje + "El campo Nombre está vacío  "
	}
	if len(request.FormValue("email")) == 0 {
		mensaje = mensaje + " . El campo E-Mail está vacío "
	}

	if validations.Regex_correo.FindStringSubmatch(request.FormValue("email")) == nil {
		mensaje = mensaje + " . El E-Mail ingresado no es válido "
	}
	if len(request.FormValue("phone")) == 0 {
		mensaje = mensaje + " . El campo Teléfono está vacío "
	}
	if mensaje != "" {
		utilities.CreateFlashMessage(response, request, "danger", mensaje)
		http.Redirect(response, request, "/mysql/create", http.StatusSeeOther)
		return
	}
	db.Connect()
	sql := "INSERT into clientes values(null,? , ?, ?);"
	_, err := db.Db.Exec(sql, request.FormValue("name"), request.FormValue("email"), request.FormValue("phone"))
	if err != nil {
		fmt.Fprintln(response, err)
	}

	utilities.CreateFlashMessage(response, request, "success", "Se creó el registro exitosamente")
	http.Redirect(response, request, "/mysql/create", http.StatusSeeOther)
}

func Mysql_edit(response http.ResponseWriter, request *http.Request) {
	template := template.Must(template.ParseFiles("templates/mysql/edit.html", utilities.Frontend))

	db.Connect()
	sql := "SELECT id, nombre, correo, telefono FROM clientes where id=?"
	//clientes := modelos.Clientes{}

	vars := mux.Vars(request)
	datos, err := db.Db.Query(sql, vars["id"])
	if err != nil {
		fmt.Println(err)
	}
	defer db.CloseDB()
	var dato models.Client
	for datos.Next() {

		err := datos.Scan(&dato.Id, &dato.Name, &dato.Email, &dato.Phone)

		if err != nil {
			log.Fatal(err)
		}

	}
	css_sesion, css_mensaje := utilities.ReturnFlashMessage(response, request)
	clientehttp := models.ClientHttp2{
		Css:     css_sesion,
		Message: css_mensaje,
		Data:    dato,
	}
	template.Execute(response, clientehttp)
}

func Mysql_edit_post(response http.ResponseWriter, request *http.Request) {
	mensaje := ""
	vars := mux.Vars(request)
	if len(request.FormValue("name")) == 0 {
		mensaje = mensaje + "El campo Nombre está vacío  "
	}
	if len(request.FormValue("email")) == 0 {
		mensaje = mensaje + " . El campo E-Mail está vacío "
	}

	if validations.Regex_correo.FindStringSubmatch(request.FormValue("email")) == nil {
		mensaje = mensaje + " . El E-Mail ingresado no es válido "
	}
	if len(request.FormValue("phone")) == 0 {
		mensaje = mensaje + " . El campo Teléfono está vacío "
	}

	if mensaje != "" {
		utilities.CreateFlashMessage(response, request, "danger", mensaje)
		http.Redirect(response, request, "/mysql/edit/"+vars["id"], http.StatusSeeOther)
		return
	}
	db.Connect()

	sql := "update clientes set nombre=?, correo=?, telefono=? where id=?;"
	_, err := db.Db.Exec(sql, request.FormValue("name"), request.FormValue("email"), request.FormValue("phone"), vars["id"])
	if err != nil {
		fmt.Println(err)
	}
	utilities.CreateFlashMessage(response, request, "success", "Se modificó el registro exitosamente")
	http.Redirect(response, request, "/mysql/edit/"+vars["id"], http.StatusSeeOther)
}

func Mysql_delete(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	db.Connect()
	sql := "delete from clientes where id=?;"
	_, err := db.Db.Exec(sql, vars["id"])
	if err != nil {
		fmt.Fprintln(response, err)
	}
	utilities.CreateFlashMessage(response, request, "success", "Se eliminó el registro exitosamente")

	http.Redirect(response, request, "/mysql", http.StatusSeeOther)
}

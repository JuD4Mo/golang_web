package routes

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	"github.com/JuD4Mo/golang-web/db"
	"github.com/JuD4Mo/golang-web/models"
	"github.com/JuD4Mo/golang-web/utilities"
	"github.com/JuD4Mo/golang-web/validations"
	"golang.org/x/crypto/bcrypt"
)

func Register(response http.ResponseWriter, request *http.Request) {
	template := template.Must(template.ParseFiles("templates/security/register.html", utilities.Frontend))
	css_sesion, css_mensaje := utilities.ReturnFlashMessage(response, request)
	data := map[string]string{

		"css":     css_sesion,
		"message": css_mensaje,
	}
	template.Execute(response, data)
}

func Register_post(response http.ResponseWriter, request *http.Request) {
	mensaje := ""
	if len(request.FormValue("name")) == 0 {
		mensaje = mensaje + "El campo Nombre está vacío<br/> "
	}
	if len(request.FormValue("email")) == 0 {
		mensaje = mensaje + " . El campo E-Mail está vacío "
	}

	if validations.Regex_correo.FindStringSubmatch(request.FormValue("email")) == nil {
		mensaje = mensaje + " . El E-Mail ingresado no es válido "
	}
	if validations.ValidatePassword(request.FormValue("password")) == false {
		mensaje = mensaje + " . La contraseña debe tener al menos 1 número, una mayúscula, y un largo entre 6 y 20 caracteres "
	}
	if mensaje != "" {
		utilities.CreateFlashMessage(response, request, "danger", mensaje)
		http.Redirect(response, request, "/security/register", http.StatusSeeOther)
		return
	}

	db.Connect()

	defer db.CloseDB()

	sql := "INSERT INTO usuarios VALUES (null, ?,?,?,?);"

	//hash de password con bcrypt

	cost := 8
	bytes, _ := bcrypt.GenerateFromPassword([]byte(request.FormValue("password")), cost)

	_, err := db.Db.Exec(sql, request.FormValue("name"), request.FormValue("email"), request.FormValue("phone"),
		string(bytes))

	if err != nil {
		fmt.Fprintln(response, err)
	}

	utilities.CreateFlashMessage(response, request, "success", "User has been created")
	http.Redirect(response, request, "/security/register", http.StatusSeeOther)
}

func Login(response http.ResponseWriter, request *http.Request) {
	template := template.Must(template.ParseFiles("templates/security/login.html", utilities.Frontend))
	css_sesion, css_mensaje := utilities.ReturnFlashMessage(response, request)
	data := map[string]string{

		"css":     css_sesion,
		"message": css_mensaje,
	}
	template.Execute(response, data)
}

func Login_post(response http.ResponseWriter, request *http.Request) {
	mensaje := ""

	if len(request.FormValue("email")) == 0 {
		mensaje = mensaje + " . El campo E-Mail está vacío "
	}

	if validations.Regex_correo.FindStringSubmatch(request.FormValue("email")) == nil {
		mensaje = mensaje + " . El E-Mail ingresado no es válido "
	}
	if validations.ValidatePassword(request.FormValue("password")) == false {
		mensaje = mensaje + " . La contraseña debe tener al menos 1 número, una mayúscula, y un largo entre 6 y 20 caracteres "
	}
	if mensaje != "" {
		utilities.CreateFlashMessage(response, request, "danger", mensaje)
		http.Redirect(response, request, "/seguridad/login", http.StatusSeeOther)
	}

	db.Connect()

	sql := "SELECT id, nombre, correo, telefono, password FROM usuarios where correo=?;"
	data, err := db.Db.Query(sql, request.FormValue("email"))
	if err != nil {
		fmt.Fprintln(response, err)
	}

	defer db.CloseDB()

	var x models.User

	for data.Next() {
		err := data.Scan(&x.Id, &x.Name, &x.Email, &x.Phone, &x.Password)
		if err != nil {
			utilities.CreateFlashMessage(response, request, "danger", "Invalid credentials")
			http.Redirect(response, request, "/security/login", http.StatusSeeOther)
		}
	}

	//hash compare

	passwordBytes := []byte(request.FormValue("password"))
	passwordBd := []byte(x.Password)

	err = bcrypt.CompareHashAndPassword(passwordBd, passwordBytes)

	if err != nil {
		utilities.CreateFlashMessage(response, request, "danger", "invalid credentials")
		http.Redirect(response, request, "/security/login", http.StatusSeeOther)
		return
	}

	//Crear la sesión
	session, _ := utilities.Store.Get(request, "session-name")
	session.Values["id"] = strconv.Itoa(x.Id)
	session.Values["name"] = x.Name
	err = session.Save(request, response)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(response, request, "/security/protected", http.StatusSeeOther)
}

func Security_protected(response http.ResponseWriter, request *http.Request) {
	template := template.Must(template.ParseFiles("templates/security/protected.html", utilities.Frontend))
	css_sesion, css_mensaje := utilities.ReturnFlashMessage(response, request)
	id, name := utilities.ReturnLogin(request)
	data := map[string]string{

		"css":     css_sesion,
		"message": css_mensaje,
		"id":      id,
		"name":    name,
	}
	template.Execute(response, data)
}

func Logout(response http.ResponseWriter, request *http.Request) {
	session, _ := utilities.Store.Get(request, "session-name")
	session.Values["id"] = nil
	session.Values["name"] = nil
	err := session.Save(request, response)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	utilities.CreateFlashMessage(response, request, "primary", "Session closed")
	http.Redirect(response, request, "/security/login", http.StatusSeeOther)
}

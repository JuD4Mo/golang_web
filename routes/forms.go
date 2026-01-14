package routes

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/JuD4Mo/golang-web/utilities"
	"github.com/JuD4Mo/golang-web/validations"
)

func Forms_get(response http.ResponseWriter, request *http.Request) {
	template := template.Must(template.ParseFiles("templates/forms/form.html", utilities.Frontend))

	css_session, css_message := utilities.ReturnFlashMessage(response, request)
	data := map[string]string{
		"css":     css_session,
		"message": css_message,
	}

	template.Execute(response, data)

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
		// fmt.Fprintln(response, msg)
		// return
		utilities.CreateFlashMessage(response, request, "danger", msg)
		http.Redirect(response, request, "/forms", http.StatusSeeOther)
	}
}

func Forms_upload(response http.ResponseWriter, request *http.Request) {
	template := template.Must(template.ParseFiles("templates/forms/upload.html", utilities.Frontend))

	css_session, css_message := utilities.ReturnFlashMessage(response, request)
	data := map[string]string{
		"css":     css_session,
		"message": css_message,
	}

	template.Execute(response, data)

}

func Forms_upload_post(response http.ResponseWriter, request *http.Request) {
	file, handler, err := request.FormFile("image")
	if err != nil {
		utilities.CreateFlashMessage(response, request, "danger", "could not read uploaded file")
		http.Redirect(response, request, "/forms/upload", http.StatusSeeOther)
		return
	}
	defer file.Close()

	// Extract extension safely
	parts := strings.Split(handler.Filename, ".")
	if len(parts) < 2 {
		utilities.CreateFlashMessage(response, request, "danger", "invalid file name")
		http.Redirect(response, request, "/forms/upload", http.StatusSeeOther)
		return
	}
	extension := parts[len(parts)-1]

	// Use a timestamp for unique filename
	timestamp := time.Now().Format("20060102150405")
	image := timestamp + "." + extension

	// Ensure upload directory exists (relative path)
	uploadDir := "public/uploads/images"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		utilities.CreateFlashMessage(response, request, "danger", "could not create upload directory")
		http.Redirect(response, request, "/forms/upload", http.StatusSeeOther)
		return
	}

	fileImage := uploadDir + "/" + image

	f, err := os.OpenFile(fileImage, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		utilities.CreateFlashMessage(response, request, "danger", "could not create destination file")
		http.Redirect(response, request, "/forms/upload", http.StatusSeeOther)
		return
	}
	defer f.Close()

	if _, err := io.Copy(f, file); err != nil {
		utilities.CreateFlashMessage(response, request, "danger", "could not save uploaded file")
		http.Redirect(response, request, "/forms/upload", http.StatusSeeOther)
		return
	}

	utilities.CreateFlashMessage(response, request, "success", "file was uploaded correctly")
	http.Redirect(response, request, "/forms/upload", http.StatusSeeOther)
}

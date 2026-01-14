package routes

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/JuD4Mo/golang-web/utilities"
	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
	"github.com/xuri/excelize/v2"
	"gopkg.in/gomail.v2"
)

func Resources_get(response http.ResponseWriter, request *http.Request) {
	template := template.Must(template.ParseFiles("templates/resources/home.html", utilities.Frontend))

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

// ------------
func Resources_pdf(response http.ResponseWriter, request *http.Request) {
	template := template.Must(template.ParseFiles("templates/resources/pdf.html", utilities.Frontend))

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

func ImageFile(fileStr string) string {
	return filepath.Join(gofpdfDir, "public/images", fileStr)
}

var gofpdfDir string

func Filename(baseStr string) string {
	return PdfFile(baseStr + ".pdf")
}
func PdfFile(fileStr string) string {
	return filepath.Join(PdfDir(), fileStr)
}
func PdfDir() string {
	return filepath.Join(gofpdfDir, "public/pdf")
}
func Summary(err error, fileStr string) {
	if err == nil {
		fileStr = filepath.ToSlash(fileStr)
		fmt.Printf("Successfully generated %s\n", fileStr)
	} else {
		fmt.Println(err)
	}
}

func Resources_pdf_generate_better(response http.ResponseWriter, request *http.Request) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	// First page: manual local link
	pdf.AddPage()
	pdf.SetFont("Helvetica", "", 20)
	_, lineHt := pdf.GetFontSize()
	pdf.Write(lineHt, "To find out what's new in this tutorial, click ")
	pdf.SetFont("", "U", 0)
	link := pdf.AddLink()
	pdf.WriteLinkID(lineHt, "here", link)
	pdf.SetFont("", "", 0)
	// Second page: image link and basic HTML with link
	pdf.AddPage()
	pdf.SetLink(link, 0, -1)
	pdf.Image(ImageFile("logo.png"), 10, 12, 30, 0, false, "", 0, "http://www.fpdf.org")
	pdf.SetLeftMargin(45)
	pdf.SetFontSize(14)
	_, lineHt = pdf.GetFontSize()
	htmlStr := `You can now easily print text mixing different styles: <b>bold</b>, ` +
		`<i>italic</i>, <u>underlined</u>, or <b><i><u>all at once</u></i></b>!<br><br>` +
		`<center>You can also center text.</center>` +
		`<right>Or align it to the right.</right>` +
		`You can also insert links on text, such as ` +
		`<a href="http://www.fpdf.org">www.fpdf.org</a>, or on an image: click on the logo.`
	html := pdf.HTMLBasicNew()
	html.Write(lineHt, htmlStr)
	time := strings.Split(time.Now().String(), " ")
	nombre := string(time[4][6:14])
	fileStr := Filename(nombre)
	err := pdf.OutputFileAndClose(fileStr)
	Summary(err, fileStr)

	mensaje := "Se creó el documento PDF " + nombre + ".pdf de forma correcta"

	utilities.CreateFlashMessage(response, request, "success", mensaje)
	http.Redirect(response, request, "/resources/pdf", http.StatusSeeOther)
}

//---------------------------------------------

func Resources_pdf_generate(response http.ResponseWriter, request *http.Request) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Hello PDF")

	err := pdf.OutputFileAndClose("hello.pdf")
	if err != nil {

	}

	utilities.CreateFlashMessage(response, request, "success", "document created")
	http.Redirect(response, request, "/resources/pdf", http.StatusSeeOther)
}

func Resources_excel(response http.ResponseWriter, request *http.Request) {
	template := template.Must(template.ParseFiles("templates/resources/excel.html", utilities.Frontend))

	//Excel
	f := excelize.NewFile()
	defer func() {
		err := f.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	index, err := f.NewSheet("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}

	f.SetCellValue("Sheet1", "A1", "id")
	f.SetCellValue("Sheet1", "B1", "name")
	f.SetCellValue("Sheet1", "C1", "email")
	f.SetActiveSheet(index)

	//build document
	time := strings.Split(time.Now().String(), " ")
	name := string(time[4][6:14]) + ".xlsx"

	err = f.SaveAs("public/excel/" + name)
	if err != nil {
		fmt.Println(err)
	}
	//return

	css_session, css_message := utilities.ReturnFlashMessage(response, request)
	data := map[string]string{
		"css":     css_session,
		"message": css_message,
		"name":    name,
	}

	template.Execute(response, data)
}

func Resources_qr(response http.ResponseWriter, request *http.Request) {
	template := template.Must(template.ParseFiles("templates/resources/qr.html", utilities.Frontend))

	//Generate qr
	dataCodeImageQR, err := qrcode.Encode("https://www.youtube.com", qrcode.High, 256)
	if err != nil {
		log.Fatalln("Error generating QR", err)
	}
	img := base64.StdEncoding.EncodeToString(dataCodeImageQR)
	//Return

	css_session, css_message := utilities.ReturnFlashMessage(response, request)
	data := map[string]string{
		"css":     css_session,
		"message": css_message,
		"image":   img,
	}

	template.Execute(response, data)

	// template, err := template.ParseFiles("templates/example/home.html", "templates/layout/frontend.html")
	// if err != nil {
	// 	panic(err)
	// }
	// template.Execute(response, nil)
}

func Resources_send_email(response http.ResponseWriter, request *http.Request) {
	template := template.Must(template.ParseFiles("templates/resources/send_email.html", utilities.Frontend))

	//Generate qr
	dataCodeImageQR, err := qrcode.Encode("https://www.youtube.com", qrcode.High, 256)
	if err != nil {
		log.Fatalln("Error generating QR", err)
	}
	img := base64.StdEncoding.EncodeToString(dataCodeImageQR)
	//Return

	css_session, css_message := utilities.ReturnFlashMessage(response, request)
	data := map[string]string{
		"css":     css_session,
		"message": css_message,
		"image":   img,
	}

	template.Execute(response, data)

	// template, err := template.ParseFiles("templates/example/home.html", "templates/layout/frontend.html")
	// if err != nil {
	// 	panic(err)
	// }
	// template.Execute(response, nil)
}

func SendEmail() {
	msg := gomail.NewMessage()
	msg.SetHeader("From", "noreply@devcorebits.com")
	msg.SetHeader("To", "something@gmail.com")
	msg.SetHeader("Subject", "Go course")
	msg.SetBody("text/html", "something")
	//Configurar conexión con cliente SMTP
	n := gomail.NewDialer("host", 587, "username@something", "password") //envs

	err := n.DialAndSend()
	if err != nil {
		panic(err)
	}
}

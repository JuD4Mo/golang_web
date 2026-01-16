package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/JuD4Mo/golang-web/routes"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	mux := mux.NewRouter()
	mux.HandleFunc("/", routes.Home)
	mux.HandleFunc("/params/{id}/{slug}", routes.Params)
	mux.HandleFunc("/params_querystring", routes.ParamsQueryString)
	mux.HandleFunc("/structs", routes.Structs)

	mux.HandleFunc("/forms", routes.Forms_get)
	mux.HandleFunc("/forms-post", routes.Forms_post).Methods("POST")

	mux.HandleFunc("/forms/upload", routes.Forms_upload)
	mux.HandleFunc("/forms/upload-post", routes.Forms_upload_post).Methods("POST")

	mux.HandleFunc("/resources", routes.Resources_get)
	mux.HandleFunc("/resources/pdf", routes.Resources_pdf)
	mux.HandleFunc("/resources/pdf-generate", routes.Resources_pdf_generate_better)

	mux.HandleFunc("/resources/excel", routes.Resources_excel)

	mux.HandleFunc("/resources/qr", routes.Resources_qr)

	mux.HandleFunc("/resources/send-email", routes.Resources_send_email)

	mux.HandleFunc("/client-http", routes.Client_http)
	mux.HandleFunc("/client-http/create", routes.Client_http_create)
	mux.HandleFunc("/client-http/create-post", routes.Client_http_create_post).Methods("POST")
	mux.HandleFunc("/client-http/edit/{id}", routes.Client_http_edit)
	mux.HandleFunc("/client-http/edit-post/{id}", routes.Client_http_edit_post).Methods("POST")

	mux.HandleFunc("/client-http/delete/{id}", routes.Client_http_delete)

	//Archivos estáticos hacia mux
	s := http.StripPrefix("/public/", http.FileServer(http.Dir("./public/")))
	mux.PathPrefix("/public/").Handler(s)

	//404 error
	mux.NotFoundHandler = mux.NewRoute().HandlerFunc(routes.Page404).GetHandler()

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	server := &http.Server{
		Addr:         "localhost:" + os.Getenv("PORT"),
		Handler:      mux,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}

// func main() {
// 	//mux := http.NewServeMux()

// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprintln(w, "hello world")
// 	})

// 	log.Fatal(http.ListenAndServe("localhost:8081", nil))
// }

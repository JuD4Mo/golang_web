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

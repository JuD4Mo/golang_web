package protected

import (
	"net/http"

	"github.com/JuD4Mo/golang-web/utilities"
)

func Protected(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := utilities.Store.Get(r, "session-name")
		if session.Values["id"] != nil {
			next.ServeHTTP(w, r)
		} else {
			utilities.CreateFlashMessage(w, r, "warning", "You must be logged in to see this content")
			http.Redirect(w, r, "/security/login", http.StatusSeeOther)
		}

	}
}

// func Protected(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		next.ServeHTTP(w, r)
// 	}
// }

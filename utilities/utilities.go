package utilities

import (
	"net/http"

	"github.com/gorilla/sessions"
)

var Frontend string = "templates/layout/frontend.html"

var Store = sessions.NewCookieStore([]byte("session-name"))

func CreateFlashMessage(response http.ResponseWriter, request *http.Request, css, message string) {
	session, err := Store.Get(request, "flash-session")
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	session.AddFlash(css, "css")
	session.AddFlash(message, "message")
	session.Save(request, response)
}

func ReturnFlashMessage(response http.ResponseWriter, request *http.Request) (string, string) {
	// Use the same session name used when creating flashes
	session, _ := Store.Get(request, "flash-session")

	// Read both flash queues, then save the session once
	fm := session.Flashes("css")
	fm2 := session.Flashes("message")
	_ = session.Save(request, response)

	css_session := ""
	if len(fm) > 0 {
		css_session = fm[0].(string)
	}

	css_message := ""
	if len(fm2) > 0 {
		css_message = fm2[0].(string)
	}

	return css_session, css_message
}

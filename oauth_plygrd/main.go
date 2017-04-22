package main

import (
	"net/http"

	"github.com/goji/httpauth"
	"github.com/gorilla/mux"
	"strings"
)

func main() {

	myUnauthorizedHandler := genHandler()

	authOpts := httpauth.AuthOptions{
		Realm: "DevCo",
		AuthFunc: myAuthFunc,
		UnauthorizedHandler: myUnauthorizedHandler,
	}

	r := mux.NewRouter()

	r.HandleFunc("/", YourHandler)
	http.Handle("/", httpauth.BasicAuth(authOpts)(r))

	http.ListenAndServe(":7000", nil)
}

func YourHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Gorilla!\n"))
}

func myAuthFunc(user, pass string, r *http.Request) bool {
	return pass == strings.Repeat(user, 3)
}

func genHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized) // 404
		w.Write([]byte("Unauthorized"))
	})
}
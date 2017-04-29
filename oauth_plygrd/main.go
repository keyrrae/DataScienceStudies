package main

import (
	"net/http"

	"fmt"
	"github.com/goji/httpauth"
	_ "github.com/goji/httpauth"
	"github.com/gorilla/mux"
)

var cred map[string]string

func init() {
	cred = make(map[string]string)
	cred["signup"] = "signup"
}

func main() {

	myUnauthorizedHandler := genHandler()

	authOpts := httpauth.AuthOptions{
		Realm:               "DevCo",
		AuthFunc:            myAuthFunc,
		UnauthorizedHandler: myUnauthorizedHandler,
	}

	r := mux.NewRouter()

	r.HandleFunc("/", RootHandler) //.Method("GET")

	r.HandleFunc("/signup", SignupHandler).Methods("POST")
	//http.Handle("/signup", httpauth.BasicAuth(authOpts)(r))

	r.HandleFunc("/hh", YourHandler).Methods("GET")
	http.Handle("/", httpauth.BasicAuth(authOpts)(r))

	http.ListenAndServe(":7000", nil)
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Gorilla!\n"))

}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	rlen := r.ContentLength
	body := make([]byte, rlen)
	r.Body.Read(body)
	fmt.Fprintln(w, string(body))
}

func YourHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Gorilla!\n"))
}

func myAuthFunc(user, pass string, r *http.Request) bool {
	return pass == cred[user]
}

func genHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized) // 404
		w.Write([]byte("Unauthorized"))
	})
}

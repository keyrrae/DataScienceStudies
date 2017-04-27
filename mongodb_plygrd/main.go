package main

import (
	"encoding/base64"
	"net/http"
	"github.com/keyrrae/monimenta_backend/mongodb_plygrd/controllers"
	"gopkg.in/mgo.v2"
	"github.com/gorilla/mux"
	"strings"
	"fmt"
)

var cred map[string]string

func init() {
	cred = make(map[string]string)
	cred["signup"] = "signup"
}

func main() {
	r := mux.NewRouter()

	// Get a UserController instance
	uc := controllers.NewUserController(getSession())

	// Get a user resource
	r.HandleFunc("/user/:id", uc.GetUser).Methods("GET")

	r.HandleFunc("/", WrapAuthenticator(RootHandler, BasicAuth)).Methods("GET")

	// Create a new user
	r.HandleFunc("/user", uc.CreateUser).Methods("POST")

	// Remove an existing user
	r.HandleFunc("/user/:id", uc.RemoveUser).Methods("DELETE")

	// Fire up the server
	http.ListenAndServe("localhost:3000", r)
}

// getSession creates a new mongo session and panics if connection error occurs
func getSession() *mgo.Session {
	// Connect to our local mongo
	s, err := mgo.Dial("mongodb://admin:c0ffee@ds119091.mlab.com:19091/heroku_tqfnq24p")

	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}

	// Deliver session
	return s
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World, %s!", r.URL.Path[1:])
}


func AuthFunc(user, pass string) bool {
	return pass == cred[user]
}

func genHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized) // 404
		w.Write([]byte("Unauthorized"))
	})
}

func WrapAuthenticator(h http.HandlerFunc, middleware ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	for _, m := range middleware {
		h = m(h)
	}
	return h
}

func BasicAuth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)

		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			http.Error(w, "Not authorized", http.StatusUnauthorized)
			return
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			http.Error(w, err.Error(),  http.StatusUnauthorized)
			return
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			http.Error(w, "Not authorized",  http.StatusUnauthorized)
			return
		}

		if !AuthFunc(pair[0], pair[1]) {
			http.Error(w, "Not authorized",  http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	}
}

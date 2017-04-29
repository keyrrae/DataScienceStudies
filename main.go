package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/keyrrae/monimenta_backend/basicauth"
	"github.com/keyrrae/monimenta_backend/controllers"
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
	"os"
)

var cred map[string]string
var dbname = "heroku_tqfnq24p"

func init() {
	cred = make(map[string]string)
	cred["admin"] = "admin"
}

func main() {
	r := mux.NewRouter()

	// Get a UserController instance
	uc := controllers.NewUserController(getSession(), dbname)
	auth := basicauth.NewBasicAuthenticator(uc, cred)

	// Get a user resource
	authHandler := basicauth.WrapAuthenticator(RootHandler, auth.AdminAuth)
	r.HandleFunc("/auth", authHandler).Methods("POST")

	getUserHandler := basicauth.WrapAuthenticator(uc.CreateUser, auth.UserAuth)
	r.HandleFunc("/user/{userid}", getUserHandler).Methods("GET")

	r.HandleFunc("/", basicauth.WrapAuthenticator(RootHandler, auth.UserAuth)).Methods("GET")

	// Create a new user
	newUserHandler := basicauth.WrapAuthenticator(uc.CreateUser, auth.AdminAuth)
	r.HandleFunc("/user", newUserHandler).Methods("POST")

	// Remove an existing user
	removeUserHandler := basicauth.WrapAuthenticator(uc.CreateUser, auth.AdminAuth)
	r.HandleFunc("/user/{userid}", removeUserHandler).Methods("DELETE")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// Fire up the server
	http.ListenAndServe(":"+port, r)
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

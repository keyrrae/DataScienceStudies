package main


import (
	"github.com/gorilla/mux"
	"github.com/keyrrae/monimenta_backend/basicauth"
	"os"
	"log"
	"net/http"
	"gopkg.in/mgo.v2"
	"fmt"
	"github.com/keyrrae/monimenta_backend/controllers"
)

var cred map[string]string

func init() {
	cred = make(map[string]string)
	cred["admin"] = "admin"
}

func main() {
	r := mux.NewRouter()

	// Get a UserController instance
	uc := controllers.NewUserController(getSession())
	auth := basicauth.NewBasicAuthenticator(uc, cred)

	// Get a user resource

	r.HandleFunc("/user/{userid}", uc.GetUser).Methods("GET")

	r.HandleFunc("/", basicauth.WrapAuthenticator(RootHandler, auth.UserAuth)).Methods("GET")

	// Create a new user
	NewUserHandler := basicauth.WrapAuthenticator(uc.CreateUser, auth.AdminAuth)
	r.HandleFunc("/user", NewUserHandler).Methods("POST")

	// Remove an existing user
	r.HandleFunc("/user/:id", uc.RemoveUser).Methods("DELETE")

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	// Fire up the server
	http.ListenAndServe(":" + port, r)
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

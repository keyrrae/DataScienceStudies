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
	session := getSession()
	dbc := controllers.NewDbController(session, dbname)
	auth := basicauth.NewBasicAuthenticator(dbc, cred)

	r.HandleFunc("/", basicauth.WrapAuthenticator(RootHandler, auth.UserAuth)).Methods("GET")

	// ================================= User APIs ===========================================
	// Authenticate user login
	authHandler := basicauth.WrapAuthenticator(RootHandler, auth.AdminAuth)
	r.HandleFunc("/auth", authHandler).Methods("POST")

	// Create a new user
	newUserHandler := basicauth.WrapAuthenticator(dbc.CreateUser, auth.AdminAuth)
	r.HandleFunc("/user", newUserHandler).Methods("POST")

	// Get a user
	getUserHandler := basicauth.WrapAuthenticator(dbc.CreateUser, auth.UserAuth)
	r.HandleFunc("/user/{userid}", getUserHandler).Methods("GET")

	// Remove an existing user
	removeUserHandler := basicauth.WrapAuthenticator(dbc.CreateUser, auth.AdminAuth)
	r.HandleFunc("/user/{userid}", removeUserHandler).Methods("DELETE")

	// ================================= Town APIs ===========================================
	// Create a new town
	newTownHandler := basicauth.WrapAuthenticator(dbc.CreateTown, auth.UserAuth)
	r.HandleFunc("/town", newTownHandler).Methods("POST")

	// Remove an existing town
	removeTownHandler := basicauth.WrapAuthenticator(dbc.DeleteTown, auth.UserAuth)
	r.HandleFunc("/town/{townid}", removeTownHandler).Methods("DELETE")

	// Remove an existing town
	modifyTownHandler := basicauth.WrapAuthenticator(dbc.EditTown, auth.UserAuth)
	r.HandleFunc("/town/{townid}", modifyTownHandler).Methods("POST")

	// Get a town
	getTownHandler := basicauth.WrapAuthenticator(dbc.GetTown, auth.UserAuth)
	r.HandleFunc("/town/{townid}", getTownHandler).Methods("GET")

	userVisitHandler := basicauth.WrapAuthenticator(dbc.UserVisitTown, auth.UserAuth)
	r.HandleFunc("/town/{townid}/user/{userid}/visit", userVisitHandler).Methods("POST")

	updateSketchHandler := basicauth.WrapAuthenticator(dbc.UpdateSketch, auth.UserAuth)
	r.HandleFunc("/town/{townid}/sketch", updateSketchHandler).Methods("POST")

	geoSearchHandler := basicauth.WrapAuthenticator(dbc.GeoSearch, auth.UserAuth)
	r.HandleFunc("/geosearch", geoSearchHandler).Methods("GET")


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
	s, err := mgo.Dial("mongodb://admin:c0ffee@ds119091.mlab.com:19091/" + dbname)

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

package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/keyrrae/monimenta_backend/basicauth"
	"github.com/keyrrae/monimenta_backend/controllers"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Config struct {
	DbName        string `json:"dbname"`
	MdbUser       string `json:"mdbuser"`
	MdbHash       string `json:"mdbhash"`
	MdbInstance   string `json:"mdbinstance"`
	RedisInstance string `json:"redis_instance"`
	RedisHost     string `json:"redis_host"`
	RedisHash     string `json:"redis_hash"`
}

var cred map[string]string
var config Config

func init() {
	config = Config{}

	configJ, err := ioutil.ReadFile("./config.json")
	check(err)

	json.Unmarshal(configJ, &config)

	cred = make(map[string]string)
	cred["admin"] = "admin"
}

func main() {
	r := mux.NewRouter()

	// Get a UserController instance
	mgoSession := NewMongoSession()
	redisClient := NewRedisClient()
	dbc := controllers.NewDbController(mgoSession, redisClient, config.DbName)

	auth := basicauth.NewBasicAuthenticator(dbc, cred)

	r.HandleFunc("/", basicauth.WrapAuthenticator(RootHandler, auth.UserAuth)).Methods("GET")

	// ================================= User APIs ===========================================
	// Authenticate user login
	authHandler := basicauth.WrapAuthenticator(dbc.AuthUser, auth.AdminAuth)
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

	// Modify an existing town
	modifyTownHandler := basicauth.WrapAuthenticator(dbc.EditTown, auth.UserAuth)
	r.HandleFunc("/town/{townid}", modifyTownHandler).Methods("PUT")

	// Get a town
	getTownHandler := basicauth.WrapAuthenticator(dbc.GetTown, auth.UserAuth)
	r.HandleFunc("/town/{townid}", getTownHandler).Methods("GET")

	userVisitHandler := basicauth.WrapAuthenticator(dbc.UserVisitTown, auth.UserAuth)
	r.HandleFunc("/town/{townid}/user/{userid}/visit", userVisitHandler).Methods("POST")

	updateSketchHandler := basicauth.WrapAuthenticator(dbc.UpdateSketch, auth.UserAuth)
	r.HandleFunc("/town/{townid}/sketch", updateSketchHandler).Methods("PUT")

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
func NewMongoSession() *mgo.Session {
	// Connect to our local mongo
	mgoUrl := "mongodb://" + config.MdbUser + ":" + config.MdbHash + "@" + config.MdbInstance
	s, err := mgo.Dial( mgoUrl + config.DbName)

	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}

	// Deliver session
	return s
}

func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     config.RedisInstance + config.RedisHost,
		Password: config.RedisHash,
		DB:       0,
	})
	return client
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World, %s!", r.URL.Path[1:])
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

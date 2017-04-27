package controllers

import (
	"fmt"
	"net/http"
	"gopkg.in/mgo.v2"
	"github.com/keyrrae/monimenta_backend/models"
	"gopkg.in/mgo.v2/bson"
	"errors"
	"encoding/json"
	"github.com/gorilla/mux"
)

type (
	// UserController represents the controller for operating on the User resource
	UserController struct {
		session *mgo.Session
	}
)

// NewUserController provides a reference to a UserController with provided mongo session
func NewUserController(s *mgo.Session) *UserController {
	return &UserController{s}
}

func (uc UserController) GetUserFromDB(id string) (user models.User, err error) {
	u := models.User{}

	if !bson.IsObjectIdHex(id) {
		return u, errors.New("Invalid user id")
	}

	bid := bson.ObjectIdHex(id)

	// Fetch user
	if err := uc.session.DB("heroku_tqfnq24p").C("users").FindId(bid).One(&u); err != nil {
		return u, nil
	}

	return u, nil
}

// GetUser retrieves an individual user resource
func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	// Grab id
	params := mux.Vars(r)
	userid := params["userid"]

	if !bson.IsObjectIdHex(userid) {
		w.WriteHeader(404)
		return
	}

	oid := bson.ObjectIdHex(userid)

	// Stub user
	u := models.User{}

	// Fetch user
	if err := uc.session.DB("heroku_tqfnq24p").C("users").FindId(oid).One(&u); err != nil {
		w.WriteHeader(404)
		return
	}

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}

// CreateUser creates a new user resource
func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Stub an user to be populated from the body
	u := models.User{}

	// Populate the user data
	json.NewDecoder(r.Body).Decode(&u)

	fmt.Println(u)

	// Add an Id
	u.Id = bson.NewObjectId()

	// Write the user to mongo
	uc.session.DB("heroku_tqfnq24p").C("users").Insert(u)

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)
}

// RemoveUser removes an existing user resource
func (uc UserController) RemoveUser(w http.ResponseWriter, r *http.Request) {
	// Grab id
	params := mux.Vars(r)
	userid := params["userid"]

	if !bson.IsObjectIdHex(userid) {
		w.WriteHeader(404)
		return
	}

	oid := bson.ObjectIdHex(userid)

	// Remove user
	if err := uc.session.DB("heroku_tqfnq24p").C("users").RemoveId(oid); err != nil {
		w.WriteHeader(404)
		return
	}

	// Write status
	w.WriteHeader(200)
}
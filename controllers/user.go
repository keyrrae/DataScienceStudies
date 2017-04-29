package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/keyrrae/monimenta_backend/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

type (
	// UserController represents the controller for operating on the User resource
	UserController struct {
		dbname    string
		session   *mgo.Session
		tablename string
	}
)

// NewUserController provides a reference to a UserController with provided mongo session
func NewUserController(s *mgo.Session, db_name string) *UserController {
	return &UserController{dbname: db_name, session: s, tablename: "users"}
}

func (uc UserController) GetUserFromDB(id string) (user models.User, err error) {
	u := models.User{}

	if !bson.IsObjectIdHex(id) {
		return u, errors.New("Invalid user id")
	}

	bid := bson.ObjectIdHex(id)

	// Fetch user
	if err := uc.session.DB(uc.dbname).C(uc.tablename).FindId(bid).One(&u); err != nil {
		return u, nil
	}

	return u, nil
}

func (uc UserController) AuthUser(w http.ResponseWriter, r *http.Request) {
	u := models.User{}

	// Populate the user data
	json.NewDecoder(r.Body).Decode(&u)

	uindb := models.User{}
	err := uc.session.DB(uc.dbname).C(uc.tablename).Find(bson.M{"email": u.Email}).One(&uindb)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if u.Cred != uindb.Cred {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", uj)

}

// GetUser retrieves an individual user resource
func (uc UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	// Grab id
	params := mux.Vars(r)
	userid := params["userid"]

	if !bson.IsObjectIdHex(userid) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	oid := bson.ObjectIdHex(userid)

	// Stub user information
	// No credential information in UserInfo
	u := models.UserInfo{}

	// Fetch user
	if err := uc.session.DB(uc.dbname).C(uc.tablename).FindId(oid).One(&u); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(u)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", uj)
}

// CreateUser creates a new user resource
func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Stub an user to be populated from the body
	u := models.User{}

	// Populate the user data
	json.NewDecoder(r.Body).Decode(&u)

	fmt.Println(u)

	uindb := models.User{}
	err := uc.session.DB(uc.dbname).C(uc.tablename).Find(bson.M{"email": u.Email}).One(&uindb)
	if err != nil {
		//log.Fatal(err)
		// Add an Id
		u.Id = bson.NewObjectId()

		// Write the user to mongo
		uc.session.DB(uc.dbname).C(uc.tablename).Insert(u)

		// Marshal provided interface into JSON structure
		uj, _ := json.Marshal(u)

		// Write content-type, statuscode, payload
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "%s", uj)
		return
	}

	if u.Email == uindb.Email {
		w.WriteHeader(http.StatusConflict)
	}
}

// RemoveUser removes an existing user resource
func (uc UserController) RemoveUser(w http.ResponseWriter, r *http.Request) {
	// Grab id
	params := mux.Vars(r)
	userid := params["userid"]

	if !bson.IsObjectIdHex(userid) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	oid := bson.ObjectIdHex(userid)

	// Remove user
	if err := uc.session.DB(uc.dbname).C(uc.tablename).RemoveId(oid); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write status
	w.WriteHeader(http.StatusOK)
}

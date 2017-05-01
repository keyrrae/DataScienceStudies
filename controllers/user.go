package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/keyrrae/monimenta_backend/models"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func (c DbController) GetUserFromDB(id string) (user models.User, err error) {
	u := models.User{}

	if !bson.IsObjectIdHex(id) {
		return u, errors.New("Invalid user id")
	}

	bid := bson.ObjectIdHex(id)

	// Fetch user
	if err := c.Session.DB(c.DbName).C(c.UserTable).FindId(bid).One(&u); err != nil {
		return u, nil
	}

	return u, nil
}

func (c DbController) AuthUser(w http.ResponseWriter, r *http.Request) {
	// TODO: generate token instead of using credentials
	u := models.User{}

	// Populate the user data
	json.NewDecoder(r.Body).Decode(&u)

	uindb := models.User{}
	err := c.Session.DB(c.DbName).C(c.UserTable).Find(bson.M{"email": u.Email}).One(&uindb)
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
func (c DbController) GetUser(w http.ResponseWriter, r *http.Request) {
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
	if err := c.Session.DB(c.DbName).C(c.UserTable).FindId(oid).One(&u); err != nil {
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
func (c DbController) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Stub an user to be populated from the body
	u := models.User{}

	// Populate the user data
	json.NewDecoder(r.Body).Decode(&u)

	fmt.Println(u)

	uindb := models.User{}
	err := c.Session.DB(c.DbName).C(c.UserTable).Find(bson.M{"email": u.Email}).One(&uindb)
	if err != nil {
		//log.Fatal(err)
		// Add an Id
		u.Id = bson.NewObjectId()

		// Write the user to mongo
		c.Session.DB(c.DbName).C(c.UserTable).Insert(u)

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
func (c DbController) RemoveUser(w http.ResponseWriter, r *http.Request) {
	// Grab id
	params := mux.Vars(r)
	userid := params["userid"]

	if !bson.IsObjectIdHex(userid) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	oid := bson.ObjectIdHex(userid)

	// Remove user
	if err := c.Session.DB(c.DbName).C(c.UserTable).RemoveId(oid); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write status
	w.WriteHeader(http.StatusOK)
}

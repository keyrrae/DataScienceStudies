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
	// PostController represents the controller for operating on the Post resource
	PostController struct {
		dbname    string
		session   *mgo.Session
		tablename string
	}
)

// NewUserController provides a reference to a UserController with provided mongo session
func NewPostController(s *mgo.Session, databasename string) *PostController {
	return &PostController{dbname: databasename, session: s, tablename: "posts"}
}

func (c PostController) GetPostFromDB(id string) (user models.User, err error) {
	u := models.User{}

	if !bson.IsObjectIdHex(id) {
		return u, errors.New("Invalid user id")
	}

	bid := bson.ObjectIdHex(id)

	// Fetch user
	if err := c.session.DB(c.dbname).C(c.tablename).FindId(bid).One(&u); err != nil {
		return u, nil
	}

	return u, nil
}

// GetUser retrieves an individual user resource
func (c PostController) GetPost(w http.ResponseWriter, r *http.Request) {
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
	if err := c.session.DB(c.dbname).C(c.tablename).FindId(oid).One(&u); err != nil {
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
func (c PostController) CreatePost(w http.ResponseWriter, r *http.Request) {
	// Stub an user to be populated from the body
	u := models.User{}

	// Populate the user data
	json.NewDecoder(r.Body).Decode(&u)

	fmt.Println(u)

	uindb := models.User{}
	err := c.session.DB(c.dbname).C(c.tablename).Find(bson.M{"email": u.Email}).One(&uindb)

	if err != nil {
		//log.Fatal(err)
		// Add an Id
		u.Id = bson.NewObjectId()

		// Write the user to mongo
		c.session.DB("heroku_tqfnq24p").C("users").Insert(u)

		// Marshal provided interface into JSON structure
		uj, _ := json.Marshal(u)

		// Write content-type, statuscode, payload
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		fmt.Fprintf(w, "%s", uj)
		return
	}

	if u.Email == uindb.Email {
		w.WriteHeader(http.StatusConflict)
	}
}

// RemoveUser removes an existing user resource
func (c PostController) RemovePost(w http.ResponseWriter, r *http.Request) {
	// Grab id
	params := mux.Vars(r)
	userid := params["userid"]

	if !bson.IsObjectIdHex(userid) {
		w.WriteHeader(404)
		return
	}

	oid := bson.ObjectIdHex(userid)

	// Remove user
	if err := c.session.DB(c.dbname).C(c.tablename).RemoveId(oid); err != nil {
		w.WriteHeader(404)
		return
	}

	// Write status
	w.WriteHeader(200)
}

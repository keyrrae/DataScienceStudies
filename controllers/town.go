package controllers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/keyrrae/monimenta_backend/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"fmt"
)

type (
	// PostController represents the controller for operating on the Post resource
	TownController struct {
		dbname    string
		session   *mgo.Session
		tablename string
	}
)

// NewUserController provides a reference to a UserController with provided mongo session
func NewTownController(s *mgo.Session, databasename string) *TownController {
	return &TownController{dbname: databasename, session: s, tablename: "towns"}
}

// GetTown retrieves an individual user resource
func (c TownController) GetTown(w http.ResponseWriter, r *http.Request) {
	// Grab id
	params := mux.Vars(r)
	townId := params["townid"]

	if !bson.IsObjectIdHex(townId) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	databaseId := bson.ObjectIdHex(townId)

	// Stub town
	t := models.Town{}

	// Fetch user
	if err := c.session.DB(c.dbname).C(c.tablename).FindId(databaseId).One(&t); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Marshal provided interface into JSON structure
	tjson, _ := json.Marshal(t)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", tjson)
}

// CreateTown creates a new town entry
func (c TownController) CreateTown(w http.ResponseWriter, r *http.Request) {
	// Stub an user to be populated from the body
	t := models.Town{}

	// Populate the town data
	json.NewDecoder(r.Body).Decode(&t)

	// Add an Id
	t.Id = bson.NewObjectId()

	// insert in to database
	err := c.session.DB(c.dbname).C(c.tablename).Insert(t)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s", t.Id.Hex())
}

// RemoveTown removes an existing town entry
func (c TownController) DeleteTown(w http.ResponseWriter, r *http.Request) {
	// Grab id
	params := mux.Vars(r)
	townId := params["townid"]

	if !bson.IsObjectIdHex(townId) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	oid := bson.ObjectIdHex(townId)

	// Remove user
	if err := c.session.DB(c.dbname).C(c.tablename).RemoveId(oid); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write status
	w.WriteHeader(http.StatusOK)
}

// RemoveTown removes an existing town entry
func (c TownController) EditTown(w http.ResponseWriter, r *http.Request) {
	// Grab id
	params := mux.Vars(r)
	townId := params["townid"]

	if !bson.IsObjectIdHex(townId) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	t := models.Town{}

	// Populate the town data
	json.NewDecoder(r.Body).Decode(&t)

	t.Id = bson.ObjectIdHex(townId)

	// Remove user
	if err := c.session.DB(c.dbname).C(c.tablename).UpdateId(t.Id, t); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write status
	w.WriteHeader(http.StatusOK)
}
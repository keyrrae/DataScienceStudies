package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/keyrrae/monimenta_backend/models"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// GetTown retrieves an individual user resource
func (c DbController) GetTown(w http.ResponseWriter, r *http.Request) {
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
	if err := c.Session.DB(c.DbName).C(c.TownTable).FindId(databaseId).One(&t); err != nil {
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
func (c DbController) CreateTown(w http.ResponseWriter, r *http.Request) {

	// Stub an user to be populated from the body
	t := models.Town{}

	// Populate the town data
	json.NewDecoder(r.Body).Decode(&t)

	// Add an Id
	t.Id = bson.NewObjectId()

	// TODO: add a reference to this town in user
	// TODO: Locality sensitive hashing

	if t.Address == "" {

		getaddurl := "http://maps.googleapis.com/maps/api/geocode/json?latlng="
		getaddurl += strconv.FormatFloat(t.Lat, 'f', -1, 64) + "," + strconv.FormatFloat(t.Lng, 'f', -1, 64)
		getaddurl += "&sensor=false"

		req, _ := http.NewRequest("GET", getaddurl, nil)
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		add := models.Address{}
		json.Unmarshal(body, &add)

		if len(add.Results) > 0 {
			t.Address = add.Results[0].FormattedAddress
		}
	}

	// REF: http://maps.googleapis.com/maps/api/geocode/json?latlng=44.4647452,7.3553838&sensor=false

	// insert in to database
	err := c.Session.DB(c.DbName).C(c.TownTable).Insert(t)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%s", t.Id.Hex())
}

// RemoveTown removes an existing town entry
func (c DbController) DeleteTown(w http.ResponseWriter, r *http.Request) {
	// Grab id
	params := mux.Vars(r)
	townId := params["townid"]

	if !bson.IsObjectIdHex(townId) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	oid := bson.ObjectIdHex(townId)

	// Remove user
	if err := c.Session.DB(c.DbName).C(c.TownTable).RemoveId(oid); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write status
	w.WriteHeader(http.StatusOK)
}

// RemoveTown removes an existing town entry
func (c DbController) EditTown(w http.ResponseWriter, r *http.Request) {
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
	if err := c.Session.DB(c.DbName).C(c.TownTable).UpdateId(t.Id, t); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write status
	w.WriteHeader(http.StatusOK)
}

// UserVisitTown update visit town
func (c DbController) UserVisitTown(w http.ResponseWriter, r *http.Request) {
	// Grab id
	params := mux.Vars(r)
	townId := params["townid"]
	userId := params["userid"]

	if !bson.IsObjectIdHex(townId) || !bson.IsObjectIdHex(userId) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	townDbID := bson.ObjectIdHex(townId)

	// Stub town
	t := models.Town{}

	// Fetch town
	if err := c.Session.DB(c.DbName).C(c.TownTable).FindId(townDbID).One(&t); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	t.VisitedBy = append(t.VisitedBy, userId)

	// Remove user
	if err := c.Session.DB(c.DbName).C(c.TownTable).UpdateId(t.Id, t); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write status
	w.WriteHeader(http.StatusOK)
}

// UpdateSketch update the url to sketch
func (c DbController) UpdateSketch(w http.ResponseWriter, r *http.Request) {
	// Grab id
	params := mux.Vars(r)
	townId := params["townid"]

	if !bson.IsObjectIdHex(townId) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	townDbID := bson.ObjectIdHex(townId)

	// Stub town
	t := models.Town{}

	// Fetch town
	if err := c.Session.DB(c.DbName).C(c.TownTable).FindId(townDbID).One(&t); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	url, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	t.Sketch = string(url)

	if err := c.Session.DB(c.DbName).C(c.TownTable).UpdateId(t.Id, t); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write status
	w.WriteHeader(http.StatusOK)
}

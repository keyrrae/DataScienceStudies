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
	"crypto/sha1"
	"io"
)

func genGeoId(lat, lng float64) string {
	lat *= 10
	lng *= 10

	coOrd := models.CoOrd {
		Lat: int(lat),
		Lng: int(lng),
	}

	j, _ := json.Marshal(coOrd)
	h := sha1.New()
	io.WriteString(h, string(j))
	s := fmt.Sprintf("%x", h.Sum(nil))
	return s
}

func genIntGeoId(lat, lng int) string {

	coOrd := models.CoOrd {
		Lat: lat,
		Lng: lng,
	}

	j, _ := json.Marshal(coOrd)
	h := sha1.New()
	io.WriteString(h, string(j))
	s := fmt.Sprintf("%x", h.Sum(nil))
	return s
}

// GetTown retrieves an individual town entry
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

	// Fetch town
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

	// Generate geo hashing value
	t.GeoId = genGeoId(t.Lat, t.Lng)

	statusChan := make(chan int, 3)

	// Get the address from google api
	// REF: http://maps.googleapis.com/maps/api/geocode/json?latlng=44.4647452,7.3553838&sensor=false
	if t.Address == "" {
		getAddressUrl := "http://maps.googleapis.com/maps/api/geocode/json?latlng="
		getAddressUrl += strconv.FormatFloat(t.Lat, 'f', -1, 64) + "," + strconv.FormatFloat(t.Lng, 'f', -1, 64)
		getAddressUrl += "&sensor=false"

		req, _ := http.NewRequest("GET", getAddressUrl, nil)
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

	// Add a reference to this town in user
	go func() {
		u := models.User{}
		// Fetch user
		collection := c.Session.DB(c.DbName).C(c.UserTable)
		uid:= bson.ObjectIdHex(t.UserId)

		if err := collection.FindId(uid).One(&u); err != nil {
			// User not exist
			statusChan <- http.StatusNotFound
			fmt.Println(err)
			return // exit this goroutine
		}
		fmt.Println(u)
		u.TownIDs = append(u.TownIDs, string(t.Id))

		// change user
		if err := collection.UpdateId(u.Id, u); err != nil {
			statusChan <- http.StatusNotFound
			return
		}

		statusChan <- http.StatusOK
	} ()

	go func() {
		// insert in to database
		err := c.Session.DB(c.DbName).C(c.TownTable).Insert(t)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}

		statusChan <- http.StatusOK
	} ()

	for i:= 0; i < 2; i++{
		status := <- statusChan
		if status != http.StatusOK {
			w.WriteHeader(status)
			return
		}
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

// EditTown edits an existing town entry
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

	// Update town information
	if err := c.Session.DB(c.DbName).C(c.TownTable).UpdateId(t.Id, t); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Write status
	w.WriteHeader(http.StatusOK)
}

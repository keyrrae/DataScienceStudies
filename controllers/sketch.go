package controllers

import (
	"net/http"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"github.com/keyrrae/monimenta_backend/models"
)

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

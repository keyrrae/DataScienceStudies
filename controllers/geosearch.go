package controllers

import (
	"net/http"
	"github.com/keyrrae/monimenta_backend/models"
	"encoding/json"
	"math"
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"sort"
)

type TownByVisited []models.Town

func (s TownByVisited) Len() int {
	return len(s)
}
func (s TownByVisited) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s TownByVisited) Less(i, j int) bool {
	return len(s[i].VisitedBy) > len(s[i].VisitedBy)
}

func (c DbController) GeoSearch(w http.ResponseWriter, r *http.Request) {
	gs := models.GeoSearchRequest{}

	err := json.NewDecoder(r.Body).Decode(&gs)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	//geoID := genGeoId(gs.Lat, gs.Lng)
	center := models.CoOrd{
		Lat: int(gs.Lat * 10) ,
		Lng: int(gs.Lng * 10),
	}

	radius := int(math.Ceil(gs.Radius))

	step := 1
	numStep := int(2 * radius) + 1

	ul := models.CoOrd{
		Lat: center.Lat - radius,
		Lng: center.Lng - radius,
	}

	//townChan := make(chan *models.Town, numStep * numStep)

	townChan := make(chan []models.Town, numStep * numStep)

	for i := 0; i < numStep; i++ {
		for j := 0; j < numStep; j++ {
			geoId := genIntGeoId(ul.Lat + i * step, ul.Lng + j * step)
			fmt.Println(geoId)
			go func(){
				// Remove user

				var towns []models.Town
				collection := c.Session.DB(c.DbName).C(c.TownTable)
				if err := collection.Find(bson.M{"geoid": geoId}).All(&towns); err != nil {
					townChan <- nil
					return
				}
				fmt.Println(towns)

				townChan <- towns

			}()
		}
	}

	var res []models.Town

	for i := 0; i < numStep * numStep; i++ {
		town := <- townChan
		if town != nil {
			res = append(res, town...)
			//fmt.Println(town)
		}
	}

	sort.Sort(TownByVisited(res))
	if len(res) > gs.Count {
			res = res[0:gs.Count]
	}

	// Marshal provided interface into JSON structure
	townsJson, _ := json.Marshal(res)

	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", townsJson)
}
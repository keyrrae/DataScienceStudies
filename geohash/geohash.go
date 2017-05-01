package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"io"
)

type Coord struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func genGeoId(lat, lng float64) string {
	lat *= 10
	lng *= 10

	lat = float64(int(lat)) / 10
	lng = float64(int(lng)) / 10

	coord := Coord{Lat: lat, Lng: lng}

	j, _ := json.Marshal(coord)

	h := sha1.New()
	io.WriteString(h, string(j))
	sum := h.Sum(nil)

	s := fmt.Sprintf("%x", sum[0:12])
	return s
}

func main() {
	s := genGeoId(34.18, -119.234)
	fmt.Println(s)
	fmt.Println(bson.IsObjectIdHex(s))
}

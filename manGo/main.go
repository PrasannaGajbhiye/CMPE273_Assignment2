package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type LocationRequest struct {
	Name    string
	Address string
	City    string
	State   string
	Zip     string
}

type LocationLatLng struct {
	Lat float64
	Lng float64
}

type LocationResponse struct {
	Id         int64
	Name       string
	Address    string
	City       string
	State      string
	Zip        string
	Coordinate LocationLatLng
}

const (
	MongoDBHosts = "ds043694.mongolab.com:43694"
	AuthDatabase = "locationdb"
	AuthUserName = "mockrunuser"
	AuthPassword = "mockrunuser@273"
	TestDatabase = "locationdb"
)

// Create Location - POST Method
func createLocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

	// Connect to mongoDB server
	info := &mgo.DialInfo{
		Addrs:    []string{MongoDBHosts},
		Timeout:  60 * time.Second,
		Database: AuthDatabase,
		Username: AuthUserName,
		Password: AuthPassword,
	}
	session, err := mgo.DialWithInfo(info)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	db := session.DB("locationdb")
	//	collection := db.C("TestLocationCollection")
	collection := db.C("LocationCollection")

	// Fetch Input Json
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic("Error in reading body.")
	}

	var loc LocationRequest
	err = json.Unmarshal(body, &loc)
	if err != nil {
		panic("Error in unmarshalling.")
	}

	var locResp LocationResponse
	locResp.Name = loc.Name
	locResp.Address = loc.Address
	locResp.City = loc.City
	locResp.State = loc.State
	locResp.Zip = loc.Zip

	address := loc.Address + ", " + loc.City + ", " + loc.State
	address = strings.Replace(address, " ", "+", -1)

	// Fetch Coordinates
	resp, err := http.Get("http://maps.google.com/maps/api/geocode/json?address=" + address + "&sensor=false")
	if err != nil {
		fmt.Println("Error")
	}
	defer resp.Body.Close()

	body, _ = ioutil.ReadAll(resp.Body)
	var msgRes interface{}
	_ = json.Unmarshal(body, &msgRes)

	mRes := msgRes.(map[string]interface{})["results"]
	mRes0 := mRes.([]interface{})[0]
	mGeo := mRes0.(map[string]interface{})["geometry"]
	mLoc := mGeo.(map[string]interface{})["location"]

	locLat := mLoc.(map[string]interface{})["lat"].(float64)
	locLng := mLoc.(map[string]interface{})["lng"].(float64)

	var locCoordinates LocationLatLng
	locCoordinates.Lat = locLat
	locCoordinates.Lng = locLng

	locResp.Coordinate = locCoordinates

	// Check if documents exists
	var searchId LocationResponse
	err = collection.Find(nil).Sort("-id").One(&searchId)
	if err != nil {
		locResp.Id = 12345
	} else {
		locResp.Id = searchId.Id + 1
	}

	err = collection.Insert(locResp)
	if err != nil {
		log.Fatal(err)
	} else {
		mapB, _ := json.Marshal(locResp)
		fmt.Fprintf(rw, "\n"+string(mapB)+"\n")
	}

}

// Get Location - GET Method
func getLocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

	// Connect to mongoDB
	info := &mgo.DialInfo{
		Addrs:    []string{MongoDBHosts},
		Timeout:  60 * time.Second,
		Database: AuthDatabase,
		Username: AuthUserName,
		Password: AuthPassword,
	}
	session, err := mgo.DialWithInfo(info)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	db := session.DB("locationdb")
	//	collection := db.C("TestLocationCollection")
	collection := db.C("LocationCollection")

	// Fetch Input location_id
	location_idStr := p.ByName("location_id")
	location_id, _ := strconv.ParseInt(location_idStr, 10, 64)

	var locResp LocationResponse

	// Search for the location_id
	err = collection.Find(bson.M{"id": location_id}).One(&locResp)
	if err != nil {
		log.Fatal(err)
	}

	mapB, _ := json.Marshal(locResp)
	fmt.Fprintf(rw, "\n"+string(mapB)+"\n")
}

// Update Location Address - PUT Method
func updateLocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

	// Connect to mongoDB
	info := &mgo.DialInfo{
		Addrs:    []string{MongoDBHosts},
		Timeout:  60 * time.Second,
		Database: AuthDatabase,
		Username: AuthUserName,
		Password: AuthPassword,
	}
	session, err := mgo.DialWithInfo(info)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	db := session.DB("locationdb")
	//	collection := db.C("TestLocationCollection")
	collection := db.C("LocationCollection")

	// Fetch Input Json
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic("Error in reading body.")
	}

	var loc LocationRequest
	err = json.Unmarshal(body, &loc)
	if err != nil {
		panic("Error in unmarshalling.")
	}

	// Fetch Input location_id
	location_idStr := p.ByName("location_id")
	location_id, _ := strconv.ParseInt(location_idStr, 10, 64)

	address := loc.Address + ", " + loc.City + ", " + loc.State
	address = strings.Replace(address, " ", "+", -1)

	// Fetch Coordinates
	resp, err := http.Get("http://maps.google.com/maps/api/geocode/json?address=" + address + "&sensor=false")
	if err != nil {
		fmt.Println("Error")
	}
	defer resp.Body.Close()

	body, _ = ioutil.ReadAll(resp.Body)
	var msgRes interface{}
	_ = json.Unmarshal(body, &msgRes)

	mRes := msgRes.(map[string]interface{})["results"]
	mRes0 := mRes.([]interface{})[0]
	mGeo := mRes0.(map[string]interface{})["geometry"]
	mLoc := mGeo.(map[string]interface{})["location"]

	locLat := mLoc.(map[string]interface{})["lat"].(float64)
	locLng := mLoc.(map[string]interface{})["lng"].(float64)

	// Search for the location_id & update the location details
	err = collection.Update(bson.M{"id": location_id}, bson.M{"$set": bson.M{"address": loc.Address, "city": loc.City, "state": loc.State, "zip": loc.Zip, "coordinate.lat": locLat, "coordinate.lng": locLng}})
	if err != nil {
		log.Fatal(err)
	}

	var locResp LocationResponse

	// Search for the updated location_id
	err = collection.Find(bson.M{"id": location_id}).One(&locResp)
	if err != nil {
		log.Fatal(err)
	}

	mapB, _ := json.Marshal(locResp)
	fmt.Fprintf(rw, "\n"+string(mapB)+"\n")
}

// Remove Location - DELETE Method
func removeLocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {

	// Connect to mongoDB
	info := &mgo.DialInfo{
		Addrs:    []string{MongoDBHosts},
		Timeout:  60 * time.Second,
		Database: AuthDatabase,
		Username: AuthUserName,
		Password: AuthPassword,
	}
	session, err := mgo.DialWithInfo(info)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	db := session.DB("locationdb")
	//	collection := db.C("TestLocationCollection")
	collection := db.C("LocationCollection")

	// Fetch Input location_id
	location_idStr := p.ByName("location_id")
	location_id, _ := strconv.ParseInt(location_idStr, 10, 64)

	// Delete the location corresponding to the location_id
	err = collection.Remove(bson.M{"id": location_id})
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Fprintf(rw, "Location document deleted successfully.\n")
	}
}

func main() {

	mux := httprouter.New()
	mux.POST("/locations", createLocation)
	mux.GET("/locations/:location_id", getLocation)
	mux.PUT("/locations/:location_id", updateLocation)
	mux.DELETE("/locations/:location_id", removeLocation)
	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}
	server.ListenAndServe()
}

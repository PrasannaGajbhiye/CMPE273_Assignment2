# CMPE273_Assignment2
# CRUD Location Service

The location service has the following REST endpoints to store and retrieve locations. All the data persists into MongoDB. For Go application to MongoDB, driver mgo has been used.


## Usage

### Install

```
go get github.com/PrasannaGajbhiye/CMPE273_Assignment2
```

Start the  server:

```
go clean
go build
./manGo
```

### Start the client 
#### PUSH Request- For creation of new location
```
curl -H "Content-Type: application/json" -X POST -d '{"name" : "John Smith","address":"123 Main St","city": "San Francisco","state": "CA","zip":"94113"}' http://localhost:8080/locations
```
Following will be the response for the above request:
```
{"Id":12345,"Name":"John Smith","Address":"123 Main St","City":"San Francisco","State":"CA","Zip":"94113","Coordinate":{"Lat":37.7917618,"Lng":-122.3943405}}
```

#### PUSH Request- For creation of new location
```
curl -H "Content-Type: application/json" -X POST -d '{"name" : "John Smith","address":"123 Main St","city": "San Francisco","state": "CA","zip":"94113"}' http://localhost:8080/locations
```
Following will be the response for the above request:
```
{"Id":12345,"Name":"John Smith","Address":"123 Main St","City":"San Francisco","State":"CA","Zip":"94113","Coordinate":{"Lat":37.7917618,"Lng":-122.3943405}}
```

#### GET Request- For getting location details of a specific location_id
```
curl http://localhost:8080/locations/12345
```
Following will be the response for the above request:
```
{"Id":12345,"Name":"John Smith","Address":"123 Main St","City":"San Francisco","State":"CA","Zip":"94113","Coordinate":{"Lat":37.7917618,"Lng":-122.3943405}}
```

#### PUT Request- For updating location details of a specific location_id
```
curl -H "Content-Type: application/json" -X PUT -d '{"address" : "1600 Amphitheatre Parkway","city" : "Mountain View","state" : "CA","zip" : "94043"}' http://localhost:8080/locations/12345
```
Following will be the response for the above request:
```
{"Id":12345,"Name":"John Smith","Address":"1600 Amphitheatre Parkway","City":"Mountain View","State":"CA","Zip":"94043","Coordinate":{"Lat":37.4220352,"Lng":-122.0841244}}
```

#### DELETE Request- For deletion of new location
```
curl -X DELETE http://localhost:8080/locations/12345

```
Following will be the response for the above request:
```
Location document deleted successfully.
```

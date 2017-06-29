// This program provides a sample application for using MongoDB with
// the mgo driver.
package main

import (
	"log"
	"sync"
	"time"

	. "github.com/AlexsJones/choke/objects"
	mgo "gopkg.in/mgo.v2"
)

const (
	MongoDBHosts = "ds035428.mongolab.com:35428"
	AuthDatabase = "goinggo"
	AuthUserName = "guest"
	AuthPassword = "welcome"
	TestDatabase = "goinggo"
)

// main is the entry point for the application.
func main() {
	// We need this object to establish a session to our MongoDB.
	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    []string{MongoDBHosts},
		Timeout:  60 * time.Second,
		Database: AuthDatabase,
		Username: AuthUserName,
		Password: AuthPassword,
	}

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
	}

	// Reads may not be entirely up-to-date, but they will always see the
	// history of changes moving forward
	mongoSession.SetMode(mgo.Monotonic, true)

	// Create a wait group to manage the goroutines.
	var waitGroup sync.WaitGroup

	// Perform 10 concurrent queries against the database.
	waitGroup.Add(10)
	for query := 0; query < 10; query++ {
		go RunQuery(query, &waitGroup, mongoSession)
	}
	waitGroup.Wait()
	log.Println("All Queries Completed")
}

func RunQuery(query int, waitGroup *sync.WaitGroup, mongoSession *mgo.Session) {
	defer waitGroup.Done()

	// Request a socket connection from the session to process our query.
	sessionCopy := mongoSession.Copy()
	defer sessionCopy.Close()

	collection := sessionCopy.DB(TestDatabase).C("buoy_stations")

	log.Printf("RunQuery : %d : Executing\n", query)

	var buoyStations []BuoyStation
	err := collection.Find(nil).All(&buoyStations)
	if err != nil {
		log.Printf("RunQuery : ERROR : %s\n", err)
		return
	}

	log.Printf("RunQuery : %d : Count[%d]\n", query, len(buoyStations))
}

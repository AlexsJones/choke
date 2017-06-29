package mongo

import (
	"fmt"
	"log"
	"time"

	mgo "gopkg.in/mgo.v2"
)

//MongodbConnector ...
type MongodbConnector struct {
	Configuration *MongodbConfiguration
	Session       *mgo.Session
}

//MongodbConfiguration ...
type MongodbConfiguration struct {
	Hosts    []string
	Timeout  time.Duration
	Database string
	Username string
	Password string
}

//NewMongodbConfiguration creates a new configuration
func NewMongodbConfiguration(callback func(*MongodbConfiguration)) *MongodbConfiguration {
	m := &MongodbConfiguration{
		Timeout: 60 * time.Second,
	}
	callback(m)
	return m
}

//NewMongodbConnector creates a new mongodb connector
func NewMongodbConnector(configuration *MongodbConfiguration) *MongodbConnector {

	return &MongodbConnector{Configuration: configuration}
}

//Connect to interface
func (mongo *MongodbConnector) Connect() error {
	fmt.Println("Connecting to mongo")

	mongoDBDialInfo := &mgo.DialInfo{
		Addrs:    mongo.Configuration.Hosts,
		Timeout:  mongo.Configuration.Timeout,
		Database: mongo.Configuration.Database,
		Username: mongo.Configuration.Username,
		Password: mongo.Configuration.Password,
	}

	// Create a session which maintains a pool of socket connections
	// to our MongoDB.
	mongoSession, err := mgo.DialWithInfo(mongoDBDialInfo)
	if err != nil {
		log.Fatalf("CreateSession: %s\n", err)
		return err
	}
	mongo.Session = mongoSession
	mongoSession.SetMode(mgo.Monotonic, true)
	return nil
}

//Disconnect to interface
func (mongo *MongodbConnector) Disconnect() error {
	fmt.Println("Disconnecting from mongo")

	return nil
}

//Request Handle
func (mongo *MongodbConnector) Request() error {

	return nil
}

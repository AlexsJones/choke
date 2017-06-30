package database

//Connect interface to db connector
func Connect(connector Interface) error {

  return connector.Connect()
}
//Disconnect interface to db connector
func Disconnect(connector Interface) error {

  return connector.Disconnect()
}
//Request interface to db connector
func Request(connector Interface) error {

  return connector.Request()
}

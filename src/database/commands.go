package database

//Connect interface to db connector
func Connect(connector Interface) error {

  return connector.Connect()
}

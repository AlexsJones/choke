package database

//Interface is the baseclass for connections
type Interface interface {
	Connect() error
	Disconnect() error
	Request() error
}

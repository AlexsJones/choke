package main

import (
	"github.com/AlexsJones/choke/src/database"
	"github.com/AlexsJones/choke/src/database/mongo"
	"github.com/AlexsJones/choke/src/queue"
	"github.com/AlexsJones/choke/src/routes"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/view"
)

/************************************/
var configuration = mongo.NewMongodbConfiguration(func(config *mongo.MongodbConfiguration) {
	config.Hosts = []string{"localhost"}
})
var m = mongo.NewMongodbConnector(configuration)

var q = queue.NewQueue()

/************************************/
func main() {

	/*********/
	//Connect persistant connector to mongo
	database.Connect(m)
	/********/

	app := iris.New()

	app.OnErrorCode(iris.StatusInternalServerError, func(ctx context.Context) {
		errMessage := ctx.Values().GetString("error")
		if errMessage != "" {
			ctx.Writef("Internal server nerror: %s", errMessage)
			return
		}
		ctx.Writef("(Unexpected) internal server error")
	})

	app.AttachView(view.HTML("./views", ".html").Reload(true))

	routes.AddRoutes(app, m, q)

	app.Get("/data", func(ctx context.Context) {

	})

	go q.Run()

	app.Run(iris.Addr(":8080"), iris.WithCharset("UTF-8"))
}

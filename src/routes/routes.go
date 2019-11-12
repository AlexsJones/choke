package routes

import (
	"github.com/AlexsJones/choke/src/database"
	"github.com/AlexsJones/choke/src/queue"

	"github.com/kataras/iris/v12"
)

//AddRoutes application routes
func AddRoutes(app *iris.Application, databaseContext database.Interface, q *queue.Queue) {

	app.Handle("GET", "/", func(ctx iris.Context) {
		ctx.View("index.html")
	})

	app.Handle("GET", "/healthcheck", func(ctx iris.Context) {
		ctx.WriteString("Okay")
	})

	AddTeaRoutes(app, databaseContext, q)
}

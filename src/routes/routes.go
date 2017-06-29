package routes

import (
	"github.com/AlexsJones/choke/src/database"
	"github.com/AlexsJones/choke/src/queue"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
)

//AddRoutes application routes
func AddRoutes(app *iris.Application, databaseContext database.Interface, q *queue.Queue) {

	app.Handle("GET", "/", func(ctx context.Context) {
		ctx.View("index.html")
	})

	app.Handle("GET", "/healthcheck", func(ctx context.Context) {
		ctx.Write([]byte("Okay"))
	})

	AddTeaRoutes(app, databaseContext, q)
}

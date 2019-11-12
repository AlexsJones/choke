package routes

import (
	"fmt"
	"log"

	"github.com/AlexsJones/choke/src/database"
	"github.com/AlexsJones/choke/src/database/mongo"
	"github.com/AlexsJones/choke/src/database/types"
	"github.com/AlexsJones/choke/src/queue"

	"github.com/kataras/iris/v12"
	"gopkg.in/mgo.v2/bson"
)

//AddTeaRoutes ...
func AddTeaRoutes(app *iris.Application, db database.Interface, q *queue.Queue) {
	mongo, ok := db.(*mongo.MongodbConnector)
	if !ok {
		panic(mongo)
	}
	teaMiddle := func(ctx iris.Context) {
		println(ctx.Method() + ": " + ctx.Path())
		ctx.Next()
	}
	teaRoutes := app.Party("/teas", teaMiddle)
	{
		teaRoutes.Get("/", func(ctx iris.Context) {
			sessionCopy := mongo.Session.Copy()
			defer sessionCopy.Close()
			var teas []types.Tea
			err := mongo.Session.DB("development").C("teas").Find(nil).All(&teas)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			ctx.JSON(teas)
		})
		teaRoutes.Get("/{id:string}", func(ctx iris.Context) {

			id := ctx.Params().Get("id")
			sessionCopy := mongo.Session.Copy()
			defer sessionCopy.Close()

			var tea types.Tea
			sessionCopy.DB("development").C("teas").FindId(bson.ObjectIdHex(id)).One(&tea)

			ctx.Header("Content-Type", "application/vnd.api+json")
			ctx.JSON(tea.Name)
		})
		teaRoutes.Post("/create", func(ctx iris.Context) {

			request := queue.Request{}
			request.Action = func() {
				log.Printf("Writing new request to db....\n")
				sessionCopy := mongo.Session.Copy()
				defer sessionCopy.Close()
				var tea types.Tea
				ctx.ReadJSON(&tea)
				id := bson.NewObjectId()
				sessionCopy.DB("development").C("teas").UpsertId(id, tea)
			}
			go q.PushRequest(request)
			ctx.WriteString("OK")
		})
	}
}

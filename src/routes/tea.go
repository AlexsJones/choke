package routes

import (
	"fmt"

	"gopkg.in/mgo.v2/bson"

	"github.com/AlexsJones/choke/src/database"
	"github.com/AlexsJones/choke/src/database/mongo"
	"github.com/AlexsJones/choke/src/database/types"
	"github.com/AlexsJones/choke/src/queue"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
)

//AddTeaRoutes ...
func AddTeaRoutes(app *iris.Application, databaseContext database.Interface, q *queue.Queue) {

	mongo, ok := databaseContext.(*mongo.MongodbConnector)
	if !ok {
		panic(mongo)
	}
	teaMiddle := func(ctx context.Context) {
		println(ctx.Method() + ": " + ctx.Path())
		ctx.Next()
	}
	teaRoutes := app.Party("/teas", teaMiddle)
	{
		teaRoutes.Get("/", func(ctx context.Context) {
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
		teaRoutes.Get("/{id:string}", func(ctx context.Context) {
			id := ctx.Params().Get("id")
			sessionCopy := mongo.Session.Copy()
			defer sessionCopy.Close()
			collection := &types.TeaRepo{Coll: sessionCopy.DB("development").C("teas")}
			resp, err := collection.Find(id)
			if err != nil {
				ctx.Write([]byte("{}"))
				return
			}
			ctx.Header("Content-Type", "application/vnd.api+json")
			ctx.JSON(resp)
		})
		teaRoutes.Post("/create", func(ctx context.Context) {

			request := queue.Request{}
			request.Action = func() {
				sessionCopy := mongo.Session.Copy()
				defer sessionCopy.Close()
				var tea types.Tea
				ctx.ReadJSON(&tea)
				id := bson.NewObjectId()
				sessionCopy.DB("development").C("teas").UpsertId(id, tea)
			}
			go q.PushRequest(request)
		})
	}
}

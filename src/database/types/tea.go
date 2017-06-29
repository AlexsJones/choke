package types

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

//Tea ...
type Tea struct {
	Id       bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Name     string        `json:"name"`
	Category string        `json:"category"`
}
type TeaRepo struct {
	Coll *mgo.Collection
}
type TeaResource struct {
	Data Tea `json:"data"`
}

func (r *TeaRepo) Find(id string) (TeaResource, error) {
	result := TeaResource{}
	err := r.Coll.FindId(bson.ObjectIdHex(id)).One(&result.Data)
	if err != nil {
		return result, err
	}

	return result, nil
}

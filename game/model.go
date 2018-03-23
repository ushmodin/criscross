package criscross

import "gopkg.in/mgo.v2/bson"

const (
	PlayerUnknown = 0
	PlayerOwner   = 1
	PlayerGuest   = 2

	GameStatusNew    = 0
	GameStatusEnd    = 1
	GameStatusInGame = 2
)

type Game struct {
	ID      bson.ObjectId `bson:"_id,omitempty"`
	Owner   bson.ObjectId `bson:"ownerId"`
	Guest   bson.ObjectId `bson:"guestId"`
	Board   [][]int       `bson:"board"`
	WhoNext int           `bson:"whoNext"`
	Winner  int           `bson:"winner"`
	Status  int           `bson:"status"`
}

type User struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	Username string        `json:"username"`
	Password string        `json:"password"`
	Email    string        `json:"email"`
}

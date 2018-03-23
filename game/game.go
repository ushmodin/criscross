package criscross

import (
	"gopkg.in/mgo.v2/bson"
)

func StartGame(user User) (bson.ObjectId, error) {
	game := Game{
		Owner:   user.ID,
		Board:   [][]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
		WhoNext: PlayerOwner,
		Winner:  PlayerUnknown,
		Status:  GameStatusNew,
	}
	return SaveGame(game)
}

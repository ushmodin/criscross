package criscross

import (
	"crypto/md5"
	"fmt"
	"log"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var mongo *mgo.Session

func StorageConnect(str string) error {
	for {
		conn, err := mgo.Dial(str)
		if err == nil {
			log.Println("Successfully connected to mongodb")
			mongo = conn
			return nil
		}
		log.Printf("Can't connect to mongodb. Reason: %s", err)
		time.Sleep(5 * time.Second)
	}
}

func StorageClose() {
	mongo.Close()
}

func FindUser(username string) (*mgo.Query, error) {
	c := mongo.DB("criscrossgame").C("users")
	return c.Find(bson.M{"username": username}), nil
}

func CreateUser(username, password, email string) error {
	c := mongo.DB("criscrossgame").C("users")
	return c.Insert(User{username, fmt.Sprintf("%x", md5.Sum([]byte(password))), email})
}
func LoadGame(gameId bson.ObjectId) (*mgo.Query, error) {
	c := mongo.DB("criscrossgame").C("games")
	return c.Find(bson.M{"_id": gameId})
}

func SaveGame(game CrisCrossGame) error {
	id := mgo.NewObjectId()
	game.ID = id
	c := mongo.DB("criscrossgame").C("games")
	c.Insert(game)
}

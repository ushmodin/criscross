package criscross

import (
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

func FindUser(username string) *mgo.Query {
	c := mongo.DB("criscrossgame").C("users")
	return c.Find(bson.M{"username": username})
}

func FindUserByID(id bson.ObjectId) *mgo.Query {
	c := mongo.DB("criscrossgame").C("users")
	return c.Find(bson.M{"_id": id})
}

func CreateUser(username string, password string, email string) error {
	c := mongo.DB("criscrossgame").C("users")
	return c.Insert(User{
		Username: username,
		Password: password,
		Email:    email,
	})
}

func LoadGame(gameId bson.ObjectId) *mgo.Query {
	c := mongo.DB("criscrossgame").C("games")
	return c.Find(bson.M{"_id": gameId})
}

func SaveGame(game Game) (bson.ObjectId, error) {
	id := bson.NewObjectId()
	game.ID = id
	c := mongo.DB("criscrossgame").C("games")
	err := c.Insert(game)
	if err != nil {
		return id, err
	}
	return id, nil
}

func UpdateGame(game Game) error {
	c := mongo.DB("criscrossgame").C("games")
	c.Update
}

package criscross

import (
	"crypto/md5"
	"errors"
	"fmt"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type CrisCrossGame struct {
	mgoSession *mgo.Session
}

func NewCrisCrossGame(mongo string) (*CrisCrossGame, error) {
	mgoSession, err := mgo.Dial(mongo)
	if err != nil {
		return nil, err
	}
	return &CrisCrossGame{mgoSession: mgoSession}, nil
}

func (game *CrisCrossGame) Close() {
	game.mgoSession.Close()
}

func (game *CrisCrossGame) regUser(username, password, email string) error {
	c := game.mgoSession.DB("criscrossgame").C("users")
	count, err := c.Find(bson.M{"username": username}).Count()
	if err != nil {
		return errors.New("Unknow error")
	}
	if count > 0 {
		return errors.New("Username alread exists")
	}
	c.Insert(struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}{username, fmt.Sprintf("%x", md5.Sum([]byte(password))), email})
	return nil
}

func (game *CrisCrossGame) auth(username, password string) (string, error) {
	c := game.mgoSession.DB("criscrossgame").C("users")
	q := c.Find(bson.M{"username": username})
	count, err := q.Count()
	if err != nil {
		return "", errors.New("Unknow error")
	}
	if count == 0 {
		return "", errors.New("Unknow user or password")
	}
	var user struct {
		Username string
		Password string
	}
	err = q.One(&user)
	if err != nil {
		return "", errors.New("Unknow error")
	}
	if user.Password != fmt.Sprintf("%x", md5.Sum([]byte(password))) {
		return "", errors.New("Unknow user or password")
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(1)).UnixNano(),
		Issuer:    user.Username,
	})
	tokenStr, err := token.SignedString([]byte("12345678"))
	if err != nil {
		log.Println(err)
		return "", errors.New("Unknow error")
	}
	return tokenStr, nil
}

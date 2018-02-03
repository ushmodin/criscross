package criscross

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
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

func (game *CrisCrossGame) regUser(r *http.Request) error {
	var regReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	err := json.NewDecoder(r.Body).Decode(&regReq)
	if err != nil {
		return errors.New("Invalid request body")
	}

	c := game.mgoSession.DB("criscrossgame").C("users")
	count, err := c.Find(bson.M{"username": regReq.Username}).Count()
	if err != nil {
		return errors.New("Unknow error")
	}
	if count > 0 {
		return errors.New("Username alread exists")
	}
	regReq.Password = fmt.Sprintf("%x", md5.Sum([]byte(regReq.Password)))
	c.Insert(regReq)
	return nil
}

func (game *CrisCrossGame) auth(r *http.Request) (string, error) {
	var authReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&authReq)
	if err != nil {
		return "", errors.New("Invalid request body")
	}

	c := game.mgoSession.DB("criscrossgame").C("users")
	q := c.Find(bson.M{"username": authReq.Username})
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
	if user.Password != fmt.Sprintf("%x", md5.Sum([]byte(authReq.Password))) {
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

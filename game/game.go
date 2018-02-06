package criscross

import (
	"crypto/md5"
	"errors"
	"fmt"
)

type CrisCrossGame struct {
}

func NewCrisCrossGame() (*CrisCrossGame, error) {
	return &CrisCrossGame{}, nil
}

func (game *CrisCrossGame) regUser(username, password, email string) error {
	q, err := FindUser(username)
	if err != nil {
		return errors.New("Unknow error")
	}
	count, err := q.Count()
	if err != nil {
		return errors.New("Unknow error")
	}
	if count > 0 {
		return errors.New("Username alread exists")
	}
	err = CreateUser(username, password, email)
	if err != nil {
		return errors.New("Unknow error")
	}
	return nil
}

func (game *CrisCrossGame) auth(username, password string) error {
	q, err := FindUser(username)
	if err != nil {
		return errors.New("Unknow error")
	}
	count, err := q.Count()
	if err != nil {
		return errors.New("Unknow error")
	}
	if count == 0 {
		return errors.New("Unknow user or password")
	}
	var user struct {
		Password string
	}
	err = q.One(&user)
	if err != nil {
		return errors.New("Unknow error")
	}
	if user.Password != fmt.Sprintf("%x", md5.Sum([]byte(password))) {
		return errors.New("Unknow user or password")
	}
	return nil
}

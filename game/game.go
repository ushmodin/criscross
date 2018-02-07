package criscross

import (
	"crypto/md5"
	"fmt"
)

type CrisCrossGame struct {
}

func NewCrisCrossGame() (*CrisCrossGame, error) {
	return &CrisCrossGame{}, nil
}

func (game *CrisCrossGame) regUser(username, password, email string) gameError {
	q, err := FindUser(username)
	if err != nil {
		return NewGameError(UNKNOW_ERROR, "User not found")
	}
	count, err := q.Count()
	if err != nil {
		return NewGameError(UNKNOW_ERROR, "User not found")
	}
	if count > 0 {
		return NewGameError(REG_ERROR, "Username already exists")
	}
	err = CreateUser(username, password, email)
	if err != nil {
		return NewGameError(UNKNOW_ERROR, "Can't create user")
	}
	return nil
}

func (game *CrisCrossGame) auth(username, password string) gameError {
	q, err := FindUser(username)
	if err != nil {
		return NewGameError(UNKNOW_ERROR, "User not found")
	}
	count, err := q.Count()
	if err != nil {
		return NewGameError(UNKNOW_ERROR, "User not found")
	}
	if count == 0 {
		return NewGameError(AUTH_ERROR, "Unknown user or password1")
	}
	var user struct {
		Password string
	}
	err = q.One(&user)
	if err != nil {
		return NewGameError(UNKNOW_ERROR, "Unknown user or password2")
	}
	if user.Password != fmt.Sprintf("%x", md5.Sum([]byte(password))) {
		return NewGameError(AUTH_ERROR, "Unknown user or password3")
	}
	return nil
}

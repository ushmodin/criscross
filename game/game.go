package criscross

import (
	"crypto/md5"
	"fmt"
	"sync"
)

const (
	NextStepOwner = 1
	NextStepGuest = 2
)

type CrisCrossGame struct {
	sync.RWMutex
	ownerIn, ownerOut chan []byte
	guestIn, guestOut chan []byte
	board             [][]int
	whoNext           int
}

func NewCrisCrossGame(in, out chan []byte) *CrisCrossGame {
	board := [][]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}
	return &CrisCrossGame{ownerIn: in, ownerOut: out, board: board, whoNext: NextStepOwner}
}

func (game *CrisCrossGame) Join(in, out chan []byte) {
	game.guestIn = in
	game.guestOut = out
}

func (game *CrisCrossGame) NextStep(who, row, col int) error {
	game.Lock()
	defer game.Unlock()
	if game.whoNext != who {
		return NewGameError(NOT_YOUR_STEP, "Not your step")
	}
	if row > 3 || row < 0 {
		return NewGameError(VALUE_ERROR, "Invalid row")
	}
	if col > 3 || col < 0 {
		return NewGameError(VALUE_ERROR, "Invalid column")
	}
	if game.board[row][col] != 0 {
		return NewGameError(VALUE_ERROR, "Cel is busy")
	}
	game.board[row][col] = game.whoNext
	if game.whoNext == NextStepGuest {
		game.whoNext = NextStepOwner
	} else if game.whoNext == NextStepOwner {
		game.whoNext = NextStepGuest
	}
	return nil
}

func (game *CrisCrossGame) State() ([][]int, int) {
	game.RLock()
	defer game.RUnlock()
	rows := len(game.board)
	res := make([][]int, rows)
	for r := 0; r < rows; r++ {
		res[r] = make([]int, len(game.board[r]))
		copy(res[r], game.board[r])
	}
	return res, game.whoNext
}

type CrisCrossEngine struct {
	games map[interface{}]CrisCrossGame
}

func NewCrisCrossEngine() (*CrisCrossEngine, error) {
	return &CrisCrossEngine{}, nil
}

func (game *CrisCrossEngine) start(in, out chan []byte) interface{} {
	return nil
}

func (game *CrisCrossEngine) stop(id interface{}) {

}

func (game *CrisCrossEngine) join(in, out chan []byte) {

}

func (game *CrisCrossEngine) regUser(username, password, email string) gameError {
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

func (game *CrisCrossEngine) auth(username, password string) gameError {
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

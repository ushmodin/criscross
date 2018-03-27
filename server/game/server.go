package criscross

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var jwtKey = []byte("12345678")
var keyGetterFunc = func(token *jwt.Token) (interface{}, error) {
	return jwtKey, nil
}
var jwtSignMethod = jwt.SigningMethodHS256
var upgrader = websocket.Upgrader{
	ReadBufferSize:  10,
	WriteBufferSize: 10,
}

func StartHttpServer(addr string) error {
	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: keyGetterFunc,
		SigningMethod:       jwtSignMethod,
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err string) {
			writeError(w, NewGameError(AUTH_ERROR, err))
		},
	})

	r := mux.NewRouter()
	r.HandleFunc("/api/ping", pingHandler).Methods("GET")
	r.HandleFunc("/api/reg", regHandler).Methods("POST")
	r.HandleFunc("/api/auth", authHandler).Methods("POST")
	r.Handle("/api/game/start", jwtMiddleware.Handler(http.HandlerFunc(startGameHandler))).Methods("GET")
	r.Handle("/api/game/{id}", jwtMiddleware.Handler(http.HandlerFunc(gameStatusHandler))).Methods("GET")
	r.Handle("/api/game/{id}", jwtMiddleware.Handler(http.HandlerFunc(stepHandler))).Methods("POST")

	return http.ListenAndServe(addr, r)
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
	w.WriteHeader(http.StatusOK)
}

func regHandler(w http.ResponseWriter, r *http.Request) {
	var regReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	err := json.NewDecoder(r.Body).Decode(&regReq)
	if err != nil {
		writeError(w, err)
	}

	err = CreateUser(regReq.Username, fmt.Sprintf("%x", md5.Sum([]byte(regReq.Password))), regReq.Email)

	if err != nil {
		writeError(w, err)
	} else {
		w.Write([]byte("{}"))
	}
}
func authHandler(w http.ResponseWriter, r *http.Request) {
	var authReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&authReq)
	if err != nil {
		writeError(w, err)
	}
	q := FindUser(authReq.Username)
	var user User
	err = q.One(&user)
	if err != nil {
		writeError(w, NewGameError(AUTH_ERROR, "User or password not found"))
		return
	}
	if user.Password != fmt.Sprintf("%x", md5.Sum([]byte(authReq.Password))) {
		writeError(w, NewGameError(AUTH_ERROR, "User or password not found"))
		return
	}
	token, err := createToken(user)

	if err != nil {
		writeError(w, err)
	} else {
		w.Header().Add("Authorization", "Bearer "+token)
	}
}

func startGameHandler(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r)
	if err != nil {
		writeError(w, err)
		return
	}
	id, err := StartGame(user)
	if err != nil {
		writeError(w, err)
		return
	}

	rsp := struct {
		GameID string `json:"gameId"`
	}{id.String()}
	json.NewEncoder(w).Encode(rsp)
}

func stepHandler(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r)
	if err != nil {
		writeError(w, err)
		return
	}
	params := mux.Vars(r)
	gameID := params["id"]
	var game Game
	if err := LoadGame(bson.ObjectId(gameID)).One(&game); err != nil {
		writeError(w, NewGameError(GAME_ERROR, "Game not found"))
		return
	}
	if game.Owner != user.ID || game.Guest != user.ID {
		writeError(w, NewGameError(GAME_ERROR, "Game not found"))
		return
	}
	if (game.WhoNext == PlayerGuest && user.ID != game.Guest) ||
		(game.WhoNext == PlayerOwner && user.ID != game.Owner) {
		writeError(w, NewGameError(NOT_YOUR_STEP, ""))
		return
	}
	var req struct {
		Row int `json:"row"`
		Col int `json:"col"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, errors.New("Invalid request"))
		return
	}
	if game.Board[req.Row][req.Col] != PlayerUnknown {
		writeError(w, NewGameError(CELL_BUSY, ""))
		return
	}
	if user.ID == game.Owner {
		game.Board[req.Row][req.Col] = PlayerOwner
	} else {
		game.Board[req.Row][req.Col] = PlayerGuest
	}
	updateWinner(game)
	UpdateGame(game)
}

func gameStatusHandler(w http.ResponseWriter, r *http.Request) {
	user, err := getUser(r)
	if err != nil {
		writeError(w, err)
		return
	}
	params := mux.Vars(r)
	gameID := params["id"]
	var game Game
	if err := LoadGame(bson.ObjectId(gameID)).One(&game); err != nil {
		writeError(w, NewGameError(GAME_ERROR, "Game not found"))
		return
	}
	if game.Owner != user.ID || game.Guest != user.ID {
		writeError(w, NewGameError(GAME_ERROR, "Game not found"))
		return
	}
	var rsp struct {
		Board  [][]string `json:"board"`
		Next   string     `json:"next"`
		Winner string     `json:"winner"`
		Guest  string     `json:"guest"`
		Owner  string     `json:"owner"`
	}
	rsp.Board = make([][]string, len(game.Board))
	for i := 0; i < len(game.Board); i++ {
		rsp.Board[i] = make([]string, 3)
		for j := 0; j < len(game.Board[i]); j++ {
			if game.Board[i][j] == PlayerGuest {
				rsp.Board[i][j] = "O"
			} else if game.Board[i][j] == PlayerOwner {
				rsp.Board[i][j] = "X"
			} else {
				rsp.Board[i][j] = " "
			}
		}
	}
	if game.WhoNext == PlayerGuest {
		rsp.Next = "GUEST"
	} else if game.WhoNext == PlayerOwner {
		rsp.Next = "OWNER"
	} else {
		rsp.Next = "UNKNOW"
	}
	if game.Status == GameStatusEnd {
		if game.Winner == PlayerGuest {
			rsp.Winner = "GUEST"
		} else if game.Winner == PlayerOwner {
			rsp.Winner = "OWNER"
		} else {
			rsp.Winner = "UNKNOWN"
		}
	}
	var guest User
	if err := FindUserByID(game.Guest).One(&guest); err != nil {
		writeError(w, errors.New("Can't find guest"))
		return
	}
	rsp.Guest = guest.Username
	var owner User
	if err := FindUserByID(game.Owner).One(&owner); err != nil {
		writeError(w, errors.New("Can't find Owner"))
		return
	}
	rsp.Owner = owner.Username
	json.NewEncoder(w).Encode(rsp)
}

func getUser(r *http.Request) (User, error) {
	header := r.Header["Authorization"][0]
	var token string
	if strings.HasPrefix(header, "Bearer ") {
		token = header[0:7]
	}
	if len(token) == 0 {
		return User{}, errors.New("Can't parse jwt token")
	}
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, keyGetterFunc)
	if err != nil {
		return User{}, err
	}
	userID := claims["Subject"].(string)
	if len(userID) == 0 {
		return User{}, errors.New("Invalide token. Can't find Subject")
	}
	var user User
	if err := FindUserByID(bson.ObjectId(userID)).One(&user); err != nil {
		return User{}, errors.New("Invalide token. Can't find User")
	}
	return user, nil
}

func writeError(w http.ResponseWriter, err error) {
	gameError, ok := err.(gameError)
	code := UNKNOW_ERROR
	if ok {
		code = gameError.Code()
	}
	rsp := struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}{code, err.Error()}
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(rsp)
}

func createToken(user User) (string, error) {
	token := jwt.NewWithClaims(jwtSignMethod, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(1)).UnixNano(),
		Subject:   user.ID.String(),
		Issuer:    user.Username,
	})
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

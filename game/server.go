package criscross

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"

	"github.com/gorilla/mux"
)

var jwtKey = []byte("12345678")
var jwtSignMethod = jwt.SigningMethodHS256

type CrisCrossServer struct {
	game *CrisCrossGame
}

func NewCrisCrossServer(game *CrisCrossGame) (*CrisCrossServer, error) {
	return &CrisCrossServer{game: game}, nil
}

func (srv *CrisCrossServer) ListenAndServe(addr string) {
	r := srv.createRouter()
	http.ListenAndServe(addr, r)
}

func (srv *CrisCrossServer) createRouter() *mux.Router {

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		},
		SigningMethod: jwtSignMethod,
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err string) {
			writeError(w, NewGameError(AUTH_ERROR, err))
		},
	})

	r := mux.NewRouter()
	r.HandleFunc("/api/ping", pingHandler).Methods("GET")
	r.HandleFunc("/api/reg", srv.regHandler).Methods("POST")
	r.HandleFunc("/api/auth", srv.authHandler).Methods("POST")
	r.Handle("/api/game/start", jwtMiddleware.Handler(http.HandlerFunc(srv.startGame))).Methods("GET")
	return r
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
	w.WriteHeader(http.StatusOK)
}

func (srv *CrisCrossServer) regHandler(w http.ResponseWriter, r *http.Request) {
	var regReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	err := json.NewDecoder(r.Body).Decode(&regReq)
	if err != nil {
		writeError(w, err)
	}

	err = srv.game.regUser(regReq.Username, regReq.Password, regReq.Email)

	if err != nil {
		writeError(w, err)
	} else {
		w.Write([]byte("{}"))
	}
}
func (srv *CrisCrossServer) authHandler(w http.ResponseWriter, r *http.Request) {
	var authReq struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&authReq)
	if err != nil {
		writeError(w, err)
	}
	err = srv.game.auth(authReq.Username, authReq.Password)
	if err != nil {
		writeError(w, err)
	} else {
		token, err := createToken(authReq.Username)

		if err != nil {
			writeError(w, err)
		} else {
			w.Header().Add("Authorization", "Bearer "+token)
		}
	}
}

func (srv *CrisCrossServer) startGame(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
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

func createToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwtSignMethod, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(1)).UnixNano(),
		Issuer:    username,
	})
	tokenStr, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

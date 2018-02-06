package criscross

import (
	"encoding/json"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"

	"github.com/gorilla/mux"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

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
	r := mux.NewRouter()
	r.HandleFunc("/api/ping", pingHandler).Methods("GET")
	r.HandleFunc("/api/reg", srv.regHandler).Methods("POST")
	r.HandleFunc("/api/auth", srv.authHandler).Methods("POST")
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

func writeError(w http.ResponseWriter, err error) {
	rsp := ErrorResponse{
		Code:    "AUTH_ERROR",
		Message: err.Error(),
	}
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(rsp)
}

func createToken(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * time.Duration(1)).UnixNano(),
		Issuer:    username,
	})
	tokenStr, err := token.SignedString([]byte("12345678"))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

package main

import (
	"log"
	"net/http"

	"github.com/ushmodin/criscross"
)

func main() {
	game, err := criscross.NewCrisCrossGame("localhost")
	if err != nil {
		log.Fatal(err)
	}
	defer game.Close()
	http.ListenAndServe(":8080", game.CreateRouter())
}

package main

import (
	"log"

	"github.com/ushmodin/criscross"
)

func main() {
	game, err := criscross.NewCrisCrossGame("localhost")
	if err != nil {
		log.Fatal(err)
	}
	defer game.Close()
	srv, err := criscross.NewCrisCrossServer(game)
	if err != nil {
		log.Fatal(err)
	}
	srv.ListenAndServe(":8080")
}

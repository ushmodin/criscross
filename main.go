package main

import (
	"fmt"
	"github.com/ushmodin/criscross/game"
	"log"
	"os"
)

func main() {
	game, err := criscross.NewCrisCrossGame(os.Getenv("MONGODB"))
	if err != nil {
		log.Fatal(err)
	}
	defer game.Close()
	srv, err := criscross.NewCrisCrossServer(game)
	if err != nil {
		log.Fatal(err)
	}
	srv.ListenAndServe(":8080")
	fmt.Println("End")
}

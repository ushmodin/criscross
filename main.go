package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ushmodin/criscross/game"
)

func main() {
	err := criscross.StorageConnect(os.Getenv("MONGODB"))
	if err != nil {
		log.Fatal(err)
	}
	game, err := criscross.NewCrisCrossGame()
	if err != nil {
		log.Fatal(err)
	}
	defer criscross.StorageClose()
	srv, err := criscross.NewCrisCrossServer(game)
	if err != nil {
		log.Fatal(err)
	}
	srv.ListenAndServe(":8080")
	fmt.Println("End")
}

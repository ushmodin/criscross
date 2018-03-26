package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ushmodin/criscross/server/game"
)

func main() {
	if err := criscross.StorageConnect(os.Getenv("MONGODB")); err != nil {
		log.Fatal(err)
	}
	defer criscross.StorageClose()
	err := criscross.StartHttpServer("localhost:4000")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("End")
}

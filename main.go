package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ushmodin/criscross/game"
)

func main() {
	if err := criscross.RedisConnect("localhost:6379", 0); err != nil {
		log.Fatal(err)
	}
	if err := criscross.StorageConnect(os.Getenv("MONGODB")); err != nil {
		log.Fatal(err)
	}
	defer criscross.RedisClose()
	defer criscross.StorageClose()
	err := criscross.StartHttpServer("localhost:4000")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("End")
}

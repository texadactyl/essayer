package main

import (
	"essayer/bridges"
	"log"
)

func main() {
	bridges.Setup()
	tryCrcMain()
	log.Print("main: Bye-bye")
	tryFlateMain()
}

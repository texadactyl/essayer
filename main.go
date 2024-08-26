package main

import (
	"essayer/bridges"
	"log"
)

func main() {
	bridges.Setup()
	tryZip()
	log.Print("main: Bye-bye")
}

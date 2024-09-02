package main

import (
	"essayer/bridges"
	"log"
)

func threadedFunc(ch chan bool) {
	log.Print("threadedFunc: Calling tryCrcMain")
	tryCrcMain()
	log.Print("threadedFunc: Calling tryFlateMain")
	tryFlateMain()
	log.Print("threadedFunc: End, will signal main on the channel then return")
	ch <- true
}

func main() {
	ch := make(chan bool)
	bridges.Setup()
	go func() {
		threadedFunc(ch)
	}()
	<-ch
	log.Print("main: Returned from threadedFunc")
	log.Print("main: Bye-bye")
}

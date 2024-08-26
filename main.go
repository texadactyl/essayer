package main

import (
	"log"
	"main/bridges"
)

func main() {
	bridges.Setup()
	zipLibPath := bridges.DirLibs + bridges.PathStringSep + "libzip." + bridges.LibExt
	tryZip(zipLibPath)
	log.Print("main: Bye-bye")
}

package main

import (
	"essayer/bridges"
	"log"
)

func main() {
	bridges.Setup()
	var zipLibPath string
	if bridges.WindowsOS {
		zipLibPath = bridges.DirLibs + bridges.PathStringSep + "zip." + bridges.LibExt
	} else {
		zipLibPath = bridges.DirLibs + bridges.PathStringSep + "libzip." + bridges.LibExt

	}
	tryZip(zipLibPath)
	log.Print("main: Bye-bye")
}

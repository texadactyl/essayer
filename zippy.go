package main

import (
	"essayer/bridges"
	"log"
)

func tryZip() {

	log.Println("tryZip: Begin")

	var err error
	var zipLibPath string

	// Form the zip library path.
	if bridges.WindowsOS {
		zipLibPath = bridges.DirLibs + bridges.PathStringSep + "zip." + bridges.LibExt
	} else {
		zipLibPath = bridges.DirLibs + bridges.PathStringSep + "libzip." + bridges.LibExt

	}

	// Open the zip library.
	_ = bridges.ConnectLibrary(zipLibPath)
	if err != nil {
		log.Fatalf("tryZip: purego.Dlopen for [%s] failed, reason: [%s]\n", zipLibPath, err.Error())
	}
	log.Printf("tryZip: library connected for [%s] ok\n", zipLibPath)

	log.Println("tryZip: End")

}

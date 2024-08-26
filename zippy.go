package main

import (
	"log"
)

func tryZip(zipLibPath string) {

	log.Println("tryZip: Begin")

	//var handle uintptr
	var err error

	// Open the zip library.
	_ = connectLibrary(zipLibPath)
	if err != nil {
		log.Fatalf("tryZip: purego.Dlopen for [%s] failed, reason: [%s]\n", zipLibPath, err.Error())
	}
	log.Printf("tryZip: library connected for [%s] ok\n", zipLibPath)

	log.Println("tryZip: End")

}

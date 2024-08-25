package main

import (
	"github.com/ebitengine/purego"
	"log"
)

func tryZip(zipLibPath string) {

	log.Println("tryZip: Begin")

	var zipLibHandle uintptr
	var err error

	// Open the zip library.
	zipLibHandle = connectLibrary(zipLibPath)
	if err != nil {
		log.Fatalf("DoTheZip: purego.Dlopen for [%s] failed, reason: [%s]\n", zipLibPath, err.Error())
	}
	log.Printf("DoTheZip: purego.Dlopen for [%s] ok\n", zipLibPath)

	// Close it.
	err = purego.Dlclose(zipLibHandle)
	if err != nil {
		log.Fatalf("DoTheZip: purego.Dlclose(handle) failed for [%s], reason: [%s]", zipLibPath, err.Error())
	}
	log.Printf("DoTheZip: purego.Dlclose for [%s] ok\n", zipLibPath)

	log.Println("tryZip: End")

}

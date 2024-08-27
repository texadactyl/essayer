/*
	Exercise some of the JVM libzip functions.

	Library source: src/java.base/share/native/libzip/CRC32.c
	Library include file directory: src/java.base/share/native/include/

	On-line checker: https://crc32.online/
*/

package main

import (
	"essayer/bridges"
	"github.com/ebitengine/purego"
	"log"
	"unsafe"
)

func try_ZIP_CRC32(libHandle uintptr) {

	// Register the ZIP_CRC32 library function.
	var crc32UpdateFunc func(uint32, unsafe.Pointer, uint32) uint32
	purego.RegisterLibFunc(&crc32UpdateFunc, libHandle, "ZIP_CRC32")
	log.Println("tryZip: purego.RegisterLibFunc ok")

	// Data.
	observed := uint32(0)                        // initial CRC value
	data := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ") // Argument for CRC32
	datalen := uint32(len(data))
	expected := uint32(0xabf77822)

	// Execute ZIP_CRC32().
	observed = crc32UpdateFunc(observed, unsafe.Pointer(&data[0]), datalen)
	if observed != expected {
		log.Fatalf("tryZip: Oops, expected: 0x%08x, observed: 0x%08x\n", expected, observed)
	} else {
		log.Printf("tryZip: Success, observed = expected = 0x%08x\n", observed)
	}

}

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
	libHandle := bridges.ConnectLibrary(zipLibPath)
	if err != nil {
		log.Fatalf("tryZip: purego.Dlopen for [%s] failed, reason: [%s]\n", zipLibPath, err.Error())
	}
	log.Printf("tryZip: library connected for [%s] ok\n", zipLibPath)

	// Run individual tests.
	try_ZIP_CRC32(libHandle)

	log.Println("tryZip: End")

}

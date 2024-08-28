/*
	Exercise some of the JVM libzip CRC32 functions.

	Library source: src/java.base/share/native/libzip/CRC32.c
	Library include file directory: src/java.base/share/native/include/

	On-line checker for CRC32: https://crc32.online/

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
	funcName := "ZIP_CRC32"
	var crc32UpdateFunc func(uint32, unsafe.Pointer, uint32) uint32
	purego.RegisterLibFunc(&crc32UpdateFunc, libHandle, funcName)
	log.Printf("tryCrcMain: purego.RegisterLibFunc (%s) ok\n", funcName)

	// Data.
	observed := uint32(0)                        // initial CRC value
	data := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ") // Argument for CRC32
	datalen := uint32(len(data))
	expected := uint32(0xabf77822)

	// Execute ZIP_CRC32().
	observed = crc32UpdateFunc(observed, unsafe.Pointer(&data[0]), datalen)
	if observed != expected {
		log.Fatalf("tryCrcMain: Oops, expected: 0x%08x, observed: 0x%08x\n", expected, observed)
	} else {
		log.Printf("tryCrcMain: Success, observed = expected = 0x%08x\n", observed)
	}

}

func try_Java_java_util_zip_CRC32_update(libHandle uintptr) {

	// Register the ZIP_CRC32 library function.
	funcName := "Java_java_util_zip_CRC32_update"
	var crc32UpdateFunc func(unsafe.Pointer, unsafe.Pointer, uint32, uint32) uint32
	purego.RegisterLibFunc(&crc32UpdateFunc, libHandle, funcName)
	log.Printf("tryCrcMain: purego.RegisterLibFunc (%s) ok\n", funcName)

	// Data.
	observed := uint32(0) // initial CRC value
	data := uint32('A')   // Argument for CRC32
	expected := uint32(0xd3d99e8b)
	dummyData := 0
	dummyPtr := unsafe.Pointer(&dummyData)

	// Execute ZIP_CRC32().
	observed = crc32UpdateFunc(dummyPtr, dummyPtr, observed, data)
	if observed != expected {
		log.Fatalf("tryCrcMain: Oops, expected: 0x%08x, observed: 0x%08x\n", expected, observed)
	} else {
		log.Printf("tryCrcMain: Success, observed = expected = 0x%08x\n", observed)
	}

}

func tryCrcMain() {

	log.Println("tryCrcMain: Begin")

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
		log.Fatalf("tryCrcMain: purego.Dlopen for [%s] failed, reason: [%s]\n", zipLibPath, err.Error())
	}
	log.Printf("tryCrcMain: library connected for [%s] ok\n", zipLibPath)

	// Run individual tests.
	try_ZIP_CRC32(libHandle)
	try_Java_java_util_zip_CRC32_update(libHandle)

	log.Println("tryCrcMain: End")

}

/*
	Exercise some of the JVM libzip Deflate/Inflate functions.

	Library source: src/java.base/share/native/libzip/{Deflater.c, Inflater.c}
	Library include file directory: src/java.base/share/native/include/

	Deflate/Inflate parameters from src/java.base/share/native/libzip/zlib/zlib.h:
		Compression Levels:
			#define Z_NO_COMPRESSION         0
			#define Z_BEST_SPEED             1
			#define Z_BEST_COMPRESSION       9
			#define Z_DEFAULT_COMPRESSION  (-1)
		Compression strategy:
			#define Z_FILTERED            1
			#define Z_HUFFMAN_ONLY        2
			#define Z_RLE                 3
			#define Z_FIXED               4
			#define Z_DEFAULT_STRATEGY    0

*/

package main

import (
	"essayer/bridges"
	"fmt"
	"github.com/ebitengine/purego"
	"log"
	"unsafe"
)

func tryDeflater(libHandle uintptr, inBytes []byte) []byte {

	// Register the Java_java_util_zip_Deflater_init library function.
	funcName := "Java_java_util_zip_Deflater_init"
	var deflaterInit func(unsafe.Pointer, unsafe.Pointer, uint32, uint32, uint8) uint64
	purego.RegisterLibFunc(&deflaterInit, libHandle, funcName)
	log.Printf("tryDeflater: purego.RegisterLibFunc (%s) ok\n", funcName)

	// Data.
	level := uint32(9)    // best compression
	strategy := uint32(0) // default strategy
	noWrap := uint8(0)    // false
	dummyData := 0
	dummyPtr := unsafe.Pointer(&dummyData)

	// Initialise deflater.
	streamPtr := deflaterInit(dummyPtr, dummyPtr, level, strategy, noWrap)
	if streamPtr == uint64(0) {
		log.Fatalln("tryDeflater: Oops, deflaterInit failed")
	} else {
		log.Println("tryDeflater: deflaterInit ok")
	}

	// Return output.
	return inBytes

}

func tryZipMain() {

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

	// Run individual CRC tests.
	try_ZIP_CRC32(libHandle)
	try_Java_java_util_zip_CRC32_update(libHandle)

	// Deflater test.
	rawData := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ") // Inflated raw data
	rawDatalen := uint32(len(rawData))
	crunchedData := make([]byte, rawDatalen)

	crunchedData = tryDeflater(libHandle, rawData)
	fmt.Printf("DEBUG tryDeflater returned %s\n", string(crunchedData))

	log.Println("tryCrcMain: End")

}

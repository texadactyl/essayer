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

	// JNIEXPORT jint ZIP_CRC32(jint crc, const jbyte *buf, jint len)
	funcName := "ZIP_CRC32"
	var crc32UpdateFunc func(uint32, unsafe.Pointer, uint32) uint32

	// Register the ZIP_CRC32 library function.
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

	// JNIEXPORT jint JNICALL Java_java_util_zip_CRC32_update(JNIEnv *env, jclass cls, jint crc, jint b)
	funcName := "Java_java_util_zip_CRC32_update"
	var crc32UpdateFunc func(uintptr, unsafe.Pointer, uint32, uint32) uint32

	// Register the ZIP_CRC32 library function.
	purego.RegisterLibFunc(&crc32UpdateFunc, libHandle, funcName)
	log.Printf("tryCrcMain: purego.RegisterLibFunc (%s) ok\n", funcName)

	// Data.
	observed := uint32(0) // initial CRC value
	data := uint32('A')   // Argument for CRC32
	expected := uint32(0xd3d99e8b)
	dummyData := 0
	dummyPtr := unsafe.Pointer(&dummyData)

	// Execute ZIP_CRC32().
	observed = crc32UpdateFunc(bridges.HandleENV, dummyPtr, observed, data)
	if observed != expected {
		log.Fatalf("tryCrcMain: Oops, expected: 0x%08x, observed: 0x%08x\n", expected, observed)
	} else {
		log.Printf("tryCrcMain: Success, observed = expected = 0x%08x\n", observed)
	}

}

func tryCrcMain() {

	log.Println("tryCrcMain: Begin")

	var err error
	var pathZip string

	// Form the zip library path.
	if bridges.WindowsOS {
		pathZip = bridges.PathDirLibs + bridges.SepPathString + "zip." + bridges.FileExt
	} else {
		pathZip = bridges.PathDirLibs + bridges.SepPathString + "libzip." + bridges.FileExt
	}

	// Open the zip library.
	handleZip := bridges.ConnectLibrary(pathZip)
	if err != nil {
		log.Fatalf("tryCrcMain: purego.Dlopen for [%s] failed, reason: [%s]\n", pathZip, err.Error())
	}
	log.Printf("tryCrcMain: library connected for [%s] ok\n", pathZip)

	// Run individual tests.
	try_ZIP_CRC32(handleZip)
	try_Java_java_util_zip_CRC32_update(handleZip)

	log.Println("tryCrcMain: End")

}

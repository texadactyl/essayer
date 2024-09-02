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

	Allowed flush values
			#define Z_NO_FLUSH      0
			#define Z_PARTIAL_FLUSH 1
			#define Z_SYNC_FLUSH    2
			#define Z_FULL_FLUSH    3
			#define Z_FINISH        4
			#define Z_BLOCK         5
			#define Z_TREES         6
*/

package main

import (
	"essayer/bridges"
	"github.com/ebitengine/purego"
	"log"
)

func tryDeflater(libHandle uintptr, inBytes []byte) []byte {

	var streamAddr bridges.JNIlong
	tempBytes := make([]byte, len(inBytes))
	dummyClass := make([]byte, 3)
	dummyClassPtr := &dummyClass

	// -------------------------------------------------------------------------------------------
	// JNIEXPORT jlong JNICALL
	//Java_java_util_zip_Deflater_init(JNIEnv *env, jclass cls, jint level,
	//                                 jint strategy, jboolean nowrap)
	funcName := "Java_java_util_zip_Deflater_init"
	var deflaterInit func(uintptr, bridges.JNIobject, bridges.JNIint, bridges.JNIint, bridges.JNIboolean) bridges.JNIlong

	// Register the Java_java_util_zip_Deflater_init library function.
	purego.RegisterLibFunc(&deflaterInit, libHandle, funcName)
	log.Printf("tryDeflater: purego.RegisterLibFunc (%s) ok\n", funcName)

	// Data for deflaterInit.
	level := bridges.JNIint(9)      // best compression
	strategy := bridges.JNIint(0)   // default strategy
	noWrap := bridges.JNIboolean(0) // false
	flush := bridges.JNIint(3)      // full flush
	params := bridges.JNIint(0)     // no change to parameters

	// Call deflaterInit.
	streamAddr = deflaterInit(bridges.HandleENV, bridges.JNIobject(dummyClassPtr), level, strategy, noWrap)
	if streamAddr == bridges.JNIlong(0) {
		log.Fatalf("tryDeflater: Oops, %s failed\n", funcName)
	} else {
		log.Printf("tryDeflater: %s ok\n", funcName)
	}

	// -------------------------------------------------------------------------------------------
	// JNIEXPORT jlong JNICALL
	// Java_java_util_zip_Deflater_deflateBytesBytes(JNIEnv *env, jobject this, jlong addr,
	//                                         jbyteArray inputArray, jint inputOff, jint inputLen,
	//                                         jbyteArray outputArray, jint outputOff, jint outputLen,
	//                                         jint flush, jint params)
	funcName = "Java_java_util_zip_Deflater_deflateBytesBytes"
	var deflaterBytesToBytes func(uintptr, bridges.JNIobject, bridges.JNIlong,
		bridges.JNIbyteArray, bridges.JNIint, bridges.JNIint,
		bridges.JNIbyteArray, bridges.JNIint, bridges.JNIint,
		bridges.JNIint, bridges.JNIint) bridges.JNIlong

	// Register the Java_java_util_zip_Deflater_end library function.
	purego.RegisterLibFunc(&deflaterBytesToBytes, libHandle, funcName)
	log.Printf("tryDeflater: purego.RegisterLibFunc (%s) ok\n", funcName)

	// Call deflaterEnd.
	sizeDeflated := deflaterBytesToBytes(bridges.HandleENV, bridges.JNIobject(dummyClassPtr), streamAddr,
		bridges.JNIbyteArray(&inBytes[0]), bridges.JNIint(0), bridges.JNIint(len(inBytes)),
		bridges.JNIbyteArray(&tempBytes[0]), bridges.JNIint(0), bridges.JNIint(len(tempBytes)),
		flush, params)
	if sizeDeflated < 1 {
		log.Fatalf("tryDeflater: Oops, %s failed\n", funcName)
	} else {
		log.Printf("tryDeflater: %s ok\n", funcName)
	}

	// -------------------------------------------------------------------------------------------
	// JNIEXPORT void JNICALL Java_java_util_zip_Deflater_end(JNIEnv *env, jclass cls, jlong addr)
	funcName = "Java_java_util_zip_Deflater_end"
	var deflaterEnd func(uintptr, bridges.JNIobject, bridges.JNIlong)

	// Register the Java_java_util_zip_Deflater_end library function.
	purego.RegisterLibFunc(&deflaterEnd, libHandle, funcName)
	log.Printf("tryDeflater: purego.RegisterLibFunc (%s) ok\n", funcName)

	// Call deflaterEnd.
	deflaterEnd(bridges.HandleENV, bridges.JNIobject(dummyClassPtr), streamAddr)
	log.Printf("tryDeflater: %s ok\n", funcName)

	// -------------------------------------------------------------------------------------------
	// Return output.
	outBytes := make([]byte, sizeDeflated)
	_ = copy(outBytes, tempBytes[:sizeDeflated])
	return outBytes

}

func tryFlateMain() {

	log.Println("tryFlateMain: Begin")

	var pathZip string

	// Form the zip library path.
	if bridges.WindowsOS {
		pathZip = bridges.PathDirLibs + bridges.SepPathString + "zip." + bridges.FileExt
	} else {
		pathZip = bridges.PathDirLibs + bridges.SepPathString + "libzip." + bridges.FileExt
	}

	// Open the zip library.
	handleZip := bridges.ConnectLibrary(pathZip)
	log.Printf("tryFlateMain: library connected for [%s] ok\n", pathZip)

	// Deflater test.
	inflatedData := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ") // Inflated raw data
	log.Printf("tryFlateMain: Calling tryDeflater with %d inflated bytes\n", len(inflatedData))
	deflatedData := tryDeflater(handleZip, inflatedData)
	log.Printf("tryFlateMain: tryDeflater returned %d deflated bytes\n", len(deflatedData))

	log.Println("tryFlateMain: End")

}

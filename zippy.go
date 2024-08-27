/*
	Library source: src/java.base/share/native/libzip/CRC32.c
	Include file directory: src/java.base/share/native/include/

	Chaecker: https://crc32.online/
*/

package main

import (
	"essayer/bridges"
	"github.com/ebitengine/purego"
	"log"
	"unsafe"
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
	libHandle := bridges.ConnectLibrary(zipLibPath)
	if err != nil {
		log.Fatalf("tryZip: purego.Dlopen for [%s] failed, reason: [%s]\n", zipLibPath, err.Error())
	}
	log.Printf("tryZip: library connected for [%s] ok\n", zipLibPath)

	/*
		Java_java_util_zip_CRC32_update(JNIEnv *env, jclass cls, jint crc, jint b)
		{
			Bytef buf[1];

			buf[0] = (Bytef)b;
			return crc32(crc, buf, 1);
		}
	*/

	// Register the library function.
	var crc32UpdateFunc func(uint32, unsafe.Pointer, uint32) uint32
	purego.RegisterLibFunc(&crc32UpdateFunc, libHandle, "ZIP_CRC32")
	log.Println("tryZip: purego.RegisterLibFunc ok")

	// Data.
	crc := uint32(0)                             // initial CRC value
	data := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ") // Argument for CRC32
	datalen := uint32(len(data))

	// Execute ZIP_CRC32().
	crc = crc32UpdateFunc(crc, unsafe.Pointer(&data[0]), datalen)
	log.Printf("tryZip: CRC32 = 0x%08x\n", crc)

	log.Println("tryZip: End")

}

package main

import (
	"log"
	"os"
	"runtime"
)

var OperSys string
var WindowsOS = false
var DirLibs string
var LibJvm string
var LibExt string
var PathStringSep = string(os.PathSeparator)

func setup() {

	log.Println("setup: Begin")

	// Set up library file extension and library path string as a function of O/S.
	OperSys = runtime.GOOS
	switch OperSys {
	case "darwin":
		LibExt = "dylib"
	case "linux", "freebsd":
		LibExt = "so"
	case "windows":
		LibExt = "dll"
		WindowsOS = true
	default:
		log.Fatalln("setup: Unsupported O/S: %s", OperSys)
	}

	// Get Java home.
	javaHome := os.Getenv("JAVA_HOME")
	if javaHome == "" {
		log.Fatalln("setup: Environment variable JAVA_HOME missing but is required. Exiting.")
	}

	// Calculate the path of the Java lib directory.
	DirLibs = javaHome + PathStringSep + "lib"

	// Calculate the path of the server libjvm and open it.
	LibJvm = DirLibs + PathStringSep + "server" + PathStringSep + "libjvm." + LibExt
	_ = connectLibrary(LibJvm)
	log.Println("setup: End")

}

func main() {
	setup()
	zipLibPath := DirLibs + PathStringSep + "libzip." + LibExt
	tryZip(zipLibPath)
	log.Print("main: Bye-bye")
}

package bridges

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

func Setup() {

	log.Println("bridges/Setup: Begin")

	// Set up library file extension and library path string as a function of O/S.
	OperSys = runtime.GOOS
	switch OperSys {
	case "darwin":
		LibExt = "dylib"
	case "linux":
		LibExt = "so"
	case "windows":
		LibExt = "dll"
		WindowsOS = true
	default:
		log.Fatalln("bridges/Setup: Unsupported O/S: %s", OperSys)
	}

	// Get Java home.
	javaHome := os.Getenv("JAVA_HOME")
	if javaHome == "" {
		log.Fatalln("bridges/Setup: Environment variable JAVA_HOME missing but is required. Exiting.")
	}

	// Calculate the path of the Java lib directory.
	// Calculate the path of the server JVM library.
	if WindowsOS {
		DirLibs = javaHome + PathStringSep + "bin"
		LibJvm = DirLibs + PathStringSep + "jvm." + LibExt
	} else {
		DirLibs = javaHome + PathStringSep + "lib"
		LibJvm = DirLibs + PathStringSep + "libjvm." + LibExt
	}

	// Connect to libjvm.
	_ = ConnectLibrary(LibJvm)
	log.Println("bridges/Setup: End")

}

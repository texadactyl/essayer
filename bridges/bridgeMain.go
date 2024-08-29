package bridges

import (
	"github.com/ebitengine/purego"
	"log"
	"os"
	"runtime"
)

type JNIboolean uint8
type JNIchar uint16
type JNIint int32
type JNIlong int64
type JNIshort int16
type JNIfloat float32
type JNIdouble float64

// Needed for CreateJvm.
type t_JavaVMInitArgs struct {
	version            JNIint
	nOptions           JNIint
	JavaVMOption       uintptr
	ignoreUnrecognized JNIboolean
}

var JavaVMInitArgs = t_JavaVMInitArgs{version: 0x00090000, nOptions: 0, JavaVMOption: 0, ignoreUnrecognized: 0}

var OperSys string                           // One of: "darwin", "linux", "unix", "windows"
var WindowsOS = false                        // true only if OperSys = "windows"
var PathDirLibs string                       // Directory of the more common JVM libraries (E.g. libzip.so)
var PathLibjvm string                        // Full path of libjvm.so
var PathLibjava string                       // Full path of libjava.so
var FileExt string                           // File extension of a library file: "so" (Linux and Unix), "dll" (Windows), "dylib" (MacOS)
var SepPathString = string(os.PathSeparator) // ";" (Windows) or ":" (everybody else)
var HandleLibjvm uintptr                     // Handle of the open libjvm
var HandleLibjava uintptr                    // Handle of the open libjava
var HandleJVM uintptr                        // Handle of the created JVM
var HandleENV uintptr                        // Handle of the JNI environment

func Setup() {

	log.Println("bridges/Setup: Begin")

	// Set up library file extension and library path string as a function of O/S.
	OperSys = runtime.GOOS
	switch OperSys {
	case "darwin":
		FileExt = "dylib"
	case "linux":
		FileExt = "so"
	case "windows":
		FileExt = "dll"
		WindowsOS = true
	default:
		log.Fatalln("bridges/Setup: Unsupported O/S: %s", OperSys)
	}

	// Get Java home.
	javaHome := os.Getenv("JAVA_HOME")
	if javaHome == "" {
		log.Fatalln("bridges/Setup: Environment variable JAVA_HOME missing but is required. Exiting.")
	}

	// Calculate some needed paths.
	if WindowsOS {
		PathDirLibs = javaHome + SepPathString + "bin"
		PathLibjvm = PathDirLibs + SepPathString + "server" + SepPathString + "jvm." + FileExt
		PathLibjava = PathDirLibs + SepPathString + "java." + FileExt
	} else {
		PathDirLibs = javaHome + SepPathString + "lib"
		PathLibjvm = PathDirLibs + SepPathString + "server" + SepPathString + "libjvm." + FileExt
		PathLibjava = PathDirLibs + SepPathString + "libjava." + FileExt
	}

	// Connect to libjvm.
	HandleLibjvm = ConnectLibrary(PathLibjvm)
	log.Println("bridges/Setup: connect to libjvm ok")

	// Connect to libjava.
	HandleLibjava = ConnectLibrary(PathLibjava)
	log.Println("bridges/Setup: connect to libjava ok")

	// Register the JVM creator library function.
	funcName := "JNI_CreateJavaVM"
	var createJvm func(*uintptr, *uintptr, *t_JavaVMInitArgs) JNIint // (& ptr to JVM, & ptr to env, & arguments) returns JNIint
	purego.RegisterLibFunc(&createJvm, HandleLibjvm, funcName)
	log.Printf("bridges/Setup: purego.RegisterLibFunc (%s) ok\n", funcName)

	// Create the JVM.
	ret := createJvm(&HandleJVM, &HandleENV, &JavaVMInitArgs)
	if ret < 0 {
		log.Fatalln("bridges/Setup: Cannot create a JVM. Exiting.")
	}
	log.Printf("bridges/Setup: createJvm ok\n")

	// Register the GetEnv library function.
	funcName = "JNU_GetEnv"
	var getEnv func(uintptr, *uintptr, JNIint) JNIint // (ptr to JVM, & ptr to env,JNI version) returns JNIint
	purego.RegisterLibFunc(&getEnv, HandleLibjava, funcName)
	log.Printf("bridges/Setup: purego.RegisterLibFunc (%s) ok\n", funcName)

	// Get the JNI environment pointer.
	ret = getEnv(HandleJVM, &HandleENV, JavaVMInitArgs.version)
	if ret < 0 {
		log.Fatalln("bridges/Setup: Cannot get the JNI environment pointer. Exiting.")
	}
	log.Printf("bridges/Setup: getEnv ok\n")

	log.Println("bridges/Setup: End")
}

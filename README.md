# essayer

This repo is for trying out things with ```purego```.

**2024-08-24 12:00**

So far, it works as expected in a GitHub Actions job in ```ubuntu-latest``` and ```macos-latest``` for libzip open and close.
I need to get that working in ```windows-latest``` too.
Also, I need to try some real JVM library functions (E.g. libzip inflate and deflate).

**2024-08-24 13:30**

Purego does not support the same user-level APIs in every O/S. In fact, they require the user to employ https://pkg.go.dev/golang.org/x/sys/windows for opening a library in the Windows environment. 

**2024-08-25**

Note: Variadic library functions (variable number of arguments) are not supported.

**2024-08-26**

On all platforms, successfully opened the JVM library and then the ZIP library.
<br>
Windows: What a royal pain in the hindquarters!
* O/S differences (historical)
* Java library subdirectory location (%JAVA_HOME%\bin on Windows vs ${JAVA_HOME}/lib on Posix)
* Java library file ame prefixes ("" on Windows vs "lib" on Posix E.g. zip.dll vs libzip.so)

**2024-08-27**

Executed the following libzip functions successfully on all platforms. Checked results with https://crc32.online/ :
* ZIP_CRC32
* Java_java_util_zip_CRC32_update

**2024-08-28**

Working on zip deflate and inflate.
Some observations of the libzip code:
* The environment pointer is used for throwing exceptions. Maybe other uses? I haven't yet run into a function that uses the object pointer. Yet.
* In general,
<br>- Very little comments (mostly, none) in the source code.
<br>- Symbol definitions are all over the place and, generally, without comments, of course. Pardon my sarcasm.
<br>- Parameters established during initialization of a z_stream structure seem to be required to be presented again (jint flush, jint params) even though the z_stream is also a parameter.
* My assumption is that the jacobin Go code wrapper for natives will be fed all required parameter values on the stack and they must match up to the requested function's definition in our native table (TBD).
* Extract from src/java.base/{unix,windows}/native/include/jni_md.h
```
    typedef int jint;                  // int64 assuming 64-bit architecture
```
* Extract from src/java.base/share/native/include/jni.h:
```
    typedef unsigned char   jboolean;  // uint8
    typedef unsigned short  jchar;     // uint16
    typedef short           jshort;    // int16
    typedef float           jfloat;    // float32
    typedef double          jdouble;   // float64 assuming 64-bit architecture
```
* The JNI environment is required for several functions (E.g. Java_java_util_zip_Deflater_deflateBytesBytes) just to execute; a dummy pointer cannot be used. Each JNI environment belongs to a unique JVM thread.

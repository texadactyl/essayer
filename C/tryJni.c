#include <dlfcn.h>
#include <inttypes.h>
#include <jni.h>
#include <pthread.h>
#include <stdarg.h>
#include <stdio.h>
#include <stdlib.h> 
#include <string.h> 

// JVM pointer.
JavaVM *jvm = NULL;

void OOPS(const char *fmt, ...) {
    char buffer[1024];
    va_list args;
    va_start(args, fmt);
    vsnprintf(buffer, sizeof(buffer), fmt, args);
    va_end(args);
    printf("*** Oops, %s", buffer);
    exit(86);
}

void * getSymAddr(void * libHandle, char * nameFunc) {
    void * addr = dlsym(libHandle, nameFunc);
    if (addr == NULL)
        OOPS("*** Oops, getSymAddr: Failed to find function %s\n", nameFunc);
    printf("getSymAddr: Found function %s ok\n", nameFunc);
    return addr;
}

void tryDeflater(JNIEnv * env, void * libzip) {

    char * libName;
    char * funcName;

	// -------------------------------------------------------------------------------------------
	// JNIEXPORT jlong JNICALL
	// Java_java_util_zip_Deflater_init(JNIEnv *env, jclass cls, jint level,
	//                                 jint strategy, jboolean nowrap)	
    jlong (* fnInit) (JNIEnv *env, jclass cls, jint level, jint strategy, jboolean nowrap);
    
	funcName = "Java_java_util_zip_Deflater_init";	
	fnInit = getSymAddr(libzip, funcName);
	long stream = (*fnInit) (env, NULL, 9, 0, 0x00);
	printf("tryDeflater: fnInit ok, stream = %ld\n", stream);
	
	// -------------------------------------------------------------------------------------------
	// JNIEXPORT jlong JNICALL
	// Java_java_util_zip_Deflater_deflateBytesBytes(JNIEnv *env, jobject this, jlong addr,
	//                                         jbyteArray inputArray, jint inputOff, jint inputLen,
	//                                         jbyteArray outputArray, jint outputOff, jint outputLen,
	//                                         jint flush, jint params)	
    jlong (* fnDeflate) (JNIEnv *env, jclass cls, jlong stream,
                            jbyteArray inputArray, jint inputOff, jint inputLen,
                            jbyteArray outputArray, jint outputOff, jint outputLen,
                            jint flush, jint params);
                            
	funcName = "Java_java_util_zip_Deflater_deflateBytesBytes";	
	fnDeflate = getSymAddr(libzip, funcName);

    char input[256] = "This is a test input string for compression.";
    int outputBufferSize = 1024;  // Start with a larger size
    char *output = (char*)malloc(outputBufferSize);
    if (!output)
        OOPS("tryDeflater: Failed to allocate memory for output buffer\n");

    // Convert char arrays to jbyteArray
    jbyteArray inputArray = (*env)->NewByteArray(env, sizeof(input));
    jbyteArray outputArray = (*env)->NewByteArray(env, outputBufferSize);

    // Set the contents of the input jbyteArray
    (*env)->SetByteArrayRegion(env, inputArray, 0, sizeof(input), (jbyte*)input);


    jlong result = (*fnDeflate) (env, NULL, stream, 
                                        inputArray, 0, strlen(input),
                                        outputArray, 0, outputBufferSize,
                                        0, 0);

    // Retrieve the output array and the actual output size.
    (*env)->GetByteArrayRegion(env, outputArray, 0, sizeof(output), (jbyte*)output);
	printf("tryDeflater: fnDeflate ok, result = %ld\n", result);
	
    // The result indicates the number of bytes written to the output buffer
    if (result > 0 && result < outputBufferSize) {
        // Retrieve only the valid portion of the output array
        (*env)->GetByteArrayRegion(env, outputArray, 0, (jsize)result, (jbyte*)output);

        // Print the compressed data size and content
        printf("tryDeflater: Deflate result (output size): %" PRId64 " bytes\n", result);
        printf("tryDeflater: Compressed Output: %.*s\n", (int)result, output);
    } else {
        OOPS("tryDeflater: Compression failed or output buffer is too small.\n");
    }
	
	// -------------------------------------------------------------------------------------------
	// JNIEXPORT void JNICALL Java_java_util_zip_Deflater_end(JNIEnv *env, jclass cls, jlong addr)
    jlong (* fnEnd) (JNIEnv *env, jclass cls, jlong stream);

	funcName = "Java_java_util_zip_Deflater_end";
	fnEnd = getSymAddr(libzip, funcName);
	stream = (*fnEnd) (env, NULL, stream);
	printf("tryDeflater: fnEnd ok\n");
	
}

// Function that the native thread will execute
void * child_thread(void * dummy_arg) {

    JNIEnv * env = NULL;
    char *className = "java/lang/String";

    // JAVA_HOME
    char * javaHome;

    // Attach the current thread to the JVM.
    printf("child_thread: Let's try to attach myself to the JVM .....\n");
    if ((*jvm)->AttachCurrentThread(jvm, (void**)&env, NULL) != 0)
        OOPS("*** Oops, child_thread: Failed to attach thread to the JVM\n");
    printf("child_thread: Thread successfully attached to the JVM\n");

    // Find a Java class by calling a function through an environment vector.
    jclass cls = (*env)->FindClass(env, className);
    if (cls == NULL)
        OOPS("*** Oops, child_thread: Failed to find class %s\n", className);
    printf("child_thread: Found class %s\n", className);
    
    // JVM library paths abd handles.
    char pathDirLibs[256];
    char pathLibJvm[256+32];
    void * libjvm;
    char pathLibZip[256+32];
    void * libzip;

    // Get JAVA_HOME
    javaHome = getenv("JAVA_HOME");
    if (javaHome == NULL)
        OOPS("*** Oops, child_thread: Failed to find environment variable JAVA_HOME\n");
    printf("child_thread: JAVA_HOME: %s\n", javaHome);
    
    // JVM library directory for nearly all functions.
    sprintf(pathDirLibs, "%s/lib", javaHome);

    // Open libjvm.
    sprintf(pathDirLibs, "%s/lib", javaHome);
    printf("child_thread: pathDirLibs: %s\n", pathDirLibs);
    sprintf(pathLibJvm, "%s/server/libjvm.so", pathDirLibs);
    printf("child_thread: pathLibJvm: %s\n", pathLibJvm);
    libjvm = dlopen(pathLibJvm, RTLD_LAZY);
    if (libjvm == NULL) 
        OOPS("*** Oops, child_thread: dlopen(%s) failed\n", pathLibJvm);
    printf("child_thread: dlopen(%s) ok\n", pathLibJvm);
    
    // Open libzip.
    sprintf(pathLibZip, "%s/libzip.so", pathDirLibs);
    libzip = dlopen(pathLibZip, RTLD_LAZY);
    if (libzip == NULL)
        OOPS("*** Oops, child_thread: dlopen(%s) failed\n", pathLibZip);
    printf("child_thread: dlopen(%s) ok\n", pathLibZip);
    
    // Try deflater functions.
    tryDeflater(env, libzip);
    
    // Detach the thread before exiting
    (*jvm)->DetachCurrentThread(jvm);
    printf("child_thread: Thread successfully detached from the JVM\n");
    
    return NULL;
}

int main() {

    // JNI environment pointer in the ancestor thread (main).
    JNIEnv *env = NULL;

    printf("main: Begin\n");

    JavaVMInitArgs vm_args;
    JavaVMOption options[1];
    options[0].optionString = "-Djava.class.path=.";
    vm_args.version = JNI_VERSION_1_8;
    vm_args.nOptions = 1;
    vm_args.options = options;
    vm_args.ignoreUnrecognized = 0;

    // Create the JVM.
    if (JNI_CreateJavaVM(&jvm, (void**)&env, &vm_args) < 0) 
        OOPS("*** Oops, main: Failed to create JVM\n");
    printf("main: JNI_CreateJavaVM ok\n");

    // Create a native thread.
    pthread_t thread;
    if (pthread_create(&thread, NULL, child_thread, NULL) != 0)
        OOPS("*** Oops, main: Failed to create native thread\n");
    printf("main: pthread_create ok, waiting for completion .....\n");

    // Wait for the thread to finish.
    pthread_join(thread, NULL);

    // Destroy the JVM.
    (*jvm)->DestroyJavaVM(jvm);

    printf("main: End\n");
    
    return 0;
}


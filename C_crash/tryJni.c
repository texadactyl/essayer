#include <dlfcn.h>
#include <jni.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

// JVM pointer.
JavaVM *jvm = NULL;

// JNI environment pointer.
JNIEnv *env = NULL;

int main() {
    printf("main: Begin\n");

    JavaVMInitArgs vm_args;
    JavaVMOption options[1];
    options[0].optionString = "-Djava.class.path=.";
    vm_args.version = JNI_VERSION_1_8;
    vm_args.nOptions = 1;
    vm_args.options = options;
    vm_args.ignoreUnrecognized = 0;

    // Create the JVM.
    if (JNI_CreateJavaVM(&jvm, (void**)&env, &vm_args) < 0) {
        printf("*** Oops, main: Failed to create JVM\n");
        return 1;
    }

    // Find a Java class by calling a function through an environment vector.
    char * className = "java/lang/String";
    jclass cls = (*env)->FindClass(env, className);
    if (cls == NULL) {
        printf("*** Oops, main: Failed to find class %s\n", className);
        return 1;
    }
    printf("main: Found class %s\n", className);
    
    // Get JAVA_HOME
    char * javaHome = getenv("JAVA_HOME");
    if (javaHome == NULL) {
        printf("*** Oops, main: Failed to find environment variable JAVA_HOME\n");
        return 1;
    }
    
    // Open libjvm.
    char * pathDirLibs; sprintf(pathDirLibs, "%s/lib", javaHome);
    char * pathLibJvm; sprintf(pathLibJvm, "%s/server/libjvm.so", pathDirLibs);
    void * libjvm = dlopen(pathLibJvm, RTLD_LAZY);
    if (libjvm == NULL) {
        printf("*** Oops, main: dlopen(%s) failed\n", pathLibJvm);
        return 1;
    }
    printf("main: dlopen(%s) ok\n", pathLibJvm);
    
    printf("main: End\n");
    
    return 0;
}


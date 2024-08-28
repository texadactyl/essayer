#include <jni.h>
#include <pthread.h>
#include <stdio.h>

// Global variable to hold the JavaVM reference
JavaVM *jvm = NULL;

// Function that the native thread will execute
void* thread_routine(void* arg) {
    JNIEnv *env = NULL;

    // Attach the current thread to the JVM
    if ((*jvm)->AttachCurrentThread(jvm, (void**)&env, NULL) != 0) {
        printf("Failed to attach thread to the JVM\n");
        return NULL;
    }

    // Now you can use the JNIEnv pointer to interact with the JVM
    printf("Thread successfully attached to the JVM\n");

    // For example, find a Java class and call a method (simplified)
    jclass cls = (*env)->FindClass(env, "java/lang/String");
    if (cls != NULL) {
        printf("Found java/lang/String class\n");
    } else {
        printf("Failed to find java/lang/String class\n");
    }

    // Detach the thread before exiting
    (*jvm)->DetachCurrentThread(jvm);

    printf("Thread successfully detached from the JVM\n");
    return NULL;
}

int main() {
    // Normally, JVM creation is handled by the Java launcher, but here’s an example:
    JavaVMInitArgs vm_args;
    JavaVMOption options[1];
    options[0].optionString = "-Djava.class.path=.";
    vm_args.version = JNI_VERSION_1_8;
    vm_args.nOptions = 1;
    vm_args.options = options;
    vm_args.ignoreUnrecognized = 0;

    // Create the JVM
    JNIEnv *env; // Add this line to define env
    if (JNI_CreateJavaVM(&jvm, (void**)&env, &vm_args) < 0) {
        printf("Failed to create JVM\n");
        return 1;
    }

    // Create a native thread
    pthread_t thread;
    if (pthread_create(&thread, NULL, thread_routine, NULL) != 0) {
        printf("Failed to create native thread\n");
        return 1;
    }

    // Wait for the thread to finish
    pthread_join(thread, NULL);

    // Destroy the JVM (optional, usually not done until program exit)
    (*jvm)->DestroyJavaVM(jvm);

    return 0;
}


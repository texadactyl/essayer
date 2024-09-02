#include <jni.h>
#include <pthread.h>
#include <stdio.h>

// JVM reference.
JavaVM *jvm = NULL;

// Function that the native thread will execute
void* thread_routine(void* arg) {
    JNIEnv *env = NULL; // Set by AttachCurrentThread
    char *className = "java/lang/String";

    // Attach the current thread to the JVM.
    if ((*jvm)->AttachCurrentThread(jvm, (void**)&env, NULL) != 0) {
        printf("*** Oops, thread_routine: Failed to attach thread to the JVM\n");
        return NULL;
    }
    printf("Thread successfully attached to the JVM\n");

    // Find a Java class by calling a function through an environment vector.
    jclass cls = (*env)->FindClass(env, className);
    if (cls == NULL) {
        printf("*** Oops, thread_routine: Failed to find class %s\n", className);
        return NULL;
    }
    printf("thread_routine: Found class %s\n", className);
    
    // Detach the thread before exiting
    (*jvm)->DetachCurrentThread(jvm);

    printf("thread_routine: Thread successfully detached from the JVM\n");
    return NULL;
}

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
    JNIEnv *env; // Add this line to define env
    if (JNI_CreateJavaVM(&jvm, (void**)&env, &vm_args) < 0) {
        printf("*** Oops, main: Failed to create JVM\n");
        return 1;
    }

    // Create a native thread.
    pthread_t thread;
    if (pthread_create(&thread, NULL, thread_routine, NULL) != 0) {
        printf("*** Oops, main: Failed to create native thread\n");
        return 1;
    }

    // Wait for the thread to finish.
    pthread_join(thread, NULL);

    // Destroy the JVM (unnecessary in the case of this sample code).
    (*jvm)->DestroyJavaVM(jvm);

    printf("main: End\n");
    
    return 0;
}


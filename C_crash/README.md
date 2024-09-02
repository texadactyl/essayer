# Lessons Learned

The JVM and JNI are intimately intertwined. Each thread that calls JNI functions must have its own JNI environment. Trying to run a C program inside the same thread as the JVM seems to always lead to a crash. I could be confused.

Java_java_util_zip_Deflater_deflateBytesBytes will not work without using the environment utility functions. See ```essayer/C```. That one doesn't crash although the resulting output length makes no sense. The same is most likely true for the Go version of a Java_java_util_zip_Deflater_deflateBytesBytes caller. But, as far as I can see, Go cannot execute functions expressed as vectors in the JNI environment.

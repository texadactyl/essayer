# Lesson Learned

The JVM and JNI are intimately intertwined. Trying to run a C program without its own thread that creates a JVM starts threading anyways. A crash is inevitable.

Java_java_util_zip_Deflater_deflateBytesBytes will not work without using the environment utility functions. See ```essayer/C```. That one doesn't crash although the resulting output length makes no sense.

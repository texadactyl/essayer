* Lesson Learned

The JVM and JNI are inteimately intertwined. Trying to run a C program without its own thread that creates a JVM starts threading anyways. A crash is inevitable. See post-mortem log file ```hs_err_pid618262.log```.

JAVABASE=/home/elkins/BASIS/OpenJDK/src/java.base
LIBDIR_SERVER=$JAVA_HOME/lib/server
LIBDIR_ZIP=$JAVA_HOME/lib
set -x
gcc -o essai -I $JAVA_HOME/include/ -I $JAVA_HOME/include/linux/ tryAttachThread.c \
    -L $LIBDIR_SERVER -l :libjvm.so -L $LIBDIR_ZIP -l :libzip.so


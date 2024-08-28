SRCBASE=/home/elkins/BASIS/OpenJDK/src/java.base
LIBDIR=/usr/lib/jvm/java-17-openjdk-amd64/lib/server
set -x
gcc -o essai -I $SRCBASE/share/native/include/ -I $SRCBASE/unix/native/include tryAttachThread.c \
    -L $LIBDIR -l :libjvm.so


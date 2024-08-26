//go:build windows

package bridges

import (
	"golang.org/x/sys/windows"
	"log"
)

func ConnectLibrary(libPath string) uintptr {
	var handle uintptr
	var err error
	handle, err = windows.LoadLibrary(libPath)
	if err != nil {
		log.Fatalf("ConnectLibrary: purego.Dlopen for [%s] failed, reason: [%s]\n",
			libPath, err.Error())
	}
	return handle
}

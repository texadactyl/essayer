//go:build windows

package bridges

import (
	"golang.org/x/sys/windows"
	"log"
)

func ConnectLibrary(libPath string) uintptr {
	handle, err := windows.LoadLibrary(libPath)
	if err != nil {
		log.Fatalf("ConnectLibrary: windows.LoadLibrary for [%s] failed, reason: [%s]\n",
			libPath, err.Error())
	}
	return uintptr(handle)
}

//go:build windows
// +build windows

package main

import (
	"golang.org/x/sys/windows"
	"log"
)

func connectLibrary(libPath string) uintptr {
	var handle uintptr
	var err error
	handle, err = windows.LoadLibrary(libPath)
	if err != nil {
		log.Fatalf("connectLibrary: purego.Dlopen for [%s] failed, reason: [%s]\n",
			libPath, err.Error())
	}
	return handle
}

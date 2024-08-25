//go:build darwin || freebsd || linux || windows

package main

import (
	"github.com/ebitengine/purego"
	"golang.org/x/sys/windows"
	"log"
)

func connectLibrary(libPath string) uintptr {
	var handle uintptr
	var err error
	if WindowsOS {
		handle, err = windows.LoadLibrary(libPath)
		if err != nil {
			log.Fatalf("connectLibrary: purego.Dlopen for [%s] failed, reason: [%s]\n",
				LibJvm, err.Error())
		}
	} else {
		handle, err = purego.Dlopen(libPath, purego.RTLD_NOW|purego.RTLD_GLOBAL)
		if err != nil {
			log.Fatalf("connectLibrary: purego.Dlopen for [%s] failed, reason: [%s]\n",
				LibJvm, err.Error())
		}
	}
	return handle
}

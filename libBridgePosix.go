//go:build !windows
// +build !windows

package main

import (
	"github.com/ebitengine/purego"
	"log"
)

func connectLibrary(libPath string) uintptr {
	var handle uintptr
	var err error
	handle, err = purego.Dlopen(libPath, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		log.Fatalf("connectLibrary: purego.Dlopen for [%s] failed, reason: [%s]\n",
			libPath, err.Error())
	}
	return handle
}

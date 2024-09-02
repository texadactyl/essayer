//go:build !windows

package bridges

import (
	"github.com/ebitengine/purego"
	"log"
)

func ConnectLibrary(libPath string) uintptr {
	var handle uintptr
	var err error
	handle, err = purego.Dlopen(libPath, purego.RTLD_LAZY)
	if err != nil {
		log.Fatalf("ConnectLibrary: purego.Dlopen for [%s] failed, reason: [%s]\n",
			libPath, err.Error())
	}
	return handle
}

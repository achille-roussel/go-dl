// +build darwin linux

package dl

// #cgo CFLAGS: -W -Wall -Wno-unused-parameter -O3
// #cgo LDFLAGS: -ldl
//
// #include <dlfcn.h>
// #include <stdlib.h>
import "C"
import (
	"runtime"
	"strings"
	"sync"
	"syscall"
	"unsafe"
)

type dylib struct {
	mutex  sync.RWMutex
	handle unsafe.Pointer
}

var mutex sync.Mutex

func open(path string, mode Mode) (lib *dylib, err error) {
	var cpath *C.char
	var handle unsafe.Pointer

	if len(path) != 0 {
		if strings.Index(path, ext) < 0 {
			path += ext
		}
		cpath = C.CString(path)
		defer C.free(unsafe.Pointer(cpath))
	}

	mutex.Lock()
	defer mutex.Unlock()

	if handle = C.dlopen(cpath, makeMode(mode)); handle == nil {
		err = lastError()
		return
	}

	lib = &dylib{
		handle: handle,
	}

	runtime.SetFinalizer(lib, (*dylib).close)
	return
}

func (lib *dylib) Close() (err error) {
	lib.mutex.Lock()
	err = lib.close()
	lib.mutex.Unlock()
	return
}

func (lib *dylib) Symbol(name string) (addr uintptr, err error) {
	var handle unsafe.Pointer
	var ptr unsafe.Pointer
	var sym *C.char

	sym = C.CString(name)
	defer C.free(unsafe.Pointer(sym))

	lib.mutex.RLock()
	defer lib.mutex.RUnlock()

	if handle = lib.handle; handle == nil {
		err = syscall.EINVAL
		return
	}

	mutex.Lock()
	defer mutex.Unlock()

	if ptr = C.dlsym(handle, sym); ptr == nil {
		err = lastError()
		return
	}

	addr = uintptr(ptr)
	return
}

func (lib *dylib) close() (err error) {
	var handle unsafe.Pointer

	if handle, lib.handle = lib.handle, nil; handle != nil {
		mutex.Lock()
		defer mutex.Unlock()

		if C.dlclose(handle) != 0 {
			err = lastError()
		}
	}

	return
}

func makeMode(mode Mode) (c C.int) {
	if (mode & Lazy) != 0 {
		c |= C.RTLD_LAZY
	}

	if (mode & Now) != 0 {
		c |= C.RTLD_NOW
	}

	if (mode & Global) != 0 {
		c |= C.RTLD_GLOBAL
	}

	if (mode & Local) != 0 {
		c |= C.RTLD_LOCAL
	}

	return
}

func lastError() error {
	var s string

	if err := C.dlerror(); err == nil {
		s = "unknown"
	} else {
		s = C.GoString(err)
	}

	return &Error{
		Message: s,
	}
}

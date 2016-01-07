package dl

import "testing"

func TestOpenLazyGlobal(t *testing.T) {
	var lib Library
	var err error

	if lib, err = Open("libc", Lazy|Global); err != nil {
		t.Error("open:", err)
		return
	}

	if err = lib.Close(); err != nil {
		t.Error("close:", err)
		return
	}
}

func TestOpenNowLocal(t *testing.T) {
	var lib Library
	var err error

	if lib, err = Open("libc", Now|Local); err != nil {
		t.Error("open:", err)
		return
	}

	if err = lib.Close(); err != nil {
		t.Error("close:", err)
		return
	}
}

func TestSymbol(t *testing.T) {
	var lib Library
	var err error
	var ptr uintptr

	if lib, err = Open("libc", Lazy|Local); err != nil {
		t.Error("open:", err)
		return
	}

	defer lib.Close()

	if ptr, err = lib.Symbol("printf"); err != nil {
		t.Error("symbol:", err)
		return
	}

	if ptr == 0 {
		t.Error("null pointer returned by Library.Symbol")
		return
	}
}

func TestOpenError(t *testing.T) {
	if _, err := Open("something-weird", Lazy|Local); err == nil {
		t.Error("error:", err)
	} else {
		t.Log(err)
	}
}

func TestSymbolError(t *testing.T) {
	var lib Library
	var err error

	if lib, err = Open("libc", Lazy|Local); err != nil {
		t.Error("open:", err)
		return
	}

	lib.Close()

	if _, err = lib.Symbol("printf"); err == nil {
		t.Error("symbol: error expected after closing the library")
	}
}

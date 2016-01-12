go-dl [![Build Status](https://travis-ci.org/achille-roussel/go-dl.svg)](https://travis-ci.org/achille-roussel/go-dl) [![Coverage Status](https://coveralls.io/repos/achille-roussel/go-dl/badge.svg?branch=master&service=github)](https://coveralls.io/github/achille-roussel/go-dl?branch=master)
=====

*go-dl* is a package exposing dynamic library loading features to the Go language.

Dynamically Loading C Libraries
-------------------------------

The *go-dl* package exposes the `Open` function that taking a library name or
path and some options will load the code in memory and return an object allowing
interaction with the library.

Here's an example:

```go
package main

import (
    "fmt"
    "unsafe"

    "github.com/achille-roussel/go-dl"
)

func main() {
    var lib dl.Library
    var err error
    var sym unsafe.Pointer

    // Dynamically loads libc on a linux platform.
    // Note that the library names or paths are usually system-dependant.
    if lib, err = dl.Open("libc.so.6", dl.Lazy|dl.Local); err != nil {
        fmt.Println(err)
        return
    }

    defer lib.Close()

    // Get the address of the `puts` symbol, this can be used to call the
    // function using a ffi package.
    if sym, err = lib.Symbol("puts"); err != nil {
        fmt.Println(err)
        return
    }

    // ...
}
```

Finding C Libraries
-------------------

*go-dl* also comes with a function that attempts to emulate the platform's logic
for dynamic library discovery in a portable way, here's an example:

```go
package main

import (
    "fmt"

    "github.com/achille-roussel/go-dl"
)

func main() {
    var path string
    var err error

    if path, err = dl.Find("libc"); err != nil {
        fmt.Println(err)
        return
    }

    // Now `path` holds the full path to where the standard C library was found
    // on the system.
    fmt.Println(path)
}
```

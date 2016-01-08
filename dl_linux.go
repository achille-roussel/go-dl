package dl

// #cgo CFLAGS: -W -Wall -Wno-unused-parameter -O3
// #cgo LDFLAGS: -ldl
import "C"
import (
	"os"
	"path/filepath"
	"strings"
)

func find(name string) (path string, err error) {
	if strings.ContainsRune(name, '/') {
		path = name
		return
	}

	if len(filepath.Ext(name)) == 0 {
		name += ".so"
	}

	dirs := make([]string, 0, 100)
	dirs = append(dirs, getPaths("LD_LIBRARY_PATH")...)
	dirs = append(dirs, "/lib", "/usr/lib")

	ok := false

	for _, dir := range dirs {
		filepath.Walk(dir, func(p string, f os.FileInfo, e error) error {
			if !ok && f != nil && !f.IsDir() && (f.Mode() & 0600) != 0 && strings.HasPrefix(f.Name(), name) {
				path, ok = p, true
			}
			return nil
		})

		if ok {
			break
		}
	}

	if !ok {
		err = os.ErrNotExist
	}

	return
}

func getPaths(env string) []string {
	paths := strings.Split(os.Getenv(env), ":")
	i := 0

	for _, p := range paths {
		if len(p) != 0 {
			paths[i] = p
			i++
		}
	}

	return paths[:i]
}

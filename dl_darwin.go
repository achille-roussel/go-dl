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
		name += ".dylib"
	}

	p1 := getPaths("LD_LIBRARY_PATH")
	p2 := getPaths("DYLD_LIBRARY_PATH")
	p3 := getPaths("DYLD_FALLBACK_LIBRARY_PATH")

	if len(p3) == 0 {
		p3 = []string{
			filepath.Join(os.Getenv("HOME"), "lib"),
			"/usr/local/lib",
			"/usr/lib",
		}
	}

	dirs := make([]string, 0, 100)
	dirs = append(dirs, p1...)
	dirs = append(dirs, p2...)
	dirs = append(dirs, ".")
	dirs = append(dirs, p3...)

	ok := false

	for _, dir := range dirs {
		filepath.Walk(dir, func(p string, f os.FileInfo, e error) error {
			if !ok && f != nil && !f.IsDir() && strings.HasPrefix(f.Name(), name) {
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

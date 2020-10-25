package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func extension(filename string, isDir bool) string {
	if isDir {
		return ""
	}

	_, _filename := filepath.Split(filename)

	f := strings.Split(_filename, ".")
	var extension string

	if len(f) > 0 {
		extension = strings.Join(f[1:], ".")
	}

	return extension
}

func getFullPath(parentPath, filename string) string {
	pathSep := PathSep
	if parentPath == PathSep {
		pathSep = ""
	}
	if strings.HasSuffix(parentPath, PathSep) {
		pathSep = ""
	}
	if strings.HasPrefix(filename, PathSep) {
		pathSep = ""
	}

	_fullPath := fmt.Sprintf("%s%s%s", parentPath, pathSep, filename)

	return fixSlash(_fullPath)
}

func fixSlash(absFilepath string) string {
	_absFilepath := absFilepath

	if !strings.HasPrefix(_absFilepath, PathSep) {
		_absFilepath = fmt.Sprintf("%s%s", PathSep, _absFilepath)
	}

	if _absFilepath != PathSep && strings.HasSuffix(_absFilepath, PathSep) {
		_absFilepath = strings.TrimSuffix(_absFilepath, PathSep)
	}

	return _absFilepath
}

func indexExists(arr interface{}, index int) bool {
	switch value := arr.(type) {
	case *[]string:
		return len(*value) > index

	case []string:
		return len(value) > index

	default:
		log.Panic("invalid type in 'indexExists'")
	}

	return false
}

// Get Parent path of a list of directories and files
func GetParentPath(sep byte, paths ...string) string {
	// Handle special cases.
	switch len(paths) {
	case 0:
		return ""
	case 1:
		return path.Clean(paths[0])
	}

	// Note, we treat string as []byte, not []rune as is often
	// done in Go. (And sep as byte, not rune). This is because
	// most/all supported OS' treat paths as string of non-zero
	// bytes. A filename may be displayed as a sequence of Unicode
	// runes (typically encoded as UTF-8) but paths are
	// not required to be valid UTF-8 or in any normalized form
	// (e.g. "é" (U+00C9) and "é" (U+0065,U+0301) are different
	// file names.
	c := []byte(path.Clean(paths[0]))

	// We add a trailing sep to handle the case where the
	// common prefix directory is included in the path list
	// (e.g. /home/user1, /home/user1/foo, /home/user1/bar).
	// path.Clean will have cleaned off trailing / separators with
	// the exception of the root directory, "/" (in which case we
	// make it "//", but this will get fixed up to "/" bellow).
	c = append(c, sep)

	// Ignore the first path since it's already in c
	for _, v := range paths[1:] {
		// Clean up each path before testing it
		v = path.Clean(v) + string(sep)

		// Find the first non-common byte and truncate c
		if len(v) < len(c) {
			c = c[:len(v)]
		}
		for i := 0; i < len(c); i++ {
			if v[i] != c[i] {
				c = c[:i]
				break
			}
		}
	}

	// Remove trailing non-separator characters and the final separator
	for i := len(c) - 1; i >= 0; i-- {
		if c[i] == sep {
			c = c[:i]
			break
		}
	}

	return string(c)
}

func fileExistsLocal(filename string) bool {
	_, err := os.Stat(filename)

	return !os.IsNotExist(err)
}

func isFileLocal(name string) bool {
	if fi, err := os.Stat(name); err == nil {
		if fi.Mode().IsRegular() {
			return true
		}
	}

	return false
}

func isDirLocal(name string) bool {
	if fi, err := os.Stat(name); err == nil {
		if fi.Mode().IsDir() {
			return true
		}
	}
	return false
}

func isSymlinkLocal(fi os.FileInfo) bool {
	return fi.Mode()&os.ModeSymlink != 0
}

func isDisallowedFiles(filename string) bool {
	return StringContains(disallowedFiles, filename)
}

func existsLocal(filename string) bool {
	_, err := os.Stat(filename)

	return !os.IsNotExist(err)
}

func Percent(percent float32, all float32) float32 {
	return (percent / all) * 100
}

func StringFilter(x []string, f func(string) bool) []string {
	a := make([]string, 0)

	for _, v := range x {
		if f(v) && len(v) > 7 {
			a = append(a, v)
		}
	}

	return a
}

func StringContains(list []string, search string) bool {
	for _, a := range list {
		if a == search {
			return true
		}
	}

	return false
}

func subpathExists(path, searchPath string) bool {
	return path != "" && strings.HasPrefix(searchPath, path)
}

func mapLocalPathToMtpPath(
	sourcePath, sourceParentPath, destinationPath string,
) (destinationParentPath, destinationFilePath string) {
	trimmedSourcePath := strings.TrimPrefix(sourcePath, sourceParentPath)
	fullPath := getFullPath(destinationPath, trimmedSourcePath)

	return filepath.Dir(fullPath), fullPath
}

func SanitizeDosName(name string) string {
	if !strings.ContainsAny(name, disallowedFileName) {
		return name
	}
	dest := make([]byte, len(name))
	for i := 0; i < len(name); i++ {
		if strings.Contains(disallowedFileName, string(name[i])) {
			dest[i] = '_'
		} else {
			dest[i] = name[i]
		}
	}
	return string(dest)
}

func transferRateInMBs(size int64, lastSentTime time.Time, speed float64) float64 {
	var _size = float64(size)
	var elapsedTime = time.Since(lastSentTime).Seconds()

	rate := _size / elapsedTime / 1000 / 1000

	// prevent polluted values
	if elapsedTime <= 1 {
		return speed
	}

	return math.Round(rate*100) / 100
}

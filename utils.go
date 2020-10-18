package main

import (
	"fmt"
	"github.com/ganeshrvel/go-mtpfs/mtp"
	"log"
	"path/filepath"
	"strings"
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

func isObjectADir(obj *mtp.ObjectInfo) bool {
	return obj.ObjectFormat == mtp.OFC_Association
}

func getFullPath(parentPath, filename string) string {
	pathSep := PathSep
	if parentPath == PathSep {
		pathSep = ""
	}

	return fmt.Sprintf("%s%s%s", parentPath, pathSep, filename)
}

func fixDirSlash(absFilepath string) string {
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

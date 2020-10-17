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
	pathSep := "/"
	if parentPath == "/" {
		pathSep = ""
	}

	return fmt.Sprintf("%s%s%s", parentPath, pathSep, filename)
}

func fixDirSlash(absFilepath string) string {
	_absFilepath := absFilepath

	if !strings.HasPrefix(_absFilepath, "/") {
		_absFilepath = fmt.Sprintf("%s%s", "/", _absFilepath)
	}

	if _absFilepath != "/" && strings.HasSuffix(_absFilepath, "/") {
		_absFilepath = strings.TrimSuffix(_absFilepath, "/")
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

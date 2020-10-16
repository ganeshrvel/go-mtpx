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

func fixDirSlash(absFilepath string) string {
	if !strings.HasPrefix(absFilepath, "/") {
		absFilepath = fmt.Sprintf("%s%s", "/", absFilepath)
	}

	if absFilepath != "/" && strings.HasSuffix(absFilepath, "/") {
		absFilepath = strings.TrimSuffix(absFilepath, "/")
	}

	return absFilepath
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

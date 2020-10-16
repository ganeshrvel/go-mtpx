package main

import (
	"github.com/ganeshrvel/go-mtpfs/mtp"
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

func isDirectoryObject(obj *mtp.ObjectInfo) bool {
	return obj.ObjectFormat == mtp.OFC_Association
}

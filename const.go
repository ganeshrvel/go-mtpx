package main

import (
	"github.com/ganeshrvel/go-mtpfs/mtp"
	"os"
)

const PathSep = string(os.PathSeparator)

const ParentObjectId = mtp.GOH_ROOT_PARENT

const devTimeout = 15000

const disallowedFileName = ":*?\"<>|"

var disallowedFiles = []string{".DS_Store"}

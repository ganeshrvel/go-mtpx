package main

import (
	"github.com/ganeshrvel/go-mtpfs/mtp"
	"time"
)

type Init struct {
	debugMode bool
}

type StorageData struct {
	sid  uint32
	info mtp.StorageInfo
}

type FileInfo struct {
	Size       int64
	IsDir      bool
	ModTime    time.Time
	Name       string
	FullPath   string
	ParentPath string
	Extension  string
	ParentId   uint32
	ObjectId   uint32

	Info *mtp.ObjectInfo
}

type DirectoryTree map[uint32]*DirectoryInfo

type DirectoryInfo struct {
	*FileInfo

	Children   []*DirectoryTree
}

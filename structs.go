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
	FullPath   string //todo
	ParentPath string //todo
	Extension  string
	ParentId   uint32
	ObjectId   uint32

	Info *mtp.ObjectInfo
}

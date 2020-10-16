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
	Info         *mtp.ObjectInfo
	Size         int64
	IsDir        bool
	ModTime      time.Time
	Name         string
	FullPath     string //todo
	ParentPath   string //todo
	Extension    string
	ParentObject uint32
}

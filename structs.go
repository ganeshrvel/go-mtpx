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

type WalkCb func(objectId uint32, fi *FileInfo, err error) error

type TransferredFileInfo struct {
	FileInfo *FileInfo

	StartTime      time.Time
	LatestSentTime time.Time
	FilesSent      int
	Speed          float64
}

type TransferFilesCb func(fi *TransferredFileInfo, err error) error

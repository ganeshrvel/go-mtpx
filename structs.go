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

type WalkDirectoryCb func(objectId uint32, fi *FileInfo)

type UploadFileInfo struct {
	FileInfo *FileInfo

	startTime      time.Time
	latestSentTime time.Time
	sentFiles      int
	speed          float64
}

type UploadFilesCb func(uploadFi *UploadFileInfo)

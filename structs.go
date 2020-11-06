package mtpx

import (
	"github.com/ganeshrvel/go-mtpfs/mtp"
	"time"
)

type Init struct {
	DebugMode bool
}

type StorageData struct {
	Sid  uint32
	Info mtp.StorageInfo
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

type ProgressCb func(fi *ProgressInfo, err error) error

type TransferSizeInfo struct {
	// total size
	Total      int64

	// total size transferred
	Sent       int64

	// progress in percentage
	Progress float32
}

type ProgressInfo struct {
	FileInfo *FileInfo

	// transfer starting time
	StartTime time.Time

	// most recent transfer time
	LatestSentTime time.Time

	// transfer rate (in MB/s)
	Speed float64

	// total files to transfer
	// note: this value will be available only if the pre-processing of file transfer is enabled
	TotalFiles int64

	// total files transferred
	FilesSent int64

	// total file transfer progress in percentage
	FilesSentProgress float32

	// size information of the current file which is being transferred
	Current *TransferSizeInfo

	// total size information of the files for the transfer session
	Bulk *TransferSizeInfo
}

type SizeProgressCb func(total, sent int64, objectId uint32)

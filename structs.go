package mtpx

import (
	"github.com/ganeshrvel/go-mtpfs/mtp"
	"os"
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

type TransferSizeInfo struct {
	// total size to transfer
	// note: the value will be 0 if pre-processing was not allowed
	Total int64

	// total size transferred
	Sent int64

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
	// note: the value will be 0 if pre-processing was not allowed
	TotalFiles int64

	// total directories to transfer
	// note: the value will be 0 if pre-processing was not allowed
	TotalDirectories int64

	// total files transferred
	FilesSent int64

	// total file transfer progress in percentage
	FilesSentProgress float32

	// size information of the current file which is being transferred
	ActiveFileSize *TransferSizeInfo

	// total size information of the files for the transfer session
	BulkFileSize *TransferSizeInfo

	Status TransferStatus
}

type SizeProgressCb func(total, sent int64, objectId uint32, err error) error

type LocalWalkCb func(fi *os.FileInfo, err error) error

type ProgressCb func(fi *ProgressInfo, err error) error

type LocalPreprocessCb func(fi *os.FileInfo, err error) error

type MtpPreprocessCb func(fi *FileInfo, err error) error

type FileProp struct {
	ObjectId uint32
	FullPath string
}

type processDownloadFilesProps struct {
	destinationFileParentPath, destinationFilePath, sourceParentPath string
	bulkFilesSent, bulkSizeSent, totalFiles, totalSize               int64
}

type downloadFilesObjectCache map[string]downloadFilesObjectCacheContainer

type downloadFilesObjectCacheContainer struct {
	fileInfo                                                         *FileInfo
	destinationFileParentPath, destinationFilePath, sourceParentPath string
}

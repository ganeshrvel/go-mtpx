package mtpx

import (
	"errors"
	"fmt"
	"github.com/ganeshrvel/go-mtpfs/mtp"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// todo: work on documentations
// todo: hotplug
// todo: information mode -> get total files, break the buf into smaller chunks and calculate the transfer rate
// todo: implement the progress info, check for progressCount in gomtpfs
// update go.mod

// initialize the mtp device
// returns mtp device
func Initialize(init Init) (*mtp.Device, error) {
	dev, err := mtp.SelectDevice("")

	if err != nil {
		return nil, MtpDetectFailedError{error: err}
	}

	dev.MTPDebug = init.DebugMode
	dev.DataDebug = init.DebugMode
	dev.USBDebug = init.DebugMode

	dev.Timeout = devTimeout

	if err = dev.Configure(); err != nil {
		return nil, ConfigureError{error: err}
	}

	return dev, nil
}

// close the mtp device
func Dispose(dev *mtp.Device) {
	dev.Close()
}

// fetch device Info
func FetchDeviceInfo(dev *mtp.Device) (*mtp.DeviceInfo, error) {
	info := mtp.DeviceInfo{}
	err := dev.GetDeviceInfo(&info)

	if err != nil {
		return nil, DeviceInfoError{error: err}
	}

	return &info, nil
}

// fetch storages
func FetchStorages(dev *mtp.Device) ([]StorageData, error) {
	sids := mtp.Uint32Array{}
	if err := dev.GetStorageIDs(&sids); err != nil {
		return nil, StorageInfoError{error: err}
	}

	if len(sids.Values) < 1 {
		return nil, NoStorageError{error: fmt.Errorf("no storage found")}
	}

	var result []StorageData

	for _, sid := range sids.Values {
		var info mtp.StorageInfo
		if err := dev.GetStorageInfo(sid, &info); err != nil {
			return nil, StorageInfoError{error: err}
		}

		result = append(result, StorageData{
			Sid:  sids.Values[0],
			Info: info,
		})
	}

	return result, nil
}

// create a new directory recursively using [fullPath]
// The path will be created if it does not exists
func MakeDirectory(dev *mtp.Device, storageId uint32, fullPath string) (rObjectId uint32, rError error) {
	_fullPath := fixSlash(fullPath)

	if _fullPath == PathSep {
		return ParentObjectId, nil
	}
	splittedFullPath := strings.Split(_fullPath, PathSep)

	var objectId = uint32(ParentObjectId)
	const skipIndex = 1

	for _, fName := range splittedFullPath[skipIndex:] {
		// fetch the parent object and
		fi, err := GetObjectFromParentIdAndFilename(dev, storageId, objectId, fName)

		if err != nil {
			switch err.(type) {
			case FileNotFoundError:
				// if object does not exists then create a new directory
				_newObjectId, err := handleMakeDirectory(dev, storageId, objectId, fName)
				if err != nil {
					return 0, err
				}

				objectId = _newObjectId

				continue
			default:
				return 0, err
			}
		}

		// if the object exists but if it's a file then throw an error
		if !fi.IsDir {
			return 0, InvalidPathError{error: fmt.Errorf("invalid path: %s. The object is not a directory", fName)}
		}

		objectId = fi.ObjectId
	}

	return objectId, nil
}

// List the contents in a directory
// use [recursive] to fetch the whole nested tree
// Tip: use [objectId] whenever possible to avoid traversing down the whole file tree to process and find the [objectId]
// if [skipDisallowedFiles] is true then files matching the [disallowedFiles] list will be ignored
// rObjectId: objectId of the file/diectory
// rTotalFiles: total number of files and directories
func Walk(dev *mtp.Device, storageId uint32, fullPath string, recursive bool, skipDisallowedFiles bool, cb WalkCb) (rObjectId uint32, rTotalFiles int, rError error) {
	// fetch the objectId from [objectId] and/or [fullPath] parameters
	fi, err := GetObjectFromPath(dev, storageId, fullPath)
	if err != nil {
		return 0, 0, err
	}

	// if the object file name matches [disallowedFiles] list then return an error
	if skipDisallowedFiles {
		fName := (*fi).Name
		if ok := isDisallowedFiles(fName); ok {
			return 0, 0, InvalidPathError{error: fmt.Errorf("disallowed file %v", fName)}
		}
	}

	// if the object is a file then return objectId
	if !fi.IsDir {
		err := cb(fi.ObjectId, fi, nil)
		if err != nil {
			return 0, 0, err
		}

		return fi.ObjectId, 1, nil
	}

	totalFiles, err := proccessWalk(dev, storageId, fi.ObjectId, fullPath, recursive, skipDisallowedFiles, cb)
	if err != nil {
		return 0, 0, err
	}

	return fi.ObjectId, totalFiles, nil
}

// check if a file exists
// returns exists: bool, isDir: bool, objectId: uint32
// Since the [parentPath] is unavailable here the [fullPath] property of the resulting object [FileInfo] may not be valid.
func FileExists(dev *mtp.Device, storageId, objectId uint32, filePath string) (exists bool, fileInfo *FileInfo) {
	fi, err := GetObjectFromObjectIdOrPath(dev, storageId, objectId, filePath)

	if err != nil {
		return false, nil
	}

	return true, fi
}

// Delete an file/directory
// [objectId] and [fullPath] are optional parameters
// if [objectId] is not available then [fullPath] will be used to fetch the [objectId]
// dont leave both [objectId] and [fullPath] empty
// Tip: use [objectId] whenever possible to avoid traversing down the whole file tree to process and find the [objectId]
func DeleteFile(dev *mtp.Device, storageId, objectId uint32, fullPath string) error {
	exist, fi := FileExists(dev, storageId, objectId, fullPath)

	if !exist {
		return nil
	}

	if err := dev.DeleteObject(fi.ObjectId); err != nil {
		return FileObjectError{error: err}
	}

	return nil
}

// Rename a file/directory
// [objectId] and [fullPath] are optional parameters
// if [objectId] is not available then [fullPath] will be used to fetch the [objectId]
// dont leave both [objectId] and [fullPath] empty
// Tip: use [objectId] whenever possible to avoid traversing down the whole file tree to process and find the [objectId]
// rObjectId: objectId of the file/diectory
func RenameFile(dev *mtp.Device, storageId, objectId uint32, fullPath, newFileName string) (rObjectId uint32, error error) {
	exist, fi := FileExists(dev, storageId, objectId, fullPath)

	if !exist {
		return 0, InvalidPathError{error: fmt.Errorf("file not found: %s", fullPath)}
	}

	if err := dev.SetObjectPropValue(fi.ObjectId, mtp.OPC_ObjectFileName, &mtp.StringValue{Value: newFileName}); err != nil {
		switch v := err.(type) {
		case mtp.RCError:
			if v == 0x2002 {
				return fi.ObjectId, nil
			}
		}

		return 0, FileObjectError{error: err}
	}

	return fi.ObjectId, nil
}

// Transfer files from the local disk to the device
// sources: can be the list of files/directories that are to be sent to the device
// destination: fullPath to the destination directory
// preprocessFiles: if enabled, will fetch the total file size and count of the source. Use this will caution as it may take a few seconds to minutes to procress the files.
// rDestinationObjectId: objectId of [destination] directory
// rTotalFiles: total transferred files (directory count not included)
// rTotalSize: total size of the uploaded files
func UploadFiles(dev *mtp.Device, storageId uint32, sources []string, destination string, preprocessFiles bool, progressCb ProgressCb, preprocessCb PreprocessCb) (rDestinationObjectId uint32, rTotalFiles int64, rTotalSize int64, rError error) {
	_destination := fixSlash(destination)

	pInfo := ProgressInfo{
		FileInfo:          &FileInfo{},
		StartTime:         time.Now(),
		LatestSentTime:    time.Time{},
		Speed:             0,
		TotalFiles:        0,
		TotalDirectories:  0,
		FilesSent:         0,
		FilesSentProgress: 0,
		Current:           &TransferSizeInfo{},
		Bulk:              &TransferSizeInfo{},
	}

	// total number of files in this upload session
	var totalFiles int64 = 0

	// total number of files in this upload session
	var totalDirectories int64 = 0

	// total size of all the files combined in this upload session
	var totalSize int64 = 0

	// total number of files sent
	var bulkFilesSent int64 = 0

	// total size of data sent
	var bulkSizeSent int64 = 0

	if preprocessFiles {
		_totalFiles, _totalDirectories, _totalSize, err := walkLocalFiles(sources, func(fi *os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if err = preprocessCb(fi, nil); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			return 0, 0, 0, nil
		}

		totalFiles = _totalFiles
		totalDirectories = _totalDirectories
		totalSize = _totalSize
	}

	destParentId, err := MakeDirectory(dev, storageId, _destination)
	if err != nil {
		return 0, bulkFilesSent, bulkSizeSent, err
	}

	for _, source := range sources {
		_source := fixSlash(source)
		sourceParentPath := filepath.Dir(_source)

		destinationFilesDict := map[string]uint32{
			_destination: destParentId,
		}

		// walk through the source
		err = filepath.Walk(_source,
			func(path string, fInfo os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				name := fInfo.Name()

				// don't follow symlinks
				if isSymlinkLocal(fInfo) {
					return nil
				}

				// filter out disallowed files
				if isDisallowedFiles(name) {
					return nil
				}

				sourceFilePath := fixSlash(path)

				// map the local files path to the mtp files path
				destinationParentPath, destinationFilePath := mapSourcePathToDestinationPath(
					sourceFilePath, sourceParentPath, _destination,
				)

				size := fInfo.Size()
				isDir := fInfo.IsDir()

				// if the object is a directory then create a directory using [MakeDirectory] or [MakeDirectory]
				if isDir {
					// if the parent path exists within the [destinationFilesDict] then fetch the [parentId] (value) and make the destination directory
					if _, ok := destinationFilesDict[destinationParentPath]; ok {
						objId, err := MakeDirectory(dev, storageId, destinationFilePath)
						if err != nil {
							return err
						}

						// append the current objectId to [destinationFilesDict]
						destinationFilesDict[destinationFilePath] = objId

						// if the parent path DOES NOT exists within the [destinationFilesDict] create a new directory using costlier [MakeDirectory] method
						// this is a fallback situation
					} else {
						objId, err := MakeDirectory(dev, storageId, _destination)
						if err != nil {
							return err
						}

						// append the current objectId to [destinationFilesDict]
						destinationFilesDict[destinationFilePath] = objId
					}

					return nil
				}

				/// if the object is a file then create a file
				var fileParentId uint32

				_parentId, ok := destinationFilesDict[destinationParentPath]
				if ok {
					// if [destinationParentPath] exists within [destinationFilesDict] then use the value of the [destinationFilesDict] item as [parentId]
					fileParentId = _parentId

				} else {
					// if [destinationParentPath] DOES NOT exists within [destinationFilesDict] then create the parent directory using [MakeDirectory] and use the resulting objId as [parentId]
					objId, err := MakeDirectory(dev, storageId, destinationParentPath)

					if err != nil {
						return err
					}

					// append the current objectId to [destinationFilesDict]
					destinationFilesDict[destinationFilePath] = objId

					fileParentId = objId
				}

				// read the local file
				fileBuf, err := os.Open(sourceFilePath)
				if err != nil {
					return InvalidPathError{error: err}
				}
				defer fileBuf.Close()

				var compressedSize uint32

				// assign compressedSize of the file
				if size > 0xFFFFFFFF {
					compressedSize = 0xFFFFFFFF
				} else {
					compressedSize = uint32(size)
				}

				fObj := mtp.ObjectInfo{
					StorageID:        storageId,
					ObjectFormat:     mtp.OFC_Undefined,
					ParentObject:     fileParentId,
					Filename:         name,
					CompressedSize:   compressedSize,
					ModificationDate: time.Now(),
				}

				// keep track of [bulkFilesSent]
				bulkFilesSent += 1

				// keep track of [bulkSizeSent]
				bulkSizeSent += size

				pInfo.FileInfo = &FileInfo{
					Info:       &fObj,
					Size:       size,
					IsDir:      isDir,
					ModTime:    fObj.ModificationDate,
					Name:       fObj.Filename,
					FullPath:   destinationFilePath,
					ParentPath: destinationParentPath,
					Extension:  extension(fObj.Filename, isDir),
					ParentId:   fObj.ParentObject,
				}
				pInfo.LatestSentTime = time.Now()

				// create file
				var prevSentSize int64 = 0
				objId, err := handleMakeFile(
					dev, storageId, &fObj, &fInfo, fileBuf, true, func(total, sent int64, objId uint32) {
						pInfo.FileInfo.ObjectId = objId
						pInfo.Current = &TransferSizeInfo{
							Total:    total,
							Sent:     sent,
							Progress: Percent(float32(sent), float32(total)),
						}

						pInfo.Speed = transferRate(sent-prevSentSize, pInfo.LatestSentTime)

						_ = progressCb(&pInfo, nil)

						pInfo.LatestSentTime = time.Now()
						prevSentSize = sent
					},
				)

				if err != nil {
					return err
				}

				pInfo.TotalFiles = totalFiles
				pInfo.TotalDirectories = totalDirectories
				pInfo.FilesSent = bulkFilesSent
				pInfo.FilesSentProgress = Percent(float32(bulkFilesSent), float32(totalFiles))
				pInfo.Bulk = &TransferSizeInfo{
					Total:    totalSize,
					Sent:     bulkSizeSent,
					Progress: Percent(float32(bulkSizeSent), float32(totalSize)),
				}

				pInfo.FileInfo.ObjectId = objId
				if err := progressCb(&pInfo, nil); err != nil {
					return err
				}

				// append the current objectId to [destinationFilesDict]
				destinationFilesDict[destinationFilePath] = objId

				return nil
			},
		)

		if err != nil {
			switch err.(type) {
			case InvalidPathError:
				return destParentId, bulkFilesSent, bulkSizeSent, err

			case *os.PathError:
				if errors.Is(err, os.ErrPermission) {
					return destParentId, bulkFilesSent, bulkSizeSent, FilePermissionError{error: err}
				}

				if errors.Is(err, os.ErrNotExist) {
					return destParentId, bulkFilesSent, bulkSizeSent, InvalidPathError{error: err}
				}

				return destParentId, bulkFilesSent, bulkSizeSent, LocalFileError{error: err}
			default:
				return destParentId, bulkFilesSent, bulkSizeSent,
					FileTransferError{error: fmt.Errorf("an error occured while uploading files. %+v", err.Error())}
			}
		}
	}

	return destParentId, bulkFilesSent, bulkSizeSent, nil
}

// Transfer files from the device to the local disk
// sources: can be the list of files/directories that are to be sent to the local disk
// destination: fullPath to the destination directory
// rTotalFiles: total transferred files (directory count not included)
// rTotalSize: total size of the uploaded files
func DownloadFiles(dev *mtp.Device, storageId uint32, sources []string, destination string, cb TransferFilesCb) (rTotalFiles int, rTotalSize int64, rError error) {
	_destination := fixSlash(destination)

	downloadFi := TransferredFileInfo{
		StartTime:      time.Now(),
		LatestSentTime: time.Now(),
	}

	totalFiles := 0
	var totalSize int64 = 0

	for _, source := range sources {
		_source := fixSlash(source)
		sourceParentPath := filepath.Dir(_source)

		_, err := GetObjectFromPath(dev, storageId, _source)
		if err != nil {
			return totalFiles, totalSize, err
		}

		_, _, err = Walk(dev, storageId, _source, true, false,
			func(objectId uint32, fi *FileInfo, err error) error {
				destinationFileParentPath, destinationFilePath := mapSourcePathToDestinationPath(
					fi.FullPath, sourceParentPath, _destination,
				)

				if err != nil {
					return err
				}

				// filter out disallowed files
				if isDisallowedFiles(fi.Name) {
					return nil
				}

				// if the object is a directory then create a local directory
				if fi.IsDir {
					err := makeLocalDirectory(destinationFilePath)
					if err != nil {
						return err
					}

					return nil
				}

				/// if the object is a file then create one

				// if the local parent directory does not exists then create one
				if !fileExistsLocal(destinationFileParentPath) {
					err := makeLocalDirectory(destinationFileParentPath)
					if err != nil {
						return err
					}
				}

				// create the local file
				err = handleMakeLocalFile(dev, fi, destinationFilePath)
				if err != nil {
					return err
				}

				// keep track of [totalFiles]
				totalFiles += 1

				// keep track of [totalSize]
				totalSize += fi.Size

				downloadFi.FileInfo = fi
				downloadFi.FilesSent = totalFiles
				downloadFi.Speed = transferRateInMBs(fi.Size, downloadFi.LatestSentTime, downloadFi.Speed)
				downloadFi.LatestSentTime = time.Now()

				if err := cb(&downloadFi, nil); err != nil {
					return err
				}

				return nil
			})

		if err != nil {
			switch err.(type) {
			case InvalidPathError:
				return totalFiles, totalSize, err

			case *os.PathError:
				if errors.Is(err, os.ErrPermission) {
					return totalFiles, totalSize, FilePermissionError{error: err}
				}

				if errors.Is(err, os.ErrNotExist) {
					return totalFiles, totalSize, InvalidPathError{error: err}
				}

				return totalFiles, totalSize, LocalFileError{error: err}
			default:
				return totalFiles, totalSize,
					FileTransferError{error: fmt.Errorf("an error occured while downloading the files. %+v", err.Error())}
			}
		}
	}

	return totalFiles, totalSize, nil
}

func main() {}

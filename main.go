package main

import (
	"errors"
	"fmt"
	"github.com/ganeshrvel/go-mtpfs/mtp"
	"github.com/kr/pretty"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// todo: work on documentations
// todo: hotplug
// todo: information mode -> get total files, break the buf into smaller chunks and calculate the transfer rate

// initialize the mtp device
// returns mtp device
func Initialize(init Init) (*mtp.Device, error) {
	dev, err := mtp.SelectDevice("")

	if err != nil {
		return nil, MtpDetectFailedError{error: err}
	}

	dev.MTPDebug = init.debugMode
	dev.DataDebug = init.debugMode
	dev.USBDebug = init.debugMode

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

// fetch device info
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
			sid:  sids.Values[0],
			info: info,
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
// rDestinationObjectId: objectId of [destination] directory
// rTotalFiles: total transferred files (directory count not included)
// rTotalSize: total size of the uploaded files
func UploadFiles(dev *mtp.Device, storageId uint32, sources []string, destination string, cb TransferFilesCb) (rDestinationObjectId uint32, rTotalFiles int, rTotalSize int64, rError error) {
	_destination := fixSlash(destination)

	uploadFi := TransferredFileInfo{
		StartTime:      time.Now(),
		LatestSentTime: time.Now(),
	}

	totalFiles := 0
	var totalSize int64 = 0

	destParentId, err := MakeDirectory(dev, storageId, _destination)
	if err != nil {
		return 0, totalFiles, totalSize, err
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

				objId, err := handleMakeFile(
					dev, storageId, &fObj, &fInfo, fileBuf, true,
				)
				if err != nil {
					return err
				}

				// keep track of [totalFiles]
				totalFiles += 1

				// keep track of [totalSize]
				totalSize += size

				uploadFi.FileInfo = &FileInfo{
					Info:       &fObj,
					Size:       size,
					IsDir:      isDir,
					ModTime:    fObj.ModificationDate,
					Name:       fObj.Filename,
					FullPath:   destinationFilePath,
					ParentPath: destinationParentPath,
					Extension:  extension(fObj.Filename, isDir),
					ParentId:   fObj.ParentObject,
					ObjectId:   objId,
				}

				uploadFi.FilesSent = totalFiles
				uploadFi.Speed = transferRateInMBs(size, uploadFi.LatestSentTime, uploadFi.Speed)
				uploadFi.LatestSentTime = time.Now()

				if err := cb(&uploadFi, nil); err != nil {
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
				return destParentId, totalFiles, totalSize, err

			case *os.PathError:
				if errors.Is(err, os.ErrPermission) {
					return destParentId, totalFiles, totalSize, FilePermissionError{error: err}
				}

				if errors.Is(err, os.ErrNotExist) {
					return destParentId, totalFiles, totalSize, InvalidPathError{error: err}
				}

				return destParentId, totalFiles, totalSize, LocalFileError{error: err}
			default:
				return destParentId, totalFiles, totalSize,
					FileTransferError{error: fmt.Errorf("an error occured while uploading files. %+v", err.Error())}
			}
		}
	}

	return destParentId, totalFiles, totalSize, nil
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

func main() {
	dev, err := Initialize(Init{debugMode: false})

	if err != nil {
		log.Panic(err)
	}

	_, err = FetchDeviceInfo(dev)
	if err != nil {
		log.Panic(err)
	}

	storages, err := FetchStorages(dev)
	if err != nil {
		log.Panic(err)
	}

	sid := storages[0].sid
	pretty.Println("storage id: ", sid)

	//totalFiles, err := dev.GetNumObjects(sid, mtp.GOH_ALL_ASSOCS, ParentObjectId)
	//if err != nil {
	//	log.Panic(err)
	//}
	//
	//pretty.Println(int64(totalFiles))
	//

	/*objectId, totalFiles, err := Walk(dev, sid, 0, "/mtp-test-files/mock_dir1", true, func(objectId uint32, fi *FileInfo) {
		pretty.Println("objectId is: ", objectId)
	})

	if err != nil {
		log.Panic(err)
	}

	pretty.Println("totalFiles: ", totalFiles)
	pretty.Println("objectId: ", objectId)*/

	////MakeDirectory
	//objectId, err := MakeDirectory(dev, sid, ParentObjectId, "/", "name")
	//if err != nil {
	//	log.Panic(err)
	//}
	//
	//pretty.Println(objectId)

	//GetObjectFromPath
	//fileObj, err := GetObjectFromPath(dev, sid, "/tests/s")
	//if err != nil {
	//	log.Panic(err)
	//}
	//
	//pretty.Println("======\n")
	//pretty.Println(fileObj)
	//

	// FileExists
	//exists := FileExists(dev, sid, 0, "/tests/test.txt")
	//
	//pretty.Println("======\n")
	//pretty.Println("Does File exists:", exists)

	///DeleteFile
	//err = DeleteFile(dev, sid, 0, "/mtp-test-files/temp_dir/this is a test")
	//if err != nil {
	//	log.Panic(err)
	//}

	////RenameFile
	//objId, err := RenameFile(dev, sid, 0, "/mtp-test-files/temp_dir/b.txt", "b.txt")
	//if err != nil {
	//	log.Panic(err)
	//}
	//pretty.Println(objId)

	//UploadFiles
	//uploadFile2 := getTestMocksAsset("test-large-file")
	//start := time.Now()
	//
	//totalFiles, totalSize, err := DownloadFiles(dev, sid,
	//	[]string{sourceFile1}, downloadFile,
	//	func(downloadFi *TransferredFileInfo, err error) error {
	//		fmt.Printf("Current filepath: %s\n", downloadFi.FileInfo.FullPath)
	//		fmt.Printf("%f MB/s\n", downloadFi.Speed)
	//
	//		return nil
	//	},
	//)
	//if err != nil {
	//	log.Panic(err)
	//}
	//
	//pretty.Println(totalFiles)
	//pretty.Println(totalSize)
	//pretty.Println("time elapsed: ", time.Since(start).Seconds())

	Dispose(dev)
}

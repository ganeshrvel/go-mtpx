package main

import (
	"fmt"
	mtp "github.com/ganeshrvel/go-mtpfs/mtp"
	"github.com/kr/pretty"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

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

// create a new directory using [parentPath] and [name]
// this is non recursive
// if the parentpath does not exists then an error will be thrown
// [parentPath] -> path to the parent directory
// [name] -> directory name
func MakeDirectory(dev *mtp.Device, storageId, parentId uint32, parentPath, name string) (rObjectId uint32, rError error) {
	if name == "" {
		return 0, InvalidPathError{error: fmt.Errorf("invalid path: %s. The name cannot be empty", parentPath)}
	}

	// check if the path exists in the MTP device
	fi, err := GetObjectFromParentIdAndFilename(dev, storageId, parentId, name)

	if err == nil {
		// the object exists and if it's a file then throw an error
		if !fi.IsDir {
			return 0, InvalidPathError{error: fmt.Errorf("invalid path: %s. The object is not a directory", parentPath)}
		}

		return fi.ObjectId, nil
	}

	// check if the parent exists
	parentFi, err := GetObjectFromObjectIdOrPath(dev, storageId, parentId, parentPath)

	if err != nil {
		return 0, err
	}

	// if the parent object is a file then throw an error
	if !parentFi.IsDir {
		return 0, InvalidPathError{error: fmt.Errorf("invalid path: %s. The object is not a directory", parentPath)}
	}

	// create the directory
	return handleMakeDirectory(dev, storageId, parentFi.ObjectId, name)
}

// create a new directory recursively using [fullPath]
// The path will be created if it does not exists
func MakeDirectoryRecursive(dev *mtp.Device, storageId uint32, fullPath string) (rObjectId uint32, rError error) {
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
// [objectId] and [fullPath] are optional parameters
// if [objectId] is not available then [fullPath] will be used to fetch the [objectId]
// dont leave both [objectId] and [fullPath] empty
// Tips: use [objectId] whenever possible to avoid traversing down the whole file tree to process and find the [objectId]
// returns total number of objects
func WalkDirectory(dev *mtp.Device, storageId, objectId uint32, fullPath string, recursive bool, cb WalkDirectoryCb) (rObjectId uint32, rTotalFiles int, rError error) {
	// fetch the objectId from [objectId] and/or [fullPath] parameters
	fi, err := GetObjectFromObjectIdOrPath(dev, storageId, objectId, fullPath)
	if err != nil {
		return 0, 0, err
	}

	if !fi.IsDir {
		return 0, 0, InvalidPathError{error: fmt.Errorf("invalid path: %s. The object is not a directory", fullPath)}
	}

	totalFiles, err := proccessWalkDirectory(dev, storageId, fi.ObjectId, fullPath, recursive, cb)
	if err != nil {
		return 0, 0, err
	}

	return fi.ObjectId, totalFiles, nil
}

// Send local files to the device
// sources: can be the list of files/directories that are to be sent to the device
// destination: fullPath to the destination directory
func UploadFiles(dev *mtp.Device, storageId uint32, sources []string, destination string, cb UploadFilesCb) (rDestinationObjectId uint32, rTotalFiles int, rTotalSize int64, rError error) {
	_destination := fixSlash(destination)

	uploadFi := UploadFileInfo{
		StartTime:      time.Now(),
		LatestSentTime: time.Now(),
	}

	totalFiles := 0
	var totalSize int64 = 0

	destParentId, err := MakeDirectoryRecursive(dev, storageId, _destination)
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
				destinationParentPath, destinationFilePath := mapLocalPathToMtpPath(
					sourceFilePath, sourceParentPath, _destination,
				)

				size := fInfo.Size()
				isDir := fInfo.IsDir()

				// if the object is a directory then create a directory using [MakeDirectory] or [MakeDirectoryRecursive]
				if isDir {
					// if the parent path exists within the [destinationFilesDict] then fetch the [parentId] (value) and make the destination directory
					if parentId, ok := destinationFilesDict[destinationParentPath]; ok {
						objId, err := MakeDirectory(
							dev, storageId, parentId, destinationParentPath, name,
						)
						if err != nil {
							return err
						}

						// append the current objectId to [destinationFilesDict]
						destinationFilesDict[destinationFilePath] = objId

						// if the parent path DOES NOT exists within the [destinationFilesDict] create a new directory using costlier [MakeDirectoryRecursive] method
						// this is a fallback situation
					} else {
						objId, err := MakeDirectoryRecursive(dev, storageId, _destination)
						if err != nil {
							return err
						}

						// append the current objectId to [destinationFilesDict]
						destinationFilesDict[destinationFilePath] = objId
					}

					return nil
				}

				// if the object is a file then create a file

				// keep track of [totalFiles]
				totalFiles += 1

				// keep track of [totalSize]
				totalSize += size

				var fileParentId uint32

				_parentId, ok := destinationFilesDict[destinationParentPath]
				if ok {
					// if [destinationParentPath] exists within [destinationFilesDict] then use the value of the [destinationFilesDict] item as [parentId]
					fileParentId = _parentId

				} else {
					// if [destinationParentPath] DOES NOT exists within [destinationFilesDict] then create the parent directory using [MakeDirectoryRecursive] and use the resulting objId as [parentId]
					objId, err := MakeDirectoryRecursive(dev, storageId, destinationParentPath)

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

				cb(&uploadFi)

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
				return destParentId, totalFiles, totalSize, InvalidPathError{error: err}

			default:
				return destParentId, totalFiles, totalSize,
					UploadFileError{error: fmt.Errorf("an error occured while uploading files. %+v", err)}
			}
		}
	}

	return destParentId, totalFiles, totalSize, nil
}

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

func RenameFile(dev *mtp.Device, storageId, objectId uint32, fullPath, newFileName string) (rObjectId uint32, error error) {
	exist, fi := FileExists(dev, storageId, objectId, fullPath)

	if !exist {
		return 0, InvalidPathError{error: fmt.Errorf("file not found: %s", fullPath)}
	}

	if err := dev.SetObjectPropValue(fi.ObjectId, mtp.OPC_ObjectFileName, &mtp.StringValue{Value: newFileName}); err != nil {
		return 0, FileObjectError{error: err}
	}

	return fi.ObjectId, nil
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

	/*objectId, totalFiles, err := WalkDirectory(dev, sid, 0, "/mtp-test-files/mock_dir1", true, func(objectId uint32, fi *FileInfo) {
		pretty.Println("objectId is: ", objectId)
	})

	if err != nil {
		log.Panic(err)
	}

	pretty.Println("totalFiles: ", totalFiles)
	pretty.Println("objectId: ", objectId)*/

	///MakeDirectory
	//objectId, err := MakeDirectory(dev, sid, "/", "name")
	//if err != nil {
	//	log.Panic(err)
	//}
	//
	//pretty.Println(objectId)

	///ListDirectory
	//files, err := ListDirectory(dev, sid, 0, "/")
	//if err != nil {
	//	log.Panic(err)
	//}
	//
	//pretty.Println("Listing directory test: ", files)

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

	//////RenameFile
	//objId, err := RenameFile(dev, sid, 0, "/mtp-test-files/temp_dir/b.txt", "b.txt")
	//if err != nil {
	//	log.Panic(err)
	//}
	//pretty.Println(objId)

	//UploadFiles
	uploadFile := getTestMocksAsset("mock_dir1")
	uploadFile2 := getTestMocksAsset("mock_dir2")
	start := time.Now()

	objId, totalFiles, totalSize, err := UploadFiles(dev, sid,
		[]string{uploadFile, uploadFile2, uploadFile}, "/mtp-test-files/temp_dir/test_UploadFiles",
		func(uploadFi *UploadFileInfo) {
			fmt.Printf("Current filepath: %s\n", uploadFi.FileInfo.FullPath)
			fmt.Printf("%x MB/s\n", uploadFi.Speed)
		},
	)
	if err != nil {

		log.Panic(err)
	}

	pretty.Println(objId)
	pretty.Println(totalFiles)
	pretty.Println(totalSize)
	pretty.Println("time elapsed: ", time.Since(start).Seconds())

	Dispose(dev)
}

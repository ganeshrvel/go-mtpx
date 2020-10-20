package main

import (
	"fmt"
	mtp "github.com/ganeshrvel/go-mtpfs/mtp"
	"github.com/kr/pretty"
	"log"
	"strings"
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
		if err := dev.GetStorageInfo(uint32(sid), &info); err != nil {
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
func MakeDirectory(dev *mtp.Device, storageId uint32, parentPath, name string) (rObjectId uint32, rError error) {
	if name == "" {
		return 0, InvalidPathError{error: fmt.Errorf("invalid path: %s. The name cannot be empty", parentPath)}
	}

	parentId, isDir, err := GetObjectFromPath(dev, storageId, parentPath)

	if err != nil {
		return 0, err
	}

	// if the object exists but if it's a file then throw an error
	if !isDir {
		return 0, InvalidPathError{error: fmt.Errorf("invalid path: %s. The object is not a directory", parentPath)}
	}

	fullPath := getFullPath(parentPath, name)

	exist, isDir, _objectId := FileExists(dev, storageId, fullPath)

	if exist {
		// if the object exists but if it's a file then throw an error
		if !isDir {
			return 0, InvalidPathError{error: fmt.Errorf("invalid path: %s. The object is not a directory", parentPath)}
		}

		return _objectId, nil
	}

	return handleMakeDirectory(dev, storageId, parentId, name)
}

// create a new directory recursively using [filePath]
// The path will be created if it does not exists
func MakeDirectoryRecursive(dev *mtp.Device, storageId uint32, filePath string) (rObjectId uint32, rError error) {
	_filePath := fixSlash(filePath)

	if _filePath == PathSep {
		return ParentObjectId, nil
	}

	splittedFilePath := strings.Split(_filePath, PathSep)

	var parentId = uint32(ParentObjectId)
	const skipIndex = 1

	for _, fName := range splittedFilePath[skipIndex:] {
		// fetch the parent object and
		_parentId, isDir, err := GetObjectFromParentIdAndFilename(dev, storageId, parentId, fName)

		if err != nil {
			switch err.(type) {
			case FileNotFoundError:
				// if object does not exists then create a new directory
				_newObjectId, err := handleMakeDirectory(dev, storageId, parentId, fName)
				if err != nil {
					return 0, err
				}

				parentId = _newObjectId

				continue
			default:
				return 0, err
			}
		}

		// if the object exists but if it's a file then throw an error
		if !isDir {
			return 0, InvalidPathError{error: fmt.Errorf("invalid path: %s. The object is not a directory", fName)}
		}

		parentId = _parentId
	}

	return parentId, nil
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
	objId, err := GetObjectFromObjectIdOrPath(dev, storageId, objectId, fullPath)
	if err != nil {
		return objId, 0, err
	}

	fi, err := GetObjectFromObjectId(dev, objId, fullPath)
	if err != nil {
		return objId, 0, err
	}

	cb(objId, fi)

	totalFiles, err := proccessWalkDirectory(dev, storageId, objId, fullPath, recursive, cb)
	if err != nil {
		return objId, 0, err
	}

	return objId, totalFiles, nil
}

func main() {
	dev, err := Initialize(Init{})

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

	objectId, totalFiles, err := WalkDirectory(dev, sid, 0, "/mtp-test-files/mock_dir1", true, func(objectId uint32, fi *FileInfo) {
		pretty.Println("objectId is: ", objectId)
	})

	if err != nil {
		log.Panic(err)
	}

	pretty.Println("totalFiles: ", totalFiles)
	pretty.Println("objectId: ", objectId)

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
	//exists := FileExists(dev, sid, "/tests/test.txt")
	//
	//pretty.Println("======\n")
	//pretty.Println("Does File exists:", exists)
	Dispose(dev)
}

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

	for sid := range sids.Values {
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

// fetch file info using object id
func FetchFile(dev *mtp.Device, objectId uint32, parentPath string) (*FileInfo, error) {
	obj := mtp.ObjectInfo{}
	if err := dev.GetObjectInfo(objectId, &obj); err != nil {
		return nil, FileObjectError{error: err}
	}

	size, _ := GetFileSize(dev, &obj, objectId)
	isDir := isObjectADir(&obj)
	filename := obj.Filename
	_parentPath := fixSlash(parentPath)
	fullPath := getFullPath(_parentPath, filename)

	return &FileInfo{
		Info:       &obj,
		Size:       size,
		IsDir:      isDir,
		ModTime:    obj.ModificationDate,
		Name:       obj.Filename,
		FullPath:   fullPath,
		ParentPath: _parentPath,
		Extension:  extension(obj.Filename, isDir),
		ParentId:   obj.ParentObject,
		ObjectId:   objectId,
	}, nil
}

// List the contents in a directory
// [objectId] and [parentPath] are optional parameters
// if [objectId] is not available then parentPath is used to fetch objectId
// dont leave both [objectId] and [parentPath] empty
// Tips: use [objectId] whenever possible to avoid traversing down the file tree
func ListDirectory(dev *mtp.Device, storageId, objectId uint32, parentPath string) (*[]FileInfo, error) {
	_objectId := objectId

	// if objectId is not available then fetch the objectId from parentPath
	if _objectId == 0 {
		objId, isDir, err := GetPathObject(dev, storageId, parentPath)

		if err != nil {
			return nil, err
		}

		// if the object is not a directory throw an error
		if !isDir {
			return nil, InvalidPathError{error: fmt.Errorf("invalid path: %s. The object is not a directory", parentPath)}
		}

		_objectId = objId
	} else {
		if _objectId != ParentObjectId {
			f, err := FetchFile(dev, _objectId, parentPath)

			if err != nil {
				return nil, err
			}

			// if the object is not a directory throw an error
			if !f.IsDir {
				return nil, InvalidPathError{error: fmt.Errorf("invalid path: %s. The object is not a directory", parentPath)}
			}
		}
	}

	handles := mtp.Uint32Array{}
	if err := dev.GetObjectHandles(storageId, mtp.GOH_ALL_ASSOCS, _objectId, &handles); err != nil {
		return nil, ListDirectoryError{error: err}
	}

	var fileInfoList []FileInfo

	for _, objectId := range handles.Values {
		fi, err := FetchFile(dev, objectId, parentPath)

		if err != nil {
			continue
		}

		fileInfoList = append(fileInfoList, *fi)
	}

	return &fileInfoList, nil
}

func MakeDirectory(dev *mtp.Device, storageId uint32, parentPath, filename string) (objectId uint32, error error) {
	if filename == "" {
		return 0, InvalidPathError{error: fmt.Errorf("invalid path: %s. The filename cannot be empty", parentPath)}
	}

	parentId, isDir, err := GetPathObject(dev, storageId, parentPath)

	if err != nil {
		return 0, err
	}

	// if the object exists but if it's a file then throw an error
	if !isDir {
		return 0, InvalidPathError{error: fmt.Errorf("invalid path: %s. The object is not a directory", parentPath)}
	}

	fullPath := getFullPath(parentPath, filename)

	exist, isDir, _objectId := FileExists(dev, storageId, fullPath)

	if exist {
		// if the object exists but if it's a file then throw an error
		if !isDir {
			return 0, InvalidPathError{error: fmt.Errorf("invalid path: %s. The object is not a directory", parentPath)}
		}

		return _objectId, nil
	}

	return handleMakeDirectory(dev, storageId, parentId, filename)
}

func MakeDirectoryRecursive(dev *mtp.Device, storageId uint32, filePath string) (objectId uint32, error error) {
	_filePath := fixSlash(filePath)

	if _filePath == PathSep {
		return ParentObjectId, nil
	}

	splittedFilePath := strings.Split(_filePath, PathSep)

	var parentId = uint32(ParentObjectId)
	const skipIndex = 1

	for _, fName := range splittedFilePath[skipIndex:] {
		// fetch the parent object and
		_parentId, isDir, err := GetParentObject(dev, storageId, parentId, fName)

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

	objectId, err := MakeDirectory(dev, sid, "/", "name")
	if err != nil {
		log.Panic(err)
	}

	pretty.Println(objectId)

	/*	files, err := ListDirectory(dev, sid, 0, "/test/")
		if err != nil {
			log.Panic(err)
		}

		pretty.Println(files)*/

	/*
		fileObj, err := GetPathObject(dev, sid, "/tests/s")
		if err != nil {
			log.Panic(err)
		}

		pretty.Println("======\n")
		pretty.Println(fileObj)
	*/

	/*exists := FileExists(dev, sid, "/tests/test.txt")

	pretty.Println("======\n")
	pretty.Println("Does File exists:", exists)*/
	Dispose(dev)
}

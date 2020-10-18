package main

import (
	"fmt"
	"github.com/ganeshrvel/go-mtpfs/mtp"
	"strings"
)

func GetFileSize(dev *mtp.Device, obj *mtp.ObjectInfo, objectId uint32) (int64, error) {
	var size int64
	if obj.CompressedSize == 0xffffffff {
		var val mtp.Uint64Value
		if err := dev.GetObjectPropValue(objectId, mtp.OPC_ObjectSize, &val); err != nil {
			return 0, FileObjectError{
				fmt.Errorf("GetObjectPropValue handle %d failed: %v", objectId, err),
			}
		}

		size = int64(val.Value)
	} else {
		size = int64(obj.CompressedSize)
	}

	return size, nil
}

func GetPathObject(dev *mtp.Device, storageId uint32, filePath string) (objectId uint32, isDir bool, error error) {
	_filePath := fixSlash(filePath)

	splittedFilePath := strings.Split(_filePath, PathSep)

	if _filePath == PathSep {
		return ParentObjectId, true, nil
	}

	var parentId = uint32(ParentObjectId)
	isDir = true
	var resultCount = 0
	const skipIndex = 1

	for i, fName := range splittedFilePath[skipIndex:] {
		objectId, _isDir, err := GetParentObject(dev, storageId, parentId, fName)

		if err != nil {
			switch err.(type) {
			case FileNotFoundError:
				return 0, false, InvalidPathError{
					error: fmt.Errorf("path not found: %s\nreason: %v", filePath, err),
				}

			default:
				return 0, false, err
			}
		}

		if !_isDir && indexExists(splittedFilePath, i+1+skipIndex) {
			return 0, false, InvalidPathError{error: fmt.Errorf("path not found: %s", filePath)}
		}

		parentId = objectId
		isDir = _isDir
		resultCount += 1
	}

	if resultCount < 1 {
		return 0, false, InvalidPathError{error: fmt.Errorf("file not found: %s", filePath)}
	}

	return parentId, isDir, nil
}

func GetParentObject(dev *mtp.Device, storageId uint32, parentId uint32, filename string) (objectId uint32, isDir bool, error error) {
	handles := mtp.Uint32Array{}
	if err := dev.GetObjectHandles(storageId, mtp.GOH_ALL_ASSOCS, parentId, &handles); err != nil {
		return 0, false, FileObjectError{error: err}
	}

	for _, objectId := range handles.Values {
		obj := mtp.ObjectInfo{}
		if err := dev.GetObjectInfo(objectId, &obj); err != nil {
			return 0, false, FileObjectError{error: err}
		}

		// return the current objectId if the filename == obj.Filename
		if obj.Filename == filename {
			return objectId, isObjectADir(&obj), nil
		}
	}

	return 0, false, FileNotFoundError{error: fmt.Errorf("file not found: %s", filename)}
}

func FileExists(dev *mtp.Device, storageId uint32, filePath string) (exists bool, objectId uint32) {
	objectId, _, err := GetPathObject(dev, storageId, filePath)

	if err != nil {
		return false, objectId
	}

	return true, objectId
}

func isObjectADir(obj *mtp.ObjectInfo) bool {
	return obj.ObjectFormat == mtp.OFC_Association
}

package main

import (
	"fmt"
	"github.com/ganeshrvel/go-mtpfs/mtp"
	"strings"
	"time"
)

// fetch the file size of the object
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

// fetch an object using [objectId]
// [parentPath] is required to keep track of the [fullPath] of the object
func GetObjectFromObjectId(dev *mtp.Device, objectId uint32, parentPath string) (*FileInfo, error) {
	obj := mtp.ObjectInfo{}

	// if the [objectId] is root then return the basic root directory information
	if objectId == ParentObjectId {
		return &FileInfo{
			Size:     0,
			IsDir:    true,
			FullPath: "/",
			ObjectId: ParentObjectId,
			Info:     &mtp.ObjectInfo{},
		}, nil
	}

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

// fetch the object using [parentId] and [filename]
// it matches the [filename] to the list of files in the directory
// Since the [parentPath] is unavailable here the [fullPath] property of the resulting object [FileInfo] may not be valid.
func GetObjectFromParentIdAndFilename(dev *mtp.Device, storageId uint32, parentId uint32, filename string) (*FileInfo, error) {
	handles := mtp.Uint32Array{}
	if err := dev.GetObjectHandles(storageId, mtp.GOH_ALL_ASSOCS, parentId, &handles); err != nil {
		return nil, FileObjectError{error: err}
	}

	for _, objectId := range handles.Values {
		fi, err := GetObjectFromObjectId(dev, objectId, "")

		if err != nil {
			return nil, FileObjectError{error: err}
		}

		// return the current objectId if the filename == fi.Name
		if fi.Name == filename {
			return fi, nil
		}
	}

	return nil, FileNotFoundError{error: fmt.Errorf("file not found: %s", filename)}
}

// fetch the object information using [fullPath]
// Since the [parentPath] is unavailable here the [fullPath] property of the resulting object [FileInfo] may not be valid.
func GetObjectFromPath(dev *mtp.Device, storageId uint32, fullPath string) (*FileInfo, error) {
	_filePath := fixSlash(fullPath)

	if _filePath == PathSep {
		return GetObjectFromObjectId(dev, ParentObjectId, "")
	}

	splittedFilePath := strings.Split(_filePath, PathSep)

	var objectId = uint32(ParentObjectId)
	var resultCount = 0
	var fi *FileInfo
	const skipIndex = 1

	for i, fName := range splittedFilePath[skipIndex:] {
		_fi, err := GetObjectFromParentIdAndFilename(dev, storageId, objectId, fName)

		if err != nil {
			switch err.(type) {
			case FileNotFoundError:
				return nil, InvalidPathError{
					error: fmt.Errorf("path not found: %s\nreason: %v", fullPath, err),
				}

			default:
				return nil, err
			}
		}

		if !_fi.IsDir && indexExists(splittedFilePath, i+1+skipIndex) {
			return nil, InvalidPathError{error: fmt.Errorf("path not found: %s", fullPath)}
		}

		// updating [fi] to current [_fi]
		fi = _fi

		// updating the [objectId] to follow the nested directory
		objectId = _fi.ObjectId

		// keeping a tab on total results
		resultCount += 1
	}

	if resultCount < 1 {
		return nil, InvalidPathError{error: fmt.Errorf("file not found: %s", fullPath)}
	}

	return fi, nil
}

// fetch an object using [objectId] and/or [fullPath]
// Since the [parentPath] is unavailable here the [fullPath] property of the resulting object [FileInfo] may not be valid.
func GetObjectFromObjectIdOrPath(dev *mtp.Device, storageId, objectId uint32, fullPath string) (*FileInfo, error) {
	_objectId := objectId

	if _objectId == 0 && fullPath == "" {
		return nil, InvalidPathError{error: fmt.Errorf("invalid path: %s. both objectId and fullPath cannot be empty", fullPath)}
	}

	// if objectId is not available then fetch the objectId from fullPath
	if _objectId == 0 {
		fp, err := GetObjectFromPath(dev, storageId, fullPath)

		if err != nil {
			return nil, err
		}

		return fp, nil
	}

	fo, err := GetObjectFromObjectId(dev, _objectId, fullPath)
	if err != nil {
		return nil, err
	}

	return fo, nil
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

// check if the object is a directory
func isObjectADir(obj *mtp.ObjectInfo) bool {
	return obj.ObjectFormat == mtp.OFC_Association
}

// helper function to create a directory
func handleMakeDirectory(dev *mtp.Device, storageId, parentId uint32, filename string) (rObjectId uint32, rError error) {
	send := mtp.ObjectInfo{
		StorageID:        storageId,
		ObjectFormat:     mtp.OFC_Association,
		ParentObject:     parentId,
		Filename:         filename,
		CompressedSize:   uint32(0),
		ModificationDate: time.Now(),
	}

	_, _, objId, err := dev.SendObjectInfo(storageId, parentId, &send)
	if err != nil {
		return 0, SendObjectError{error: err}
	}

	return objId, nil
}

// helper function to create a file
func handleMakeFile(dev *mtp.Device, storageId uint32, obj *mtp.ObjectInfo) (rObjectId uint32, rError error) {
	_, _, objId, err := dev.SendObjectInfo(storageId, obj.ParentObject, obj)
	if err != nil {
		return 0, SendObjectError{error: err}
	}

	return objId, nil
}

// helper function to fetch the contents inside a directory
// use [recursive] to fetch the whole nested tree
// [objectId] and [fullPath] are optional parameters
// if [objectId] is not available then [fullPath] will be used to fetch the [objectId]
// dont leave both [objectId] and [fullPath] empty
// Tips: use [objectId] whenever possible to avoid traversing down the whole file tree to process and find the [objectId]
// returns total number of objects
func proccessWalkDirectory(dev *mtp.Device, storageId, objectId uint32, fullPath string, recursive bool, cb WalkDirectoryCb) (rTotalFiles int, rError error) {
	fi, err := GetObjectFromObjectIdOrPath(dev, storageId, objectId, fullPath)

	if err != nil {
		return 0, err
	}

	handles := mtp.Uint32Array{}
	if err := dev.GetObjectHandles(storageId, mtp.GOH_ALL_ASSOCS, fi.ObjectId, &handles); err != nil {
		return 0, ListDirectoryError{error: err}
	}

	totalFiles := 0

	for _, objId := range handles.Values {
		fi, err := GetObjectFromObjectId(dev, objId, fullPath)
		if err != nil {
			continue
		}

		totalFiles += 1

		cb(objId, fi)

		// don't traverse down the tree if [recursive] is false
		if !recursive {
			continue
		}

		// don't traverse down the tree if the object is not a directory
		if !fi.IsDir {
			continue
		}

		_totalFiles, err := proccessWalkDirectory(
			dev, storageId, objId, fi.FullPath, recursive, cb,
		)
		if err != nil {
			continue
		}

		totalFiles += _totalFiles
	}

	return totalFiles, nil
}

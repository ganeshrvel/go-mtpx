package main

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"math/rand"
	"testing"
)

func TestMakeDirectory(t *testing.T) {
	dev, err := Initialize(Init{})
	if err != nil {
		log.Panic(err)
	}

	storages, err := FetchStorages(dev)
	if err != nil {
		log.Panic(err)
	}

	sid := storages[0].sid

	var _objectId uint32
	var _objectId2 uint32
	Convey("Creating a new dir | MakeDirectory", t, func() {
		// test the directory '/mtp-test-files/temp_dir/test-MakeDirectory'
		objectId, err := MakeDirectory(dev, sid, 0, "/mtp-test-files/temp_dir", "test-MakeDirectory")

		_objectId = objectId

		So(err, ShouldBeNil)
		So(objectId, ShouldBeGreaterThan, 0)
	})

	Convey("Creating a new dir | using parentId | MakeDirectory", t, func() {
		// test the directory '/mtp-test-files/temp_dir/test-MakeDirectoryUsingParentId'
		parentId, isDir, err := GetObjectFromPath(dev, sid, "/mtp-test-files/temp_dir")
		So(err, ShouldBeNil)
		So(parentId, ShouldBeGreaterThan, 0)
		So(isDir, ShouldEqual, true)

		objectId, err := MakeDirectory(dev, sid, parentId, "", "test-MakeDirectoryUsingParentId")

		_objectId2 = objectId

		So(err, ShouldBeNil)
		So(objectId, ShouldBeGreaterThan, 0)
	})

	Convey("Testing MakeDirectory for an existing directory | MakeDirectory", t, func() {
		// test the directory '/mtp-test-files/temp_dir/test-MakeDirectory'
		objectId, err := MakeDirectory(dev, sid, 0, "/mtp-test-files/temp_dir", "test-MakeDirectory")

		So(err, ShouldBeNil)
		So(objectId, ShouldEqual, _objectId)
	})

	Convey("Testing MakeDirectory for an existing directory | using parentId | MakeDirectory", t, func() {
		// test the directory '/mtp-test-files/temp_dir/test-MakeDirectoryUsingParentId'
		parentId, isDir, err := GetObjectFromPath(dev, sid, "/mtp-test-files/temp_dir")
		So(err, ShouldBeNil)
		So(parentId, ShouldBeGreaterThan, 0)
		So(isDir, ShouldEqual, true)

		objectId, err := MakeDirectory(dev, sid, parentId, "", "test-MakeDirectoryUsingParentId")

		So(err, ShouldBeNil)
		So(objectId, ShouldEqual, _objectId2)
	})

	Convey("Creating a new random dir | fullpath | MakeDirectory", t, func() {
		// test the directory '/mtp-test-files/temp_dir/test-MakeDirectory/{random}'
		filename := fmt.Sprintf("%x", rand.Int31())

		objectId, err := MakeDirectory(dev, sid, 0, "/mtp-test-files/temp_dir/test-MakeDirectory", filename)

		So(err, ShouldBeNil)
		So(objectId, ShouldBeGreaterThan, 0)

		exists, isDir, _existingObjectId := FileExists(dev, sid, 0, getFullPath("/mtp-test-files/temp_dir/test-MakeDirectory", filename))

		So(err, ShouldBeNil)
		So(exists, ShouldEqual, true)
		So(_existingObjectId, ShouldEqual, objectId)
		So(isDir, ShouldEqual, true)
	})

	Convey("Creating a new random dir | parentId | MakeDirectory", t, func() {
		// test the directory '/mtp-test-files/temp_dir/test-MakeDirectory/{random}'
		filename := fmt.Sprintf("%x", rand.Int31())
		parentId, isDir, err := GetObjectFromPath(dev, sid, "/mtp-test-files/temp_dir/test-MakeDirectoryUsingParentId")
		So(err, ShouldBeNil)
		So(parentId, ShouldBeGreaterThan, 0)
		So(isDir, ShouldEqual, true)

		objectId, err := MakeDirectory(dev, sid, parentId, "", filename)

		So(err, ShouldBeNil)
		So(objectId, ShouldBeGreaterThan, 0)

		exists, isDir, _existingObjectId := FileExists(dev, sid, 0, getFullPath("/mtp-test-files/temp_dir/test-MakeDirectoryUsingParentId", filename))

		So(err, ShouldBeNil)
		So(exists, ShouldEqual, true)
		So(_existingObjectId, ShouldEqual, objectId)
		So(isDir, ShouldEqual, true)
	})

	Convey("invalid path | MakeDirectory | fullPath | It should throw an error", t, func() {
		// test the directory '/fake/test'
		objectId, err := MakeDirectory(dev, sid, 0, "fake", "test")

		So(err, ShouldBeError)
		So(objectId, ShouldEqual, 0)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
	})

	Convey("invalid path | MakeDirectory | parentId | It should throw an error", t, func() {
		// test the directory '/fake/test'
		objectId, err := MakeDirectory(dev, sid, 1234561234, "/mtp-test-files", "test")

		So(err, ShouldBeError)
		So(objectId, ShouldEqual, 0)
		So(err, ShouldHaveSameTypeAs, FileObjectError{})
	})

	Convey("empty folder name | MakeDirectory | It should throw an error", t, func() {
		// test the directory '/'
		objectId, err := MakeDirectory(dev, sid, 0, "fake", "")

		So(err, ShouldBeError)
		So(objectId, ShouldEqual, 0)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
	})

	Convey("filename in the path | MakeDirectory | It should throw an error", t, func() {
		// test the directory '/mtp-test-files/a.txt'
		objectId, err := MakeDirectory(dev, sid, 0, "/mtp-test-files/", "a.txt")

		So(err, ShouldBeError)
		So(objectId, ShouldEqual, 0)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
	})

	Dispose(dev)
}

func TestMakeDirectoryRecursive(t *testing.T) {
	dev, err := Initialize(Init{})
	if err != nil {
		log.Panic(err)
	}

	storages, err := FetchStorages(dev)
	if err != nil {
		log.Panic(err)
	}

	sid := storages[0].sid

	var _objectId uint32
	Convey("Creating a new dir | MakeDirectoryRecursive", t, func() {
		// test the directory '/mtp-test-files/temp_dir/test-MakeDirectoryRecursive'
		objectId, err := MakeDirectoryRecursive(dev, sid, "/mtp-test-files/temp_dir/test-MakeDirectoryRecursive")

		_objectId = objectId

		So(err, ShouldBeNil)
		So(objectId, ShouldBeGreaterThan, 0)
	})

	Convey("Testing MakeDirectoryRecursive for an existing directory | MakeDirectoryRecursive", t, func() {
		// test the directory '/mtp-test-files/temp_dir/test-MakeDirectoryRecursive'
		objectId, err := MakeDirectoryRecursive(dev, sid, "/mtp-test-files/temp_dir/test-MakeDirectoryRecursive")

		So(err, ShouldBeNil)
		So(objectId, ShouldEqual, _objectId)
	})

	Convey("Testing fullpath='/' | MakeDirectoryRecursive", t, func() {
		// test the directory '/'
		objectId, err := MakeDirectoryRecursive(dev, sid, "/")

		So(err, ShouldBeNil)
		So(objectId, ShouldEqual, ParentObjectId)
	})

	Convey("Testing fullpath='' | MakeDirectoryRecursive", t, func() {
		// test the directory ''
		objectId, err := MakeDirectoryRecursive(dev, sid, "")

		So(err, ShouldBeNil)
		So(objectId, ShouldEqual, ParentObjectId)
	})

	Convey("Creating a new random dir | MakeDirectoryRecursive", t, func() {
		// test the directory '/mtp-test-files/temp_dir/test-MakeDirectoryRecursive/{random}'
		fullpath := fmt.Sprintf("/mtp-test-files/temp_dir/test-MakeDirectoryRecursive/%x", rand.Int31())

		objectId, err := MakeDirectoryRecursive(dev, sid, fullpath)

		So(err, ShouldBeNil)
		So(objectId, ShouldBeGreaterThan, 0)

		exists, isDir, _existingObjectId := FileExists(dev, sid, 0, fullpath)

		So(err, ShouldBeNil)
		So(exists, ShouldEqual, true)
		So(_existingObjectId, ShouldEqual, objectId)
		So(isDir, ShouldEqual, true)
	})

	Convey("filename in the path | 1 | MakeDirectoryRecursive | It should throw an error", t, func() {
		// test the directory '/mtp-test-files/a.txt/folder'
		objectId, err := MakeDirectoryRecursive(dev, sid, "/mtp-test-files/a.txt/folder")

		So(err, ShouldBeError)
		So(objectId, ShouldEqual, 0)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
	})

	Convey("filename in the path | 2 | MakeDirectoryRecursive | It should throw an error", t, func() {
		// test the directory '/mtp-test-files/a.txt/folder'
		objectId, err := MakeDirectoryRecursive(dev, sid, "/mtp-test-files/a.txt")

		So(err, ShouldBeError)
		So(objectId, ShouldEqual, 0)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
	})

	Dispose(dev)
}

func TestWalkDirectory(t *testing.T) {
	dev, err := Initialize(Init{})
	if err != nil {
		log.Panic(err)
	}

	storages, err := FetchStorages(dev)
	if err != nil {
		log.Panic(err)
	}

	sid := storages[0].sid

	Convey("Testing valid directory | with objectId | objectId should be picked up instead of fullPath | WalkDirectory", t, func() {
		// test the root directory [ParentObjectId] | empty [fullPath]
		objectId, totalFiles, err := WalkDirectory(dev, sid, ParentObjectId, "", false, func(objectId uint32, fi *FileInfo) {
			So(objectId, ShouldBeGreaterThan, 0)
			So(fi, ShouldNotBeNil)
			So(fi.ParentId, ShouldEqual, 0)
		})

		So(err, ShouldBeNil)

		So(totalFiles, ShouldBeGreaterThan, 0)
		So(objectId, ShouldEqual, ParentObjectId)

		// test the root directory [ParentObjectId] | [fullPath]='/fake'
		objectId, totalFiles, err = WalkDirectory(dev, sid, ParentObjectId, "/fake", false, func(objectId uint32, fi *FileInfo) {
			So(objectId, ShouldBeGreaterThan, 0)
			So(fi, ShouldNotBeNil)
			So(fi.ParentId, ShouldEqual, 0)
		})

		So(err, ShouldBeNil)
		So(totalFiles, ShouldBeGreaterThan, 0)
		So(objectId, ShouldEqual, ParentObjectId)
	})

	Convey("Testing valid directory | without objectId | fullPath should be picked up instead of objectId | WalkDirectory", t, func() {
		/////////////////
		// test the directory '/mtp-test-files'
		/////////////////
		fullPath := "/mtp-test-files"

		var children []*FileInfo
		objectId1, totalFiles1, err := WalkDirectory(dev, sid, 0, fullPath, false, func(objectId uint32, fi *FileInfo) {
			So(fi.FullPath, ShouldNotEqual, "/mtp-test-files")
			So(objectId, ShouldBeGreaterThan, 0)
			So(fi, ShouldNotBeNil)
			So(fi.ParentId, ShouldBeGreaterThan, 0)
			So(objectId, ShouldEqual, fi.ObjectId)

			children = append(children, fi)
		})

		So(err, ShouldBeNil)
		So(totalFiles1, ShouldBeGreaterThanOrEqualTo, 4)

		// test if [objectId] == [objectId1] of '/mtp-test-files'
		objIdFromPath, _, err := GetObjectFromPath(dev, sid, fullPath)
		So(err, ShouldBeNil)

		So(objectId1, ShouldEqual, objIdFromPath)

		So(len(children), ShouldEqual, totalFiles1)

		/////////////////
		// test the directory '/mtp-test-files/'
		/////////////////
		fullPath = "/mtp-test-files/"
		children = []*FileInfo{}
		objectId2, totalFiles2, err := WalkDirectory(dev, sid, 0, fullPath, false, func(objectId uint32, fi *FileInfo) {
			So(fi.FullPath, ShouldNotEqual, "/mtp-test-files")
			So(objectId, ShouldBeGreaterThan, 0)
			So(fi, ShouldNotBeNil)
			So(fi.ParentId, ShouldBeGreaterThan, 0)
			So(objectId, ShouldEqual, fi.ObjectId)

			children = append(children, fi)
		})

		So(err, ShouldBeNil)
		So(totalFiles2, ShouldBeGreaterThanOrEqualTo, totalFiles1)

		// test if [objectId2] == [objectId1] of [fullPath]
		objIdFromPath, _, err = GetObjectFromPath(dev, sid, fullPath)
		So(err, ShouldBeNil)

		So(objectId1, ShouldEqual, objIdFromPath)
		So(objectId1, ShouldEqual, objectId2)

		So(len(children), ShouldEqual, totalFiles2)

		/////////////////
		// test the directory 'mtp-test-files/'
		/////////////////
		fullPath = "mtp-test-files/"
		children = []*FileInfo{}
		objectId3, totalFiles3, err := WalkDirectory(dev, sid, 0, fullPath, false, func(objectId uint32, fi *FileInfo) {
			So(fi.FullPath, ShouldNotEqual, "/mtp-test-files")
			So(objectId, ShouldBeGreaterThan, 0)
			So(fi, ShouldNotBeNil)
			So(fi.ParentId, ShouldBeGreaterThan, 0)
			So(objectId, ShouldEqual, fi.ObjectId)

			children = append(children, fi)
		})

		So(err, ShouldBeNil)
		So(totalFiles3, ShouldBeGreaterThanOrEqualTo, totalFiles1)

		// test if [objectId3] == [objectId] of [fullPath]
		objIdFromPath, _, err = GetObjectFromPath(dev, sid, fullPath)
		So(err, ShouldBeNil)

		So(objectId3, ShouldEqual, objIdFromPath)

		So(len(children), ShouldEqual, totalFiles3)

		/////////////////
		// test the directory 'mtp-test-files/mock_dir3/'
		/////////////////
		fullPath = "mtp-test-files/mock_dir3/"
		children = []*FileInfo{}

		objectId4, totalFiles4, err := WalkDirectory(dev, sid, 0, fullPath, false, func(objectId uint32, fi *FileInfo) {
			So(fi.FullPath, ShouldNotEqual, "/mtp-test-files/mock_dir3")
			So(objectId, ShouldBeGreaterThan, 0)
			So(fi, ShouldNotBeNil)
			So(fi.ParentId, ShouldBeGreaterThan, 0)
			So(objectId, ShouldEqual, fi.ObjectId)

			children = append(children, fi)
		})

		So(err, ShouldBeNil)
		So(totalFiles4, ShouldBeGreaterThanOrEqualTo, totalFiles1)

		// test if [objectId4] == [objectId] of [fullPath]
		objIdFromPath, _, err = GetObjectFromPath(dev, sid, fullPath)
		So(err, ShouldBeNil)

		So(objectId4, ShouldEqual, objIdFromPath)

		So(len(children), ShouldEqual, 5)
	})

	Convey("Testing valid directory | 1 | recursive=false | WalkDirectory", t, func() {
		//test the directory '/mtp-test-files/mock_dir1/1'
		fullPath := "/mtp-test-files/mock_dir1/1"

		var children []*FileInfo
		objectId, totalFiles, err := WalkDirectory(dev, sid, 0, fullPath, false, func(objectId uint32, fi *FileInfo) {
			So(fi.FullPath, ShouldNotEqual, "/mtp-test-files/mock_dir1/1")
			So(objectId, ShouldBeGreaterThan, 0)
			So(fi, ShouldNotBeNil)
			So(fi.ParentId, ShouldBeGreaterThan, 0)
			So(objectId, ShouldEqual, fi.ObjectId)

			children = append(children, fi)
		})

		So(err, ShouldBeNil)

		So(children, ShouldNotBeNil)
		So(len(children), ShouldEqual, totalFiles)
		So(len(children), ShouldEqual, 1)

		_file0 := children[0]

		So(_file0.ObjectId, ShouldBeGreaterThan, 0)
		So(_file0.Name, ShouldEqual, "a.txt")
		So(_file0.ParentId, ShouldEqual, objectId)
		So(_file0.Info.Filename, ShouldEqual, "a.txt")
		So(_file0.Extension, ShouldEqual, "txt")
		So(_file0.Size, ShouldEqual, 8)
		So(_file0.IsDir, ShouldEqual, false)
		So(_file0.FullPath, ShouldEqual, "/mtp-test-files/mock_dir1/1/a.txt")
		So(_file0.ModTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)
	})

	Convey("Testing valid directory | 2 | recursive=false | WalkDirectory", t, func() {
		//test the directory '/mtp-test-files/mock_dir1/'
		fullPath := "/mtp-test-files/mock_dir1/"

		var children []*FileInfo
		objectId, totalFiles, err := WalkDirectory(dev, sid, 0, fullPath, false, func(objectId uint32, fi *FileInfo) {
			So(fi.FullPath, ShouldNotEqual, "/mtp-test-files/mock_dir1")
			So(objectId, ShouldBeGreaterThan, 0)
			So(fi, ShouldNotBeNil)
			So(fi.ParentId, ShouldBeGreaterThan, 0)
			So(objectId, ShouldEqual, fi.ObjectId)

			children = append(children, fi)
		})

		So(err, ShouldBeNil)

		So(children, ShouldNotBeNil)
		So(len(children), ShouldEqual, totalFiles)
		So(len(children), ShouldEqual, 4)

		_file0 := children[0]

		So(_file0.ObjectId, ShouldBeGreaterThan, 0)
		So(_file0.Name, ShouldEqual, "1")
		So(_file0.ParentId, ShouldEqual, objectId)
		So(_file0.Info.Filename, ShouldEqual, "1")
		So(_file0.Extension, ShouldEqual, "")
		So(_file0.Size, ShouldEqual, 4096)
		So(_file0.IsDir, ShouldEqual, true)
		So(_file0.FullPath, ShouldEqual, "/mtp-test-files/mock_dir1/1")
		So(_file0.ModTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)
	})

	Convey("Testing valid directory | 1 | recursive=true | WalkDirectory", t, func() {
		//test the directory '/mtp-test-files/mock_dir1/'
		fullPath := "/mtp-test-files/mock_dir1/"

		var children []*FileInfo
		objectId, totalFiles, err := WalkDirectory(dev, sid, 0, fullPath, true, func(objectId uint32, fi *FileInfo) {
			So(fi.FullPath, ShouldNotEqual, "/mtp-test-files/mock_dir1")
			So(objectId, ShouldBeGreaterThan, 0)
			So(fi, ShouldNotBeNil)
			So(fi.ParentId, ShouldBeGreaterThan, 0)
			So(objectId, ShouldEqual, fi.ObjectId)

			children = append(children, fi)
		})

		So(err, ShouldBeNil)

		const childrenLength = 9

		So(children, ShouldNotBeNil)
		So(len(children), ShouldEqual, childrenLength)
		So(totalFiles, ShouldEqual, 9)

		_file0 := children[0]

		So(_file0.ObjectId, ShouldBeGreaterThan, 0)
		So(_file0.Name, ShouldEqual, "1")
		So(_file0.ParentId, ShouldEqual, objectId)
		So(_file0.Info.Filename, ShouldEqual, "1")
		So(_file0.Extension, ShouldEqual, "")
		So(_file0.Size, ShouldEqual, 4096)
		So(_file0.IsDir, ShouldEqual, true)
		So(_file0.FullPath, ShouldEqual, "/mtp-test-files/mock_dir1/1")
		So(_file0.ModTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)

		// test level 1 objects
		dirList1 := [childrenLength]string{"/mtp-test-files/mock_dir1/1", "/mtp-test-files/mock_dir1/1/a.txt", "/mtp-test-files/mock_dir1/a.txt", "/mtp-test-files/mock_dir1/3", "/mtp-test-files/mock_dir1/3/b.txt", "/mtp-test-files/mock_dir1/3/2", "/mtp-test-files/mock_dir1/3/2/b.txt", "/mtp-test-files/mock_dir1/2", "/mtp-test-files/mock_dir1/2/b.txt"}

		for i, _dir := range dirList1 {
			So(
				children[i].FullPath, ShouldEqual, _dir,
			)
		}
	})

	Convey("Testing non exisiting file | WalkDirectory | It should throw an error", t, func() {
		// test the directory '/fake' | recursive=true
		var children []*FileInfo
		objectId, totalFiles, err := WalkDirectory(dev, sid, 0, "/fake", true, func(objectId uint32, fi *FileInfo) {
			children = append(children, fi)
		})

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objectId, ShouldEqual, 0)
		So(totalFiles, ShouldEqual, 0)
		So(len(children), ShouldEqual, 0)

		// test the directory '/fake' | recursive=false
		children = []*FileInfo{}
		objectId, totalFiles, err = WalkDirectory(dev, sid, 0, "/fake", false, func(objectId uint32, fi *FileInfo) {
			children = append(children, fi)
		})

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objectId, ShouldEqual, 0)
		So(totalFiles, ShouldEqual, 0)
		So(len(children), ShouldEqual, 0)

		// test the directory '/mtp-test-files/fake' | recursive=true
		children = []*FileInfo{}
		objectId, totalFiles, err = WalkDirectory(dev, sid, 0, "/mtp-test-files/fake", true, func(objectId uint32, fi *FileInfo) {
			children = append(children, fi)
		})

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objectId, ShouldEqual, 0)
		So(totalFiles, ShouldEqual, 0)
		So(len(children), ShouldEqual, 0)

		// test the directory '/mtp-test-files/fake' | recursive=false
		children = []*FileInfo{}
		objectId, totalFiles, err = WalkDirectory(dev, sid, 0, "/mtp-test-files/fake", false, func(objectId uint32, fi *FileInfo) {
			children = append(children, fi)
		})

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objectId, ShouldEqual, 0)
		So(totalFiles, ShouldEqual, 0)
		So(len(children), ShouldEqual, 0)

		// test the directory '/mtp-test-files/a.txt'
		children = []*FileInfo{}
		objectId, totalFiles, err = WalkDirectory(dev, sid, 0, "/mtp-test-files/a.txt", true, func(objectId uint32, fi *FileInfo) {
			children = append(children, fi)
		})

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objectId, ShouldEqual, 0)
		So(totalFiles, ShouldEqual, 0)
		So(len(children), ShouldEqual, 0)

		// test the directory=''
		objectId, totalFiles, err = WalkDirectory(dev, sid, 0, "", true, func(objectId uint32, fi *FileInfo) {
			children = append(children, fi)
		})

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objectId, ShouldEqual, 0)
		So(totalFiles, ShouldEqual, 0)
		So(len(children), ShouldEqual, 0)
	})

	Dispose(dev)
}

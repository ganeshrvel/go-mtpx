package main

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"math/rand"
	"testing"
)

func TestListDirectory(t *testing.T) {
	dev, err := Initialize(Init{})
	if err != nil {
		log.Panic(err)
	}

	storages, err := FetchStorages(dev)
	if err != nil {
		log.Panic(err)
	}

	sid := storages[0].sid

	Convey("Testing valid file | with objectId | objectId should be picked up instead of fullPath | ListDirectory", t, func() {
		// test the root directory [ParentObjectId] | empty [fullPath]
		files, err := ListDirectory(dev, sid, ParentObjectId, "")

		So(err, ShouldBeNil)
		So(files, ShouldNotBeNil)
		So(len(*files), ShouldBeGreaterThan, 0)

		// test the root directory [ParentObjectId] | fake [fullPath]
		files, err = ListDirectory(dev, sid, ParentObjectId, "/fake")

		So(err, ShouldBeNil)
		So(files, ShouldNotBeNil)
		So(len(*files), ShouldBeGreaterThan, 0)
	})

	Convey("Testing valid file | without objectId | fullPath should be picked up instead of objectId | ListDirectory", t, func() {
		// test the directory '/mtp-test-files'
		files, err := ListDirectory(dev, sid, 0, "/mtp-test-files")

		So(err, ShouldBeNil)
		So(files, ShouldNotBeNil)
		So(len(*files), ShouldBeGreaterThan, 0)

		// test the directory '/mtp-test-files/'
		files, err = ListDirectory(dev, sid, 0, "/mtp-test-files/")

		So(err, ShouldBeNil)
		So(files, ShouldNotBeNil)
		So(len(*files), ShouldBeGreaterThan, 0)

		// test the directory 'mtp-test-files/'
		files, err = ListDirectory(dev, sid, 0, "mtp-test-files/")

		So(err, ShouldBeNil)
		So(files, ShouldNotBeNil)
		So(len(*files), ShouldBeGreaterThan, 0)

		// test the directory 'mtp-test-files'
		files, err = ListDirectory(dev, sid, 0, "mtp-test-files")

		So(err, ShouldBeNil)
		So(files, ShouldNotBeNil)
		So(len(*files), ShouldBeGreaterThan, 0)

		// test the directory 'mtp-test-files/mock_dir3/'
		files, err = ListDirectory(dev, sid, 0, "mtp-test-files/mock_dir3/")

		So(err, ShouldBeNil)
		So(files, ShouldNotBeNil)
		So(len(*files), ShouldEqual, 5)
	})

	Convey("Testing valid file | ListDirectory", t, func() {
		// test the directory '/mtp-test-files/mock_dir1/1'
		files, err := ListDirectory(dev, sid, 0, "/mtp-test-files/mock_dir1/1")

		So(err, ShouldBeNil)

		_files := *files

		So(_files, ShouldNotBeNil)
		So(len(_files), ShouldEqual, 1)
		So(_files[0].ObjectId, ShouldBeGreaterThan, 0)
		So(_files[0].Name, ShouldEqual, "a.txt")
		So(_files[0].ParentId, ShouldBeGreaterThan, 0)
		So(_files[0].Info.Filename, ShouldEqual, "a.txt")
		So(_files[0].Extension, ShouldEqual, "txt")
		So(_files[0].Size, ShouldEqual, 8)
		So(_files[0].IsDir, ShouldEqual, false)
		So(_files[0].FullPath, ShouldEqual, "/mtp-test-files/mock_dir1/1/a.txt")
		So(_files[0].ModTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)
	})

	Convey("Testing non exisiting file | ListDirectory | It should throw an error", t, func() {
		// test the directory '/fake'
		files, err := ListDirectory(dev, sid, 0, "/fake")

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(files, ShouldBeNil)

		// test the directory '/mtp-test-files'
		files, err = ListDirectory(dev, sid, 0, "/mtp-test-files/a.txt")

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(files, ShouldBeNil)

		// test the directory '/mtp-test-files/fake'
		files, err = ListDirectory(dev, sid, 0, "/mtp-test-files/fake")

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(files, ShouldBeNil)
	})

	Dispose(dev)
}

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
	Convey("Creating a new dir | MakeDirectory", t, func() {
		// test the directory '/mtp-test-files/temp_dir/test-MakeDirectory'
		objectId, err := MakeDirectory(dev, sid, "/mtp-test-files/temp_dir", "test-MakeDirectory")

		_objectId = objectId

		So(err, ShouldBeNil)
		So(objectId, ShouldBeGreaterThan, 0)
	})

	Convey("Testing MakeDirectory for an existing directory | MakeDirectory", t, func() {
		// test the directory '/mtp-test-files/temp_dir/test-MakeDirectory'
		objectId, err := MakeDirectory(dev, sid, "/mtp-test-files/temp_dir", "test-MakeDirectory")

		So(err, ShouldBeNil)
		So(objectId, ShouldEqual, _objectId)
	})

	Convey("Creating a new random dir | MakeDirectory", t, func() {
		// test the directory '/mtp-test-files/temp_dir/test-MakeDirectory/{random}'
		filename := fmt.Sprintf("%x", rand.Int31())

		objectId, err := MakeDirectory(dev, sid, "/mtp-test-files/temp_dir/test-MakeDirectory", filename)

		So(err, ShouldBeNil)
		So(objectId, ShouldBeGreaterThan, 0)

		exists, isDir, _existingObjectId := FileExists(dev, sid, getFullPath("/mtp-test-files/temp_dir/test-MakeDirectory", filename))

		So(err, ShouldBeNil)
		So(exists, ShouldEqual, true)
		So(_existingObjectId, ShouldEqual, objectId)
		So(isDir, ShouldEqual, true)
	})

	Convey("invalid path | MakeDirectory | It should throw an error", t, func() {
		// test the directory '/fake/test'
		objectId, err := MakeDirectory(dev, sid, "fake", "test")

		So(err, ShouldBeError)
		So(objectId, ShouldEqual, 0)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
	})

	Convey("empty folder name | MakeDirectory | It should throw an error", t, func() {
		// test the directory '/'
		objectId, err := MakeDirectory(dev, sid, "fake", "")

		So(err, ShouldBeError)
		So(objectId, ShouldEqual, 0)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
	})

	Convey("filename in the path | MakeDirectory | It should throw an error", t, func() {
		// test the directory '/mtp-test-files/a.txt'
		objectId, err := MakeDirectory(dev, sid, "/mtp-test-files/", "a.txt")

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

		exists, isDir, _existingObjectId := FileExists(dev, sid, fullpath)

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

func TestFetchDirectoryTree(t *testing.T) {
	dev, err := Initialize(Init{})
	if err != nil {
		log.Panic(err)
	}

	storages, err := FetchStorages(dev)
	if err != nil {
		log.Panic(err)
	}

	sid := storages[0].sid

	Convey("Testing valid directory | with objectId | objectId should be picked up instead of fullPath | FetchDirectoryTree", t, func() {
		// test the root directory [ParentObjectId] | empty [fullPath]
		dirListing := &DirectoryTree{}
		objectId, totalFiles, err := FetchDirectoryTree(dev, sid, ParentObjectId, "", false, dirListing)

		So(err, ShouldBeNil)
		So(dirListing, ShouldNotBeNil)
		So(totalFiles, ShouldBeGreaterThan, 0)
		So(objectId, ShouldEqual, ParentObjectId)
		pDirListing := *dirListing

		So(len(pDirListing.Children), ShouldBeGreaterThan, 0)

		// test the root directory [ParentObjectId] | [fullPath]='/fake'
		dirListing = &DirectoryTree{}
		objectId, totalFiles, err = FetchDirectoryTree(dev, sid, ParentObjectId, "/fake", false, dirListing)

		So(err, ShouldBeNil)
		So(dirListing, ShouldNotBeNil)
		So(totalFiles, ShouldBeGreaterThan, 0)
		So(objectId, ShouldEqual, ParentObjectId)
		pDirListing = *dirListing

		So(len(pDirListing.Children), ShouldBeGreaterThan, 0)
	})

	Convey("Testing valid directory | without objectId | fullPath should be picked up instead of objectId | FetchDirectoryTree", t, func() {
		/////////////////
		// test the directory '/mtp-test-files'
		/////////////////
		dirListing := &DirectoryTree{}
		fullPath := "/mtp-test-files"

		objectId1, totalFiles1, err := FetchDirectoryTree(dev, sid, 0, fullPath, false, dirListing)

		So(err, ShouldBeNil)
		So(dirListing, ShouldNotBeNil)
		So(totalFiles1, ShouldBeGreaterThanOrEqualTo, 4)
		pDirListing := *dirListing

		// test if [objectId] == [objectId1] of '/mtp-test-files'
		objIdFromPath, _, err := GetObjectFromPath(dev, sid, fullPath)
		So(err, ShouldBeNil)

		So(objectId1, ShouldEqual, objIdFromPath)

		So(len(pDirListing.Children), ShouldEqual, totalFiles1)

		/////////////////
		// test the directory '/mtp-test-files/'
		/////////////////
		dirListing = &DirectoryTree{}
		fullPath = "/mtp-test-files/"

		objectId2, totalFiles2, err := FetchDirectoryTree(dev, sid, 0, fullPath, false, dirListing)

		So(err, ShouldBeNil)
		So(dirListing, ShouldNotBeNil)
		So(totalFiles2, ShouldBeGreaterThanOrEqualTo, totalFiles1)
		pDirListing = *dirListing

		// test if [objectId2] == [objectId1] of [fullPath]
		objIdFromPath, _, err = GetObjectFromPath(dev, sid, fullPath)
		So(err, ShouldBeNil)

		So(objectId1, ShouldEqual, objIdFromPath)
		So(objectId1, ShouldEqual, objectId2)

		So(len(pDirListing.Children), ShouldEqual, totalFiles2)

		/////////////////
		// test the directory 'mtp-test-files/'
		/////////////////
		dirListing = &DirectoryTree{}
		fullPath = "mtp-test-files/"

		objectId3, totalFiles3, err := FetchDirectoryTree(dev, sid, 0, fullPath, false, dirListing)

		So(err, ShouldBeNil)
		So(dirListing, ShouldNotBeNil)
		So(totalFiles3, ShouldBeGreaterThanOrEqualTo, totalFiles1)
		pDirListing = *dirListing

		// test if [objectId3] == [objectId] of [fullPath]
		objIdFromPath, _, err = GetObjectFromPath(dev, sid, fullPath)
		So(err, ShouldBeNil)

		So(objectId3, ShouldEqual, objIdFromPath)

		So(len(pDirListing.Children), ShouldEqual, totalFiles3)

		/////////////////
		// test the directory 'mtp-test-files/mock_dir3/'
		/////////////////
		dirListing = &DirectoryTree{}
		fullPath = "mtp-test-files/mock_dir3/"

		objectId4, totalFiles4, err := FetchDirectoryTree(dev, sid, 0, fullPath, false, dirListing)

		So(err, ShouldBeNil)
		So(dirListing, ShouldNotBeNil)
		So(totalFiles4, ShouldBeGreaterThanOrEqualTo, totalFiles1)
		pDirListing = *dirListing

		// test if [objectId4] == [objectId] of [fullPath]
		objIdFromPath, _, err = GetObjectFromPath(dev, sid, fullPath)
		So(err, ShouldBeNil)

		So(objectId4, ShouldEqual, objIdFromPath)

		So(len(pDirListing.Children), ShouldEqual, 5)
	})

	Convey("Testing valid directory | 1 | recursive=false | FetchDirectoryTree", t, func() {
		//test the directory '/mtp-test-files/mock_dir1/1'
		fullPath := "/mtp-test-files/mock_dir1/1"

		dirListing := &DirectoryTree{}

		objectId, totalFiles, err := FetchDirectoryTree(dev, sid, 0, fullPath, false, dirListing)

		So(err, ShouldBeNil)

		pDirListing := *dirListing

		So(pDirListing.Children, ShouldNotBeNil)
		So(len(pDirListing.Children), ShouldEqual, totalFiles)
		So(len(pDirListing.Children), ShouldEqual, 1)
		So(pDirListing.ObjectId, ShouldEqual, objectId)

		_file0 := pDirListing.Children[0]

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

	Convey("Testing valid directory | 2 | recursive=false | FetchDirectoryTree", t, func() {
		//test the directory '/mtp-test-files/mock_dir1/'
		fullPath := "/mtp-test-files/mock_dir1/"

		dirListing := &DirectoryTree{}

		objectId, totalFiles, err := FetchDirectoryTree(dev, sid, 0, fullPath, false, dirListing)

		So(err, ShouldBeNil)

		pDirListing := *dirListing

		So(pDirListing.Children, ShouldNotBeNil)
		So(len(pDirListing.Children), ShouldEqual, totalFiles)
		So(len(pDirListing.Children), ShouldEqual, 4)
		So(pDirListing.ObjectId, ShouldEqual, objectId)

		_file0 := pDirListing.Children[0]

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

	Convey("Testing valid directory | recursive=true | FetchDirectoryTree", t, func() {
		//test the directory '/mtp-test-files/mock_dir1/'
		fullPath := "/mtp-test-files/mock_dir1/"

		dirListing := &DirectoryTree{}
		objectId, totalFiles, err := FetchDirectoryTree(dev, sid, 0, fullPath, true, dirListing)

		So(err, ShouldBeNil)

		pDirListing := *dirListing
		const l1ChildrenLength = 4

		So(pDirListing.Children, ShouldNotBeNil)
		So(len(pDirListing.Children), ShouldEqual, l1ChildrenLength)
		So(totalFiles, ShouldEqual, 9)
		So(pDirListing.ObjectId, ShouldEqual, objectId)

		_file0 := pDirListing.Children[0]

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
		dirList1 := [l1ChildrenLength]string{"/mtp-test-files/mock_dir1/1", "/mtp-test-files/mock_dir1/a.txt", "/mtp-test-files/mock_dir1/3", "/mtp-test-files/mock_dir1/2"}

		for i, _dir := range dirList1 {
			So(
				pDirListing.Children[i].FullPath, ShouldEqual, _dir,
			)
		}

		// test level 2 objects dir='1'
		const childrenLengthDir1 = 1
		dirList2 := [childrenLengthDir1]string{"/mtp-test-files/mock_dir1/1/a.txt"}

		for i, _dir := range dirList2 {
			So(
				pDirListing.Children[0].Children[i].FullPath, ShouldEqual, _dir,
			)
		}

		// test level 2 objects dir='2'
		const childrenLengthDir2 = 1
		dirList3 := [childrenLengthDir2]string{"/mtp-test-files/mock_dir1/2/b.txt"}

		for i, _dir := range dirList3 {
			So(
				pDirListing.Children[3].Children[i].FullPath, ShouldEqual, _dir,
			)
		}

		// test level 2 objects dir='3'
		const childrenLengthDir3 = 2
		dirList4 := [childrenLengthDir3]string{"/mtp-test-files/mock_dir1/3/b.txt", "/mtp-test-files/mock_dir1/3/2"}

		for i, _dir := range dirList4 {
			So(
				pDirListing.Children[2].Children[i].FullPath, ShouldEqual, _dir,
			)
		}

		// test level 3 objects dir='3/2'
		const childrenLengthDir4 = 1
		dirList5 := [childrenLengthDir4]string{"/mtp-test-files/mock_dir1/3/2/b.txt"}

		for i, _dir := range dirList5 {
			So(
				pDirListing.Children[2].Children[1].Children[i].FullPath, ShouldEqual, _dir,
			)
		}
	})

	Convey("Testing non exisiting file | FetchDirectoryTree | It should throw an error", t, func() {
		// test the directory '/fake' | recursive=true
		dirListing := &DirectoryTree{}
		objectId, totalFiles, err := FetchDirectoryTree(dev, sid, 0, "/fake", true, dirListing)

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objectId, ShouldEqual, 0)
		So(totalFiles, ShouldEqual, 0)

		// test the directory '/fake' | recursive=false
		dirListing = &DirectoryTree{}
		objectId, totalFiles, err = FetchDirectoryTree(dev, sid, 0, "/fake", false, dirListing)

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objectId, ShouldEqual, 0)
		So(totalFiles, ShouldEqual, 0)

		// test the directory '/mtp-test-files/fake' | recursive=true
		dirListing = &DirectoryTree{}
		objectId, totalFiles, err = FetchDirectoryTree(dev, sid, 0, "/mtp-test-files/fake", true, dirListing)

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objectId, ShouldEqual, 0)
		So(totalFiles, ShouldEqual, 0)

		// test the directory '/mtp-test-files/fake' | recursive=false
		dirListing = &DirectoryTree{}
		objectId, totalFiles, err = FetchDirectoryTree(dev, sid, 0, "/mtp-test-files/fake", false, dirListing)

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objectId, ShouldEqual, 0)
		So(totalFiles, ShouldEqual, 0)

		// test the directory '/mtp-test-files/a.txt'
		dirListing = &DirectoryTree{}
		objectId, totalFiles, err = FetchDirectoryTree(dev, sid, 0, "/mtp-test-files/a.txt", true, dirListing)

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objectId, ShouldEqual, 0)
		So(totalFiles, ShouldEqual, 0)

		// test the directory=''
		dirListing = &DirectoryTree{}
		objectId, totalFiles, err = FetchDirectoryTree(dev, sid, 0, "", true, dirListing)

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objectId, ShouldEqual, 0)
		So(totalFiles, ShouldEqual, 0)
	})

	Dispose(dev)
}

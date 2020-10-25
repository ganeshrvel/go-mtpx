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
		fi, err := GetObjectFromPath(dev, sid, "/mtp-test-files/temp_dir")
		So(err, ShouldBeNil)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, true)

		objectId, err := MakeDirectory(dev, sid, fi.ObjectId, "", "test-MakeDirectoryUsingParentId")

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
		fi, err := GetObjectFromPath(dev, sid, "/mtp-test-files/temp_dir")
		So(err, ShouldBeNil)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, true)

		objectId, err := MakeDirectory(dev, sid, fi.ObjectId, "", "test-MakeDirectoryUsingParentId")

		So(err, ShouldBeNil)
		So(objectId, ShouldEqual, _objectId2)
	})

	Convey("Creating a new random dir | fullpath | MakeDirectory", t, func() {
		// test the directory '/mtp-test-files/temp_dir/test-MakeDirectory/{random}'
		filename := fmt.Sprintf("%x", rand.Int31())

		objectId, err := MakeDirectory(dev, sid, 0, "/mtp-test-files/temp_dir/test-MakeDirectory", filename)

		So(err, ShouldBeNil)
		So(objectId, ShouldBeGreaterThan, 0)

		exists, fi := FileExists(dev, sid, 0, getFullPath("/mtp-test-files/temp_dir/test-MakeDirectory", filename))

		So(err, ShouldBeNil)
		So(exists, ShouldEqual, true)
		So(fi.ObjectId, ShouldEqual, objectId)
		So(fi.IsDir, ShouldEqual, true)
	})

	Convey("Creating a new random dir | parentId | MakeDirectory", t, func() {
		// test the directory '/mtp-test-files/temp_dir/test-MakeDirectory/{random}'
		filename := fmt.Sprintf("%x", rand.Int31())
		fi, err := GetObjectFromPath(dev, sid, "/mtp-test-files/temp_dir/test-MakeDirectoryUsingParentId")
		So(err, ShouldBeNil)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, true)

		objectId, err := MakeDirectory(dev, sid, fi.ObjectId, "", filename)

		So(err, ShouldBeNil)
		So(objectId, ShouldBeGreaterThan, 0)

		exists, fi := FileExists(dev, sid, 0, getFullPath("/mtp-test-files/temp_dir/test-MakeDirectoryUsingParentId", filename))

		So(err, ShouldBeNil)
		So(exists, ShouldEqual, true)
		So(fi.ObjectId, ShouldEqual, objectId)
		So(fi.IsDir, ShouldEqual, true)
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

		exists, fi := FileExists(dev, sid, 0, fullpath)

		So(err, ShouldBeNil)
		So(exists, ShouldEqual, true)
		So(fi.ObjectId, ShouldEqual, objectId)
		So(fi.IsDir, ShouldEqual, true)
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

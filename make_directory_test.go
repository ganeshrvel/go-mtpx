package main

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"math/rand"
	"testing"
)

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

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

	Convey("Testing valid file | with objectId | ListDirectory", t, func() {
		// test the directory '/mtp-test-files'
		files, err := ListDirectory(dev, sid, ParentObjectId, "/mtp-test-files")

		So(err, ShouldBeNil)
		So(files, ShouldNotBeNil)
		So(len(*files), ShouldBeGreaterThan, 0)

		// test the directory '/fake'
		files, err = ListDirectory(dev, sid, ParentObjectId, "/fake")

		So(err, ShouldBeNil)
		So(files, ShouldNotBeNil)
		So(len(*files), ShouldBeGreaterThan, 0)

	})

	Convey("Testing valid file | without objectId | ListDirectory", t, func() {
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

	var objectId uint32
	Convey("Creating a new dir | MakeDirectory", t, func() {
		// test the directory '/mtp-test-files/test-MakeDirectory'
		_objectId, err := MakeDirectory(dev, sid, "/mtp-test-files", "test-MakeDirectory")

		objectId = _objectId

		So(err, ShouldBeNil)
		So(objectId, ShouldBeGreaterThan, 0)
	})

	Convey("Testing MakeDirectory for an existing directory | MakeDirectory", t, func() {
		// test the directory '/mtp-test-files/test-MakeDirectory'
		objectId, err := MakeDirectory(dev, sid, "/mtp-test-files", "test-MakeDirectory")

		So(err, ShouldBeNil)
		So(objectId, ShouldEqual, objectId)
	})

	Convey("Creating a new random dir | MakeDirectory", t, func() {
		// test the directory '/mtp-test-files/test-MakeDirectory/{random}'
		filename := fmt.Sprintf("%x", rand.Int31())

		objectId, err := MakeDirectory(dev, sid, "/mtp-test-files/test-MakeDirectory", filename)

		So(err, ShouldBeNil)
		So(objectId, ShouldBeGreaterThan, 0)

		exists, _existingObjectId := FileExists(dev, sid, getFullPath("/mtp-test-files/test-MakeDirectory", filename))

		So(err, ShouldBeNil)
		So(exists, ShouldEqual, true)
		So(_existingObjectId, ShouldEqual, objectId)
	})

	Convey("invalid path | MakeDirectory | It should throw an error", t, func() {
		// test the directory '/fake/test'
		objectId, err := MakeDirectory(dev, sid, "fake", "test")

		So(err, ShouldBeError)
		So(objectId, ShouldEqual, 0)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
	})

	Dispose(dev)
}
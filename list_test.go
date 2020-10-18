package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"testing"
)

func TestGetObjectIdFromPath(t *testing.T) {
	dev, err := Initialize(Init{})
	if err != nil {
		log.Panic(err)
	}

	storages, err := FetchStorages(dev)
	if err != nil {
		log.Panic(err)
	}

	sid := storages[0].sid

	Convey("Testing valid file | GetObjectIdFromPath", t, func() {
		// test the directory '/mocks'
		objId, err := GetObjectIdFromPath(dev, sid, "/mocks")

		So(err, ShouldBeNil)
		So(objId, ShouldBeGreaterThan, 0)

		// test the directory '/mocks/'
		objId, err = GetObjectIdFromPath(dev, sid, "/mocks/")

		So(err, ShouldBeNil)
		So(objId, ShouldEqual, objId)

		// test the directory 'mocks/'
		objId, err = GetObjectIdFromPath(dev, sid, "mocks/")

		So(err, ShouldBeNil)
		So(objId, ShouldEqual, objId)

		// test the directory 'mocks'
		objId, err = GetObjectIdFromPath(dev, sid, "mocks")

		So(err, ShouldBeNil)
		So(objId, ShouldEqual, objId)

		// test the file '/mocks/a.txt'
		objId, err = GetObjectIdFromPath(dev, sid, "/mocks/a.txt")

		So(err, ShouldBeNil)
		So(objId, ShouldBeGreaterThan, 0)

		// test the file '/mocks/mock_dir1/a.txt'
		objId, err = GetObjectIdFromPath(dev, sid, "/mocks/mock_dir1/a.txt")

		So(err, ShouldBeNil)
		So(objId, ShouldBeGreaterThan, 0)

		// test the file 'mocks/a.txt'
		objId, err = GetObjectIdFromPath(dev, sid, "mocks/a.txt")

		So(err, ShouldBeNil)
		So(objId, ShouldEqual, objId)

		// test the file 'mocks/a.txt/'
		objId, err = GetObjectIdFromPath(dev, sid, "mocks/a.txt/")

		So(err, ShouldBeNil)
		So(objId, ShouldEqual, objId)
	})

	Convey("Testing non exisiting file | GetObjectIdFromPath | It should throw an error", t, func() {
		// test the file 'fake_file'
		objId, err := GetObjectIdFromPath(dev, sid, "fake_file")

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objId, ShouldEqual, 0)

		// test the file 'mocks/b'
		objId, err = GetObjectIdFromPath(dev, sid, "mocks/b")

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objId, ShouldEqual, 0)

		// test the file 'mocks/b'
		objId, err = GetObjectIdFromPath(dev, sid, "mocks/a.txt/1")

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objId, ShouldEqual, 0)
	})

	Dispose(dev)
}

func TestGetObjectIdFromFilename(t *testing.T) {
	dev, err := Initialize(Init{})
	if err != nil {
		log.Panic(err)
	}

	storages, err := FetchStorages(dev)
	if err != nil {
		log.Panic(err)
	}

	sid := storages[0].sid

	Convey("Testing valid file | GetObjectIdFromFilename", t, func() {
		// test the directory '/mocks'
		objId, isDir, err := GetObjectIdFromFilename(dev, sid, ParentObjectId, "mocks")

		So(err, ShouldBeNil)
		So(objId, ShouldBeGreaterThan, 0)
		So(isDir, ShouldEqual, true)

		// test the file '/mocks/a.txt'
		objId, isDir, err = GetObjectIdFromFilename(dev, sid, objId, "a.txt")

		So(err, ShouldBeNil)
		So(objId, ShouldEqual, objId)
		So(isDir, ShouldEqual, false)
	})

	Convey("Testing non exisiting file | GetObjectIdFromFilename | It should throw an error", t, func() {
		// test the file 'fake_file'
		objId, isDir, err := GetObjectIdFromFilename(dev, sid, ParentObjectId, "fake_file")

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, FileNotFoundError{})
		So(objId, ShouldEqual, 0)
		So(isDir, ShouldEqual, false)

		// test the file '/mocks'
		objId, isDir, err = GetObjectIdFromFilename(dev, sid, ParentObjectId, "/mocks")

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, FileNotFoundError{})
		So(objId, ShouldEqual, 0)
		So(isDir, ShouldEqual, false)

		// test the file 'mocks/'
		objId, isDir, err = GetObjectIdFromFilename(dev, sid, ParentObjectId, "mocks/")

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, FileNotFoundError{})
		So(objId, ShouldEqual, 0)
		So(isDir, ShouldEqual, false)
	})

	Dispose(dev)
}

func TestFileExists(t *testing.T) {
	dev, err := Initialize(Init{})
	if err != nil {
		log.Panic(err)
	}

	storages, err := FetchStorages(dev)
	if err != nil {
		log.Panic(err)
	}

	sid := storages[0].sid

	Convey("Testing valid file | FileExists", t, func() {
		// test the directory '/mocks'
		exists := FileExists(dev, sid, "/mocks/")
		So(exists, ShouldEqual, true)

		// test the file '/mocks/a.txt'
		exists = FileExists(dev, sid, "/mocks/a.txt")
		So(exists, ShouldEqual, true)

		// test the directory 'mocks/'
		exists = FileExists(dev, sid, "mocks/")
		So(exists, ShouldEqual, true)

		// test the directory 'mocks'
		exists = FileExists(dev, sid, "mocks")
		So(exists, ShouldEqual, true)

		// test the file '/mocks/a.txt/'
		exists = FileExists(dev, sid, "/mocks/a.txt/")
		So(exists, ShouldEqual, true)
	})

	Convey("Testing non existing file | FileExists | Should throw error", t, func() {
		// test the directory '/fake'
		exists := FileExists(dev, sid, "/fake/")
		So(exists, ShouldEqual, false)

		// test the file '/mocks/fake.txt'
		exists = FileExists(dev, sid, "/mocks/fake.txt")
		So(exists, ShouldEqual, false)
	})

	Dispose(dev)
}

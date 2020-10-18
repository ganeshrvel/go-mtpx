package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"testing"
)

func TestGetObjectUsingPath(t *testing.T) {
	dev, err := Initialize(Init{})
	if err != nil {
		log.Panic(err)
	}

	storages, err := FetchStorages(dev)
	if err != nil {
		log.Panic(err)
	}

	sid := storages[0].sid

	Convey("Testing valid file | GetObjectUsingPath", t, func() {
		// test the directory '/mocks'
		objId, isDir, err := GetObjectUsingPath(dev, sid, "/mocks")

		So(err, ShouldBeNil)
		So(objId, ShouldBeGreaterThan, 0)
		So(isDir, ShouldEqual, true)

		// test the directory '/mocks/'
		objId, isDir, err = GetObjectUsingPath(dev, sid, "/mocks/")

		So(err, ShouldBeNil)
		So(objId, ShouldEqual, objId)
		So(isDir, ShouldEqual, true)

		// test the directory 'mocks/'
		objId, isDir, err = GetObjectUsingPath(dev, sid, "mocks/")

		So(err, ShouldBeNil)
		So(objId, ShouldEqual, objId)
		So(isDir, ShouldEqual, true)

		// test the directory 'mocks'
		objId, isDir, err = GetObjectUsingPath(dev, sid, "mocks")

		So(err, ShouldBeNil)
		So(objId, ShouldEqual, objId)
		So(isDir, ShouldEqual, true)

		// test the file '/mocks/a.txt'
		objId, isDir, err = GetObjectUsingPath(dev, sid, "/mocks/a.txt")

		So(err, ShouldBeNil)
		So(objId, ShouldBeGreaterThan, 0)
		So(isDir, ShouldEqual, false)

		// test the file '/mocks/mock_dir1/a.txt'
		objId, isDir, err = GetObjectUsingPath(dev, sid, "/mocks/mock_dir1/a.txt")

		So(err, ShouldBeNil)
		So(objId, ShouldBeGreaterThan, 0)
		So(isDir, ShouldEqual, false)

		// test the file 'mocks/a.txt'
		objId, isDir, err = GetObjectUsingPath(dev, sid, "mocks/a.txt")

		So(err, ShouldBeNil)
		So(objId, ShouldEqual, objId)
		So(isDir, ShouldEqual, false)

		// test the file 'mocks/a.txt/'
		objId, isDir, err = GetObjectUsingPath(dev, sid, "mocks/a.txt/")

		So(err, ShouldBeNil)
		So(objId, ShouldEqual, objId)
		So(isDir, ShouldEqual, false)
	})

	Convey("Testing non exisiting file | GetObjectUsingPath | It should throw an error", t, func() {
		// test the file 'fake_file'
		objId, isDir, err := GetObjectUsingPath(dev, sid, "fake_file")

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objId, ShouldEqual, 0)
		So(isDir, ShouldEqual, false)

		// test the file 'mocks/b'
		objId, isDir, err = GetObjectUsingPath(dev, sid, "mocks/b")

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objId, ShouldEqual, 0)
		So(isDir, ShouldEqual, false)

		// test the file 'mocks/b'
		objId, isDir, err = GetObjectUsingPath(dev, sid, "mocks/a.txt/1")

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objId, ShouldEqual, 0)
		So(isDir, ShouldEqual, false)
	})

	Dispose(dev)
}

func TestGetObjectUsingParentId(t *testing.T) {
	dev, err := Initialize(Init{})
	if err != nil {
		log.Panic(err)
	}

	storages, err := FetchStorages(dev)
	if err != nil {
		log.Panic(err)
	}

	sid := storages[0].sid

	Convey("Testing valid file | GetObjectUsingParentId", t, func() {
		// test the directory '/mocks'
		objId, isDir, err := GetObjectUsingParentId(dev, sid, ParentObjectId, "mocks")

		So(err, ShouldBeNil)
		So(objId, ShouldBeGreaterThan, 0)
		So(isDir, ShouldEqual, true)

		// test the file '/mocks/a.txt'
		objId, isDir, err = GetObjectUsingParentId(dev, sid, objId, "a.txt")

		So(err, ShouldBeNil)
		So(objId, ShouldEqual, objId)
		So(isDir, ShouldEqual, false)
	})

	Convey("Testing non exisiting file | GetObjectUsingParentId | It should throw an error", t, func() {
		// test the file 'fake_file'
		objId, isDir, err := GetObjectUsingParentId(dev, sid, ParentObjectId, "fake_file")

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, FileNotFoundError{})
		So(objId, ShouldEqual, 0)
		So(isDir, ShouldEqual, false)

		// test the file '/mocks'
		objId, isDir, err = GetObjectUsingParentId(dev, sid, ParentObjectId, "/mocks")

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, FileNotFoundError{})
		So(objId, ShouldEqual, 0)
		So(isDir, ShouldEqual, false)

		// test the file 'mocks/'
		objId, isDir, err = GetObjectUsingParentId(dev, sid, ParentObjectId, "mocks/")

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

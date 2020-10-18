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
		// test the directory '/mtp-test-files'
		objId, isDir, err := GetObjectUsingPath(dev, sid, "/mtp-test-files")

		So(err, ShouldBeNil)
		So(objId, ShouldBeGreaterThan, 0)
		So(isDir, ShouldEqual, true)

		// test the directory '/mtp-test-files/'
		objId, isDir, err = GetObjectUsingPath(dev, sid, "/mtp-test-files/")

		So(err, ShouldBeNil)
		So(objId, ShouldEqual, objId)
		So(isDir, ShouldEqual, true)

		// test the directory 'mtp-test-files/'
		objId, isDir, err = GetObjectUsingPath(dev, sid, "mtp-test-files/")

		So(err, ShouldBeNil)
		So(objId, ShouldEqual, objId)
		So(isDir, ShouldEqual, true)

		// test the directory 'mtp-test-files'
		objId, isDir, err = GetObjectUsingPath(dev, sid, "mtp-test-files")

		So(err, ShouldBeNil)
		So(objId, ShouldEqual, objId)
		So(isDir, ShouldEqual, true)

		// test the file '/mtp-test-files/a.txt'
		objId, isDir, err = GetObjectUsingPath(dev, sid, "/mtp-test-files/a.txt")

		So(err, ShouldBeNil)
		So(objId, ShouldBeGreaterThan, 0)
		So(isDir, ShouldEqual, false)

		// test the file '/mtp-test-files/mock_dir1/a.txt'
		objId, isDir, err = GetObjectUsingPath(dev, sid, "/mtp-test-files/mock_dir1/a.txt")

		So(err, ShouldBeNil)
		So(objId, ShouldBeGreaterThan, 0)
		So(isDir, ShouldEqual, false)

		// test the file 'mtp-test-files/a.txt'
		objId, isDir, err = GetObjectUsingPath(dev, sid, "mtp-test-files/a.txt")

		So(err, ShouldBeNil)
		So(objId, ShouldEqual, objId)
		So(isDir, ShouldEqual, false)

		// test the file 'mtp-test-files/a.txt/'
		objId, isDir, err = GetObjectUsingPath(dev, sid, "mtp-test-files/a.txt/")

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

		// test the file 'mtp-test-files/b'
		objId, isDir, err = GetObjectUsingPath(dev, sid, "mtp-test-files/b")

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objId, ShouldEqual, 0)
		So(isDir, ShouldEqual, false)

		// test the file 'mtp-test-files/b'
		objId, isDir, err = GetObjectUsingPath(dev, sid, "mtp-test-files/a.txt/1")

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
		// test the directory '/mtp-test-files'
		objId, isDir, err := GetObjectUsingParentId(dev, sid, ParentObjectId, "mtp-test-files")

		So(err, ShouldBeNil)
		So(objId, ShouldBeGreaterThan, 0)
		So(isDir, ShouldEqual, true)

		// test the file '/mtp-test-files/a.txt'
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

		// test the file '/mtp-test-files'
		objId, isDir, err = GetObjectUsingParentId(dev, sid, ParentObjectId, "/mtp-test-files")

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, FileNotFoundError{})
		So(objId, ShouldEqual, 0)
		So(isDir, ShouldEqual, false)

		// test the file 'mtp-test-files/'
		objId, isDir, err = GetObjectUsingParentId(dev, sid, ParentObjectId, "mtp-test-files/")

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
		// test the directory '/mtp-test-files'
		exists := FileExists(dev, sid, "/mtp-test-files/")
		So(exists, ShouldEqual, true)

		// test the file '/mtp-test-files/a.txt'
		exists = FileExists(dev, sid, "/mtp-test-files/a.txt")
		So(exists, ShouldEqual, true)

		// test the directory 'mtp-test-files/'
		exists = FileExists(dev, sid, "mtp-test-files/")
		So(exists, ShouldEqual, true)

		// test the directory 'mtp-test-files'
		exists = FileExists(dev, sid, "mtp-test-files")
		So(exists, ShouldEqual, true)

		// test the file '/mtp-test-files/a.txt/'
		exists = FileExists(dev, sid, "/mtp-test-files/a.txt/")
		So(exists, ShouldEqual, true)
	})

	Convey("Testing non existing file | FileExists | Should throw error", t, func() {
		// test the directory '/fake'
		exists := FileExists(dev, sid, "/fake/")
		So(exists, ShouldEqual, false)

		// test the file '/mtp-test-files/fake.txt'
		exists = FileExists(dev, sid, "/mtp-test-files/fake.txt")
		So(exists, ShouldEqual, false)
	})

	Dispose(dev)
}

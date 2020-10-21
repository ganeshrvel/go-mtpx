package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"testing"
)

func TestGetObjectFromPath(t *testing.T) {
	dev, err := Initialize(Init{})
	if err != nil {
		log.Panic(err)
	}

	storages, err := FetchStorages(dev)
	if err != nil {
		log.Panic(err)
	}

	sid := storages[0].sid

	Convey("Testing valid file | GetObjectFromPath", t, func() {
		// test the directory '/mtp-test-files'
		objId, isDir, err := GetObjectFromPath(dev, sid, "/mtp-test-files")

		So(err, ShouldBeNil)
		So(objId, ShouldBeGreaterThan, 0)
		So(isDir, ShouldEqual, true)

		// test the directory '/mtp-test-files/'
		objId, isDir, err = GetObjectFromPath(dev, sid, "/mtp-test-files/")

		So(err, ShouldBeNil)
		So(objId, ShouldEqual, objId)
		So(isDir, ShouldEqual, true)

		// test the directory 'mtp-test-files/'
		objId, isDir, err = GetObjectFromPath(dev, sid, "mtp-test-files/")

		So(err, ShouldBeNil)
		So(objId, ShouldEqual, objId)
		So(isDir, ShouldEqual, true)

		// test the directory 'mtp-test-files'
		objId, isDir, err = GetObjectFromPath(dev, sid, "mtp-test-files")

		So(err, ShouldBeNil)
		So(objId, ShouldEqual, objId)
		So(isDir, ShouldEqual, true)

		// test the file '/mtp-test-files/a.txt'
		objId, isDir, err = GetObjectFromPath(dev, sid, "/mtp-test-files/a.txt")

		So(err, ShouldBeNil)
		So(objId, ShouldBeGreaterThan, 0)
		So(isDir, ShouldEqual, false)

		// test the file '/mtp-test-files/mock_dir1/a.txt'
		objId, isDir, err = GetObjectFromPath(dev, sid, "/mtp-test-files/mock_dir1/a.txt")

		So(err, ShouldBeNil)
		So(objId, ShouldBeGreaterThan, 0)
		So(isDir, ShouldEqual, false)

		// test the file 'mtp-test-files/a.txt'
		objId, isDir, err = GetObjectFromPath(dev, sid, "mtp-test-files/a.txt")

		So(err, ShouldBeNil)
		So(objId, ShouldEqual, objId)
		So(isDir, ShouldEqual, false)

		// test the file 'mtp-test-files/a.txt/'
		objId, isDir, err = GetObjectFromPath(dev, sid, "mtp-test-files/a.txt/")

		So(err, ShouldBeNil)
		So(objId, ShouldEqual, objId)
		So(isDir, ShouldEqual, false)
	})

	Convey("Testing non exisiting file | GetObjectFromPath | It should throw an error", t, func() {
		// test the file 'fake_file'
		objId, isDir, err := GetObjectFromPath(dev, sid, "fake_file")

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objId, ShouldEqual, 0)
		So(isDir, ShouldEqual, false)

		// test the file 'mtp-test-files/b'
		objId, isDir, err = GetObjectFromPath(dev, sid, "mtp-test-files/b")

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objId, ShouldEqual, 0)
		So(isDir, ShouldEqual, false)

		// test the file 'mtp-test-files/b'
		objId, isDir, err = GetObjectFromPath(dev, sid, "mtp-test-files/a.txt/1")

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objId, ShouldEqual, 0)
		So(isDir, ShouldEqual, false)
	})

	Dispose(dev)
}

func TestGetObjectFromParentIdAndFilename(t *testing.T) {
	dev, err := Initialize(Init{})
	if err != nil {
		log.Panic(err)
	}

	storages, err := FetchStorages(dev)
	if err != nil {
		log.Panic(err)
	}

	sid := storages[0].sid

	Convey("Testing valid file | GetObjectFromParentIdAndFilename", t, func() {
		// test the directory '/mtp-test-files'
		objId, isDir, err := GetObjectFromParentIdAndFilename(dev, sid, ParentObjectId, "mtp-test-files")

		So(err, ShouldBeNil)
		So(objId, ShouldBeGreaterThan, 0)
		So(isDir, ShouldEqual, true)

		// test the file '/mtp-test-files/a.txt'
		objId, isDir, err = GetObjectFromParentIdAndFilename(dev, sid, objId, "a.txt")

		So(err, ShouldBeNil)
		So(objId, ShouldEqual, objId)
		So(isDir, ShouldEqual, false)
	})

	Convey("Testing non exisiting file | GetObjectFromParentIdAndFilename | It should throw an error", t, func() {
		// test the file 'fake_file'
		objId, isDir, err := GetObjectFromParentIdAndFilename(dev, sid, ParentObjectId, "fake_file")

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, FileNotFoundError{})
		So(objId, ShouldEqual, 0)
		So(isDir, ShouldEqual, false)

		// test the file '/mtp-test-files'
		objId, isDir, err = GetObjectFromParentIdAndFilename(dev, sid, ParentObjectId, "/mtp-test-files")

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, FileNotFoundError{})
		So(objId, ShouldEqual, 0)
		So(isDir, ShouldEqual, false)

		// test the file 'mtp-test-files/'
		objId, isDir, err = GetObjectFromParentIdAndFilename(dev, sid, ParentObjectId, "mtp-test-files/")

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

	Convey("Testing valid file | filepath | FileExists", t, func() {
		// test the directory '/mtp-test-files'
		exists, isDir, objectId := FileExists(dev, sid, 0, "/mtp-test-files/")
		So(exists, ShouldEqual, true)
		So(objectId, ShouldBeGreaterThan, 0)
		So(isDir, ShouldEqual, true)

		// test the file '/mtp-test-files/a.txt'
		exists, isDir, objectId = FileExists(dev, sid, 0, "/mtp-test-files/a.txt")
		So(exists, ShouldEqual, true)
		So(objectId, ShouldBeGreaterThan, 0)
		So(isDir, ShouldEqual, false)

		// test the directory 'mtp-test-files/'
		exists, isDir, objectId = FileExists(dev, sid, 0, "mtp-test-files/")
		So(exists, ShouldEqual, true)
		So(objectId, ShouldBeGreaterThan, 0)
		So(isDir, ShouldEqual, true)

		// test the directory 'mtp-test-files'
		exists, isDir, objectId = FileExists(dev, sid, 0, "mtp-test-files")
		So(exists, ShouldEqual, true)
		So(objectId, ShouldBeGreaterThan, 0)
		So(isDir, ShouldEqual, true)

		// test the file '/mtp-test-files/a.txt/'
		exists, isDir, objectId = FileExists(dev, sid, 0, "/mtp-test-files/a.txt/")
		So(exists, ShouldEqual, true)
		So(objectId, ShouldBeGreaterThan, 0)
		So(isDir, ShouldEqual, false)
	})

	Convey("Testing valid file | objectId | FileExists", t, func() {
		// test the directory '/mtp-test-files'
		objectId, isDir, err := GetObjectFromPath(dev, sid, "/mtp-test-files")
		So(err, ShouldBeNil)
		So(objectId, ShouldBeGreaterThan, 0)
		So(isDir, ShouldEqual, true)

		exists, isDir, _objectId := FileExists(dev, sid, objectId, "/mtp-test-files")
		So(exists, ShouldEqual, true)
		So(objectId, ShouldEqual, _objectId)
		So(isDir, ShouldEqual, true)

		// test the file '/mtp-test-files/a.txt/'
		objectId, isDir, err = GetObjectFromPath(dev, sid, "/mtp-test-files/a.txt")
		So(err, ShouldBeNil)
		So(objectId, ShouldBeGreaterThan, 0)
		So(isDir, ShouldEqual, false)

		exists, isDir, _objectId = FileExists(dev, sid, objectId, "/mtp-test-files/a.txt")
		So(exists, ShouldEqual, true)
		So(objectId, ShouldEqual, _objectId)
		So(isDir, ShouldEqual, false)
	})

	Convey("Testing non existing file | FileExists | Should throw error", t, func() {
		// test the directory '/fake'
		exists, isDir, objectId := FileExists(dev, sid, 0, "/fake/")
		So(exists, ShouldEqual, false)
		So(objectId, ShouldEqual, 0)
		So(isDir, ShouldEqual, false)

		// test the file '/mtp-test-files/fake.txt'
		exists, isDir, objectId = FileExists(dev, sid, 0, "/mtp-test-files/fake.txt")
		So(exists, ShouldEqual, false)
		So(objectId, ShouldEqual, 0)
		So(isDir, ShouldEqual, false)
	})

	Dispose(dev)
}

func TestGetObjectFromObjectIdOrPath(t *testing.T) {
	dev, err := Initialize(Init{})
	if err != nil {
		log.Panic(err)
	}

	storages, err := FetchStorages(dev)
	if err != nil {
		log.Panic(err)
	}

	sid := storages[0].sid

	Convey("Testing valid files | GetObjectFromObjectIdOrPath", t, func() {
		// objectId=0 && fullPath="/mtp-test-files/"
		objectId, isDir, err := GetObjectFromObjectIdOrPath(dev, sid, 0, "/mtp-test-files/")
		So(err, ShouldBeNil)
		So(isDir, ShouldEqual, true)
		So(objectId, ShouldBeGreaterThan, 0)

		// objectId=0 && fullPath="mtp-test-files/"
		objectId, isDir, err = GetObjectFromObjectIdOrPath(dev, sid, 0, "mtp-test-files/")
		So(err, ShouldBeNil)
		So(isDir, ShouldEqual, true)
		So(objectId, ShouldBeGreaterThan, 0)

		// objectId=0 && fullPath="mtp-test-files"
		objectId, isDir, err = GetObjectFromObjectIdOrPath(dev, sid, 0, "mtp-test-files")
		So(err, ShouldBeNil)
		So(isDir, ShouldEqual, true)
		So(objectId, ShouldBeGreaterThan, 0)

		// objectId=parentId && fullPath="mtp-test-files"
		objectId, isDir, err = GetObjectFromObjectIdOrPath(dev, sid, ParentObjectId, "mtp-test-files")

		So(err, ShouldBeNil)
		So(isDir, ShouldEqual, true)
		So(objectId, ShouldEqual, ParentObjectId)

		// objectId=parentId && fullPath=""
		objectId, isDir, err = GetObjectFromObjectIdOrPath(dev, sid, ParentObjectId, "")
		So(err, ShouldBeNil)
		So(isDir, ShouldEqual, true)
		So(objectId, ShouldEqual, ParentObjectId)
	})

	Convey("Testing invalid files | GetObjectFromObjectIdOrPath", t, func() {
		// objectId=0 && fullPath=""
		objectId, isDir, err := GetObjectFromObjectIdOrPath(dev, sid, 0, "")
		So(err, ShouldBeError)
		So(isDir, ShouldEqual, false)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objectId, ShouldEqual, 0)

		// objectId=fake && fullPath=""
		objectId, isDir, err = GetObjectFromObjectIdOrPath(dev, sid, 1234567, "")
		So(err, ShouldBeError)
		So(isDir, ShouldEqual, false)
		So(err, ShouldHaveSameTypeAs, FileObjectError{})
		So(objectId, ShouldEqual, 0)

		// objectId=0 && fullPath="/fake"
		objectId, isDir, err = GetObjectFromObjectIdOrPath(dev, sid, 0, "/fake")
		So(err, ShouldBeError)
		So(isDir, ShouldEqual, false)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objectId, ShouldEqual, 0)
	})

	Dispose(dev)
}

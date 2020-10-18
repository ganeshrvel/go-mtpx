package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"testing"
)

func TestFetchFile(t *testing.T) {
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
		So(objId, ShouldBeGreaterThan, 0)
		So(isDir, ShouldEqual, false)
	})

	Convey("Testing non exisiting file | GetObjectIdFromFilename | It should throw an error", t, func() {
		objId, isDir, err := GetObjectIdFromFilename(dev, sid, ParentObjectId, "fake_file")

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, FileNotFoundError{})
		So(objId, ShouldEqual, 0)
		So(isDir, ShouldEqual, false)
	})

	Dispose(dev)
}

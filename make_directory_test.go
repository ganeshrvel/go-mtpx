package mtpx

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

	sid := storages[0].Sid

	var _objectId uint32
	Convey("Creating a new dir | MakeDirectory", t, func() {
		// test the directory '/mtp-test-files/temp_dir/test-MakeDirectory'
		objectId, err := MakeDirectory(dev, sid, "/mtp-test-files/temp_dir/test-MakeDirectory")

		_objectId = objectId

		So(err, ShouldBeNil)
		So(objectId, ShouldBeGreaterThan, 0)
	})

	Convey("Testing MakeDirectory for an existing directory | MakeDirectory", t, func() {
		// test the directory '/mtp-test-files/temp_dir/test-MakeDirectory'
		objectId, err := MakeDirectory(dev, sid, "/mtp-test-files/temp_dir/test-MakeDirectory")

		So(err, ShouldBeNil)
		So(objectId, ShouldEqual, _objectId)
	})

	Convey("Testing fullpath='/' | MakeDirectory", t, func() {
		// test the directory '/'
		objectId, err := MakeDirectory(dev, sid, "/")

		So(err, ShouldBeNil)
		So(objectId, ShouldEqual, ParentObjectId)
	})

	Convey("Testing fullpath='' | MakeDirectory", t, func() {
		// test the directory ''
		objectId, err := MakeDirectory(dev, sid, "")

		So(err, ShouldBeNil)
		So(objectId, ShouldEqual, ParentObjectId)
	})

	Convey("Creating a new random dir | MakeDirectory", t, func() {
		// test the directory '/mtp-test-files/temp_dir/test-MakeDirectory/{random}'
		fullpath := fmt.Sprintf("/mtp-test-files/temp_dir/test-MakeDirectory/%x", rand.Int31())

		objectId, err := MakeDirectory(dev, sid, fullpath)

		So(err, ShouldBeNil)
		So(objectId, ShouldBeGreaterThan, 0)

		exists, fi := FileExists(dev, sid, FileProp{0, fullpath})

		So(err, ShouldBeNil)
		So(exists, ShouldEqual, true)
		So(fi.ObjectId, ShouldEqual, objectId)
		So(fi.IsDir, ShouldEqual, true)
	})

	Convey("filename in the path | 1 | MakeDirectory | It should throw an error", t, func() {
		// test the directory '/mtp-test-files/a.txt/folder'
		objectId, err := MakeDirectory(dev, sid, "/mtp-test-files/a.txt/folder")

		So(err, ShouldBeError)
		So(objectId, ShouldEqual, 0)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
	})

	Convey("filename in the path | 2 | MakeDirectory | It should throw an error", t, func() {
		// test the directory '/mtp-test-files/a.txt/folder'
		objectId, err := MakeDirectory(dev, sid, "/mtp-test-files/a.txt")

		So(err, ShouldBeError)
		So(objectId, ShouldEqual, 0)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
	})

	Dispose(dev)
}

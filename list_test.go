package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"log"
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
		// test the directory '/mocks'
		files, err := ListDirectory(dev, sid, ParentObjectId, "/mocks")

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
		// test the directory '/mocks'
		files, err := ListDirectory(dev, sid, 0, "/mocks")

		So(err, ShouldBeNil)
		So(files, ShouldNotBeNil)
		So(len(*files), ShouldEqual, 4)

		// test the directory '/mocks/'
		files, err = ListDirectory(dev, sid, 0, "/mocks/")

		So(err, ShouldBeNil)
		So(files, ShouldNotBeNil)
		So(len(*files), ShouldEqual, 4)

		// test the directory 'mocks/'
		files, err = ListDirectory(dev, sid, 0, "mocks/")

		So(err, ShouldBeNil)
		So(files, ShouldNotBeNil)
		So(len(*files), ShouldEqual, 4)

		// test the directory 'mocks'
		files, err = ListDirectory(dev, sid, 0, "mocks")

		So(err, ShouldBeNil)
		So(files, ShouldNotBeNil)
		So(len(*files), ShouldEqual, 4)

		// test the directory 'mocks/mock_dir3/'
		files, err = ListDirectory(dev, sid, 0, "mocks/mock_dir3/")

		So(err, ShouldBeNil)
		So(files, ShouldNotBeNil)
		So(len(*files), ShouldEqual, 5)
	})

	Convey("Testing valid file | ListDirectory", t, func() {
		// test the directory '/mocks/mock_dir1/1'
		files, err := ListDirectory(dev, sid, 0, "/mocks/mock_dir1/1")

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
		So(_files[0].FullPath, ShouldEqual, "/mocks/mock_dir1/1/a.txt")
		So(_files[0].ModTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)
	})

	Convey("Testing non exisiting file | ListDirectory | It should throw an error", t, func() {
		// test the directory '/fake'
		files, err := ListDirectory(dev, sid, 0, "/fake")

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(files, ShouldBeNil)

		// test the directory '/mocks'
		files, err = ListDirectory(dev, sid, 0, "/mocks/a.txt")

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(files, ShouldBeNil)

		// test the directory '/mocks/fake'
		files, err = ListDirectory(dev, sid, 0, "/mocks/fake")

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(files, ShouldBeNil)
	})

	Dispose(dev)
}

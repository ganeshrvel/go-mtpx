package mtpx

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

	sid := storages[0].Sid

	Convey("Testing valid file | GetObjectFromPath", t, func() {
		// test the directory '/mtp-test-files'
		fi, err := GetObjectFromPath(dev, sid, "/mtp-test-files")

		So(err, ShouldBeNil)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, true)
		So(fi.FullPath, ShouldEqual, "/mtp-test-files")

		// test the directory '/mtp-test-files/'
		fi, err = GetObjectFromPath(dev, sid, "/mtp-test-files/")

		So(err, ShouldBeNil)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, true)
		So(fi.FullPath, ShouldEqual, "/mtp-test-files")

		// test the directory 'mtp-test-files/'
		fi, err = GetObjectFromPath(dev, sid, "mtp-test-files/")

		So(err, ShouldBeNil)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, true)
		So(fi.FullPath, ShouldEqual, "/mtp-test-files")

		// test the directory 'mtp-test-files'
		fi, err = GetObjectFromPath(dev, sid, "mtp-test-files")

		So(err, ShouldBeNil)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, true)
		So(fi.FullPath, ShouldEqual, "/mtp-test-files")

		// test the file '/mtp-test-files/a.txt'
		fi, err = GetObjectFromPath(dev, sid, "/mtp-test-files/a.txt")

		So(err, ShouldBeNil)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, false)
		So(fi.FullPath, ShouldEqual, "/mtp-test-files/a.txt")

		// test the file '/mtp-test-files/mock_dir1/a.txt'
		fi, err = GetObjectFromPath(dev, sid, "/mtp-test-files/mock_dir1/a.txt")

		So(err, ShouldBeNil)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, false)
		So(fi.FullPath, ShouldEqual, "/mtp-test-files/mock_dir1/a.txt")

		// test the file 'mtp-test-files/a.txt'
		fi, err = GetObjectFromPath(dev, sid, "mtp-test-files/a.txt")

		So(err, ShouldBeNil)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, false)
		So(fi.FullPath, ShouldEqual, "/mtp-test-files/a.txt")

		// test the file 'mtp-test-files/a.txt/'
		fi, err = GetObjectFromPath(dev, sid, "mtp-test-files/a.txt/")

		So(err, ShouldBeNil)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, false)
		So(fi.FullPath, ShouldEqual, "/mtp-test-files/a.txt")

		// test the file ''
		fi, err = GetObjectFromPath(dev, sid, "")

		So(err, ShouldBeError)
		So(fi, ShouldBeNil)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
	})

	Convey("Testing non exisiting file | GetObjectFromPath | It should throw an error", t, func() {
		// test the file 'fake_file'
		fi, err := GetObjectFromPath(dev, sid, "fake_file")

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(fi, ShouldBeNil)

		// test the file 'mtp-test-files/b'
		fi, err = GetObjectFromPath(dev, sid, "mtp-test-files/b")

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(fi, ShouldBeNil)

		// test the file 'mtp-test-files/b'
		fi, err = GetObjectFromPath(dev, sid, "mtp-test-files/a.txt/1")

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(fi, ShouldBeNil)
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

	sid := storages[0].Sid

	Convey("Testing valid file | GetObjectFromParentIdAndFilename", t, func() {
		// test the directory '/mtp-test-files'
		fi, err := GetObjectFromParentIdAndFilename(dev, sid, ParentObjectId, "mtp-test-files")

		So(err, ShouldBeNil)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, true)

		// test the file '/mtp-test-files/a.txt'
		fi, err = GetObjectFromParentIdAndFilename(dev, sid, fi.ObjectId, "a.txt")

		So(err, ShouldBeNil)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, false)
	})

	Convey("Testing non exisiting file | GetObjectFromParentIdAndFilename | It should throw an error", t, func() {
		// test the file 'fake_file'
		fi, err := GetObjectFromParentIdAndFilename(dev, sid, ParentObjectId, "fake_file")

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, FileNotFoundError{})
		So(fi, ShouldBeNil)

		// test the file '/mtp-test-files'
		fi, err = GetObjectFromParentIdAndFilename(dev, sid, ParentObjectId, "/mtp-test-files")

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, FileNotFoundError{})
		So(fi, ShouldBeNil)

		// test the file 'mtp-test-files/'
		fi, err = GetObjectFromParentIdAndFilename(dev, sid, ParentObjectId, "mtp-test-files/")

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, FileNotFoundError{})
		So(fi, ShouldBeNil)
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

	sid := storages[0].Sid

	Convey("Testing valid file | filepath | FileExists", t, func() {
		// test the directory '/mtp-test-files'
		exists, fi := FileExists(dev, sid, 0, "/mtp-test-files/")
		So(exists, ShouldEqual, true)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, true)

		// test the file '/mtp-test-files/a.txt'
		exists, fi = FileExists(dev, sid, 0, "/mtp-test-files/a.txt")
		So(exists, ShouldEqual, true)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, false)

		// test the directory 'mtp-test-files/'
		exists, fi = FileExists(dev, sid, 0, "mtp-test-files/")
		So(exists, ShouldEqual, true)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, true)

		// test the directory 'mtp-test-files'
		exists, fi = FileExists(dev, sid, 0, "mtp-test-files")
		So(exists, ShouldEqual, true)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, true)

		// test the file '/mtp-test-files/a.txt/'
		exists, fi = FileExists(dev, sid, 0, "/mtp-test-files/a.txt/")
		So(exists, ShouldEqual, true)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, false)
	})

	Convey("Testing valid file | objectId | FileExists", t, func() {
		// test the directory '/mtp-test-files'
		fi, err := GetObjectFromPath(dev, sid, "/mtp-test-files")
		So(err, ShouldBeNil)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, true)

		_objectId := fi.ObjectId

		exists, fi := FileExists(dev, sid, _objectId, "/mtp-test-files")
		So(exists, ShouldEqual, true)
		So(fi.ObjectId, ShouldEqual, _objectId)
		So(fi.IsDir, ShouldEqual, true)

		// test the file '/mtp-test-files/a.txt/'
		fi, err = GetObjectFromPath(dev, sid, "/mtp-test-files/a.txt")
		So(err, ShouldBeNil)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, false)

		_objectId = fi.ObjectId

		exists, fi = FileExists(dev, sid, _objectId, "/mtp-test-files/a.txt")
		So(exists, ShouldEqual, true)
		So(fi.ObjectId, ShouldEqual, _objectId)
		So(fi.IsDir, ShouldEqual, false)
	})

	Convey("Testing non existing file | FileExists | Should throw error", t, func() {
		// test the directory '/fake'
		exists, fi := FileExists(dev, sid, 0, "/fake/")
		So(exists, ShouldEqual, false)
		So(fi, ShouldBeNil)

		// test the file '/mtp-test-files/fake.txt'
		exists, fi = FileExists(dev, sid, 0, "/mtp-test-files/fake.txt")
		So(exists, ShouldEqual, false)
		So(fi, ShouldBeNil)
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

	sid := storages[0].Sid

	Convey("Testing valid files | GetObjectFromObjectIdOrPath", t, func() {
		// objectId=0 && fullPath="/mtp-test-files/"
		fi, err := GetObjectFromObjectIdOrPath(dev, sid, 0, "/mtp-test-files/")
		So(err, ShouldBeNil)
		So(fi.IsDir, ShouldEqual, true)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)

		// objectId=0 && fullPath="mtp-test-files/"
		fi, err = GetObjectFromObjectIdOrPath(dev, sid, 0, "mtp-test-files/")
		So(err, ShouldBeNil)
		So(fi.IsDir, ShouldEqual, true)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)

		// objectId=0 && fullPath="mtp-test-files"
		fi, err = GetObjectFromObjectIdOrPath(dev, sid, 0, "mtp-test-files")
		So(err, ShouldBeNil)
		So(fi.IsDir, ShouldEqual, true)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)

		// objectId=parentId && fullPath="mtp-test-files"
		fi, err = GetObjectFromObjectIdOrPath(dev, sid, ParentObjectId, "mtp-test-files")

		So(err, ShouldBeNil)
		So(fi.IsDir, ShouldEqual, true)
		So(fi.ObjectId, ShouldEqual, ParentObjectId)

		// objectId=parentId && fullPath=""
		fi, err = GetObjectFromObjectIdOrPath(dev, sid, ParentObjectId, "")
		So(err, ShouldBeNil)
		So(fi.IsDir, ShouldEqual, true)
		So(fi.ObjectId, ShouldEqual, ParentObjectId)
	})

	Convey("Testing invalid files | GetObjectFromObjectIdOrPath", t, func() {
		// objectId=0 && fullPath=""
		fi, err := GetObjectFromObjectIdOrPath(dev, sid, 0, "")
		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(fi, ShouldBeNil)

		// objectId=fake && fullPath=""
		fi, err = GetObjectFromObjectIdOrPath(dev, sid, 1234567, "")
		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, FileObjectError{})
		So(fi, ShouldBeNil)

		// objectId=0 && fullPath="/fake"
		fi, err = GetObjectFromObjectIdOrPath(dev, sid, 0, "/fake")
		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(fi, ShouldBeNil)
	})

	Dispose(dev)
}

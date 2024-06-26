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
		if fi.IsDir {
			So(fi.Size, ShouldEqual, 0)
		} else {
			So(fi.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}

		// test the directory '/mtp-test-files/'
		fi, err = GetObjectFromPath(dev, sid, "/mtp-test-files/")

		So(err, ShouldBeNil)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, true)
		So(fi.FullPath, ShouldEqual, "/mtp-test-files")
		if fi.IsDir {
			So(fi.Size, ShouldEqual, 0)
		} else {
			So(fi.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}

		// test the directory 'mtp-test-files/'
		fi, err = GetObjectFromPath(dev, sid, "mtp-test-files/")

		So(err, ShouldBeNil)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, true)
		So(fi.FullPath, ShouldEqual, "/mtp-test-files")
		if fi.IsDir {
			So(fi.Size, ShouldEqual, 0)
		} else {
			So(fi.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}

		// test the directory 'mtp-test-files'
		fi, err = GetObjectFromPath(dev, sid, "mtp-test-files")

		So(err, ShouldBeNil)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, true)
		So(fi.FullPath, ShouldEqual, "/mtp-test-files")
		if fi.IsDir {
			So(fi.Size, ShouldEqual, 0)
		} else {
			So(fi.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}

		// test the file '/mtp-test-files/a.txt'
		fi, err = GetObjectFromPath(dev, sid, "/mtp-test-files/a.txt")

		So(err, ShouldBeNil)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, false)
		So(fi.FullPath, ShouldEqual, "/mtp-test-files/a.txt")
		if fi.IsDir {
			So(fi.Size, ShouldEqual, 0)
		} else {
			So(fi.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}

		// test the file '/mtp-test-files/mock_dir1/a.txt'
		fi, err = GetObjectFromPath(dev, sid, "/mtp-test-files/mock_dir1/a.txt")

		So(err, ShouldBeNil)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, false)
		So(fi.FullPath, ShouldEqual, "/mtp-test-files/mock_dir1/a.txt")
		if fi.IsDir {
			So(fi.Size, ShouldEqual, 0)
		} else {
			So(fi.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}

		// test the file 'mtp-test-files/a.txt'
		fi, err = GetObjectFromPath(dev, sid, "mtp-test-files/a.txt")

		So(err, ShouldBeNil)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, false)
		So(fi.FullPath, ShouldEqual, "/mtp-test-files/a.txt")
		if fi.IsDir {
			So(fi.Size, ShouldEqual, 0)
		} else {
			So(fi.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}

		// test the file 'mtp-test-files/a.txt/'
		fi, err = GetObjectFromPath(dev, sid, "mtp-test-files/a.txt/")

		So(err, ShouldBeNil)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, false)
		So(fi.FullPath, ShouldEqual, "/mtp-test-files/a.txt")
		if fi.IsDir {
			So(fi.Size, ShouldEqual, 0)
		} else {
			So(fi.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}

		// test the file ''
		fi, err = GetObjectFromPath(dev, sid, "")

		So(err, ShouldBeError)
		So(fi, ShouldBeNil)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
	})

	Convey("Testing mixed case file names | GetObjectFromPath | It should throw an error", t, func() {
		// test the file 'MTP-TEST-FILES/A.TXT/'
		fi, err := GetObjectFromPath(dev, sid, "MTP-TEST-FILES/A.TXT/")

		So(err, ShouldBeNil)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, false)
		So(fi.FullPath, ShouldEqual, "/MTP-TEST-FILES/A.TXT")
		if fi.IsDir {
			So(fi.Size, ShouldEqual, 0)
		} else {
			So(fi.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}

		// test the file 'mtp-test-files/A.TXT/'
		fi, err = GetObjectFromPath(dev, sid, "mtp-test-files/A.TXT/")
		So(err, ShouldBeNil)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, false)
		So(fi.FullPath, ShouldEqual, "/mtp-test-files/A.TXT")
		if fi.IsDir {
			So(fi.Size, ShouldEqual, 0)
		} else {
			So(fi.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}
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

		// test the file 'mtp-test-files/a.txt/1'
		fi, err = GetObjectFromPath(dev, sid, "mtp-test-files/a.txt/1")

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(fi, ShouldBeNil)

		// test the file 'mtp-test-files/A.TXT/1'
		fi, err = GetObjectFromPath(dev, sid, "mtp-test-files/A.TXT/1")

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
		if fi.IsDir {
			So(fi.Size, ShouldEqual, 0)
		} else {
			So(fi.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}

		// test the file '/mtp-test-files/a.txt'
		fi1, err := GetObjectFromParentIdAndFilename(dev, sid, fi.ObjectId, "a.txt")
		So(err, ShouldBeNil)
		So(fi1.ObjectId, ShouldBeGreaterThan, 0)
		So(fi1.IsDir, ShouldEqual, false)
		if fi1.IsDir {
			So(fi1.Size, ShouldEqual, 0)
		} else {
			So(fi1.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}

		// test filename with different case '/mtp-test-files/A.TXT'
		fi2, err := GetObjectFromParentIdAndFilename(dev, sid, fi.ObjectId, "A.TXT")
		So(err, ShouldBeNil)
		So(fi2.ObjectId, ShouldBeGreaterThan, 0)
		So(fi2.IsDir, ShouldEqual, false)
		if fi2.IsDir {
			So(fi2.Size, ShouldEqual, 0)
		} else {
			So(fi2.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}
	})

	Convey("Testing non exisiting file | GetObjectFromParentIdAndFilename | It should throw an error", t, func() {
		// test the file 'fake_file'
		fi, err := GetObjectFromParentIdAndFilename(dev, sid, ParentObjectId, "fake_file")

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, FileNotFoundError{})
		So(fi, ShouldBeNil)

		// slashes inside the file name is invalid
		// test the file '/mtp-test-files'
		fi, err = GetObjectFromParentIdAndFilename(dev, sid, ParentObjectId, "/mtp-test-files")

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, FileNotFoundError{})
		So(fi, ShouldBeNil)

		// slashes inside the file name is invalid
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
		fc, err := FileExists(dev, sid, []FileProp{{0, "/mtp-test-files/"}})
		So(err, ShouldBeNil)
		fi := fc[0].FileInfo

		So(fc[0].Exists, ShouldEqual, true)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, true)
		if fi.IsDir {
			So(fi.Size, ShouldEqual, 0)
		} else {
			So(fi.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}

		// test the file '/mtp-test-files/a.txt'
		fc, err = FileExists(dev, sid, []FileProp{{0, "/mtp-test-files/a.txt"}})
		So(err, ShouldBeNil)
		fi1 := fc[0].FileInfo

		So(fc[0].Exists, ShouldEqual, true)
		So(fi1.ObjectId, ShouldBeGreaterThan, 0)
		So(fi1.IsDir, ShouldEqual, false)
		if fi1.IsDir {
			So(fi1.Size, ShouldEqual, 0)
		} else {
			So(fi1.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}

		// test the file '/mtp-test-files/A.TXT'
		fc, err = FileExists(dev, sid, []FileProp{{0, "/mtp-test-files/A.TXT"}})
		So(err, ShouldBeNil)
		fi2 := fc[0].FileInfo

		So(fc[0].Exists, ShouldEqual, true)
		So(fi2.ObjectId, ShouldEqual, fi1.ObjectId)
		So(fi2.IsDir, ShouldEqual, false)
		if fi2.IsDir {
			So(fi2.Size, ShouldEqual, 0)
		} else {
			So(fi2.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}

		// test the file '/MTP-TEST-FILES/A.TXT'
		fc, err = FileExists(dev, sid, []FileProp{{0, "/MTP-TEST-FILES/A.TXT"}})
		So(err, ShouldBeNil)
		fi3 := fc[0].FileInfo

		So(fc[0].Exists, ShouldEqual, true)
		So(fi3.ObjectId, ShouldEqual, fi1.ObjectId)
		So(fi3.IsDir, ShouldEqual, false)
		if fi3.IsDir {
			So(fi3.Size, ShouldEqual, 0)
		} else {
			So(fi3.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}

		// test the file '/mtp-test-files/A.TXT'
		fc, err = FileExists(dev, sid, []FileProp{{0, "/MTP-TEST-FILES/A.TXT"}})
		So(err, ShouldBeNil)
		fi4 := fc[0].FileInfo

		So(fc[0].Exists, ShouldEqual, true)
		So(fi4.ObjectId, ShouldEqual, fi1.ObjectId)
		So(fi4.IsDir, ShouldEqual, false)
		if fi4.IsDir {
			So(fi4.Size, ShouldEqual, 0)
		} else {
			So(fi4.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}

		// test the directory 'mtp-test-files/'
		fc, err = FileExists(dev, sid, []FileProp{{0, "mtp-test-files/"}})
		So(err, ShouldBeNil)
		fi = fc[0].FileInfo

		So(fc[0].Exists, ShouldEqual, true)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, true)
		if fi.IsDir {
			So(fi.Size, ShouldEqual, 0)
		} else {
			So(fi.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}

		// test the directory 'mtp-test-files'
		fc, err = FileExists(dev, sid, []FileProp{{0, "mtp-test-files"}})
		So(err, ShouldBeNil)
		fi = fc[0].FileInfo

		So(fc[0].Exists, ShouldEqual, true)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, true)
		if fi.IsDir {
			So(fi.Size, ShouldEqual, 0)
		} else {
			So(fi.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}

		// test the file '/mtp-test-files/a.txt/'
		fc, err = FileExists(dev, sid, []FileProp{{0, "/mtp-test-files/a.txt/"}})
		So(err, ShouldBeNil)
		fi = fc[0].FileInfo

		So(fc[0].Exists, ShouldEqual, true)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, false)
		if fi.IsDir {
			So(fi.Size, ShouldEqual, 0)
		} else {
			So(fi.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}
	})

	Convey("Testing valid file | objectId | FileExists", t, func() {
		// test the directory '/mtp-test-files'
		fi, err := GetObjectFromPath(dev, sid, "/mtp-test-files")
		So(err, ShouldBeNil)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, true)
		if fi.IsDir {
			So(fi.Size, ShouldEqual, 0)
		} else {
			So(fi.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}

		_objectId := fi.ObjectId

		fc, err := FileExists(dev, sid, []FileProp{{_objectId, "/mtp-test-files"}})
		So(err, ShouldBeNil)
		fi = fc[0].FileInfo

		So(fc[0].Exists, ShouldEqual, true)
		So(fi.ObjectId, ShouldEqual, _objectId)
		So(fi.IsDir, ShouldEqual, true)
		if fi.IsDir {
			So(fi.Size, ShouldEqual, 0)
		} else {
			So(fi.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}

		// test the file '/mtp-test-files/a.txt/'
		fi, err = GetObjectFromPath(dev, sid, "/mtp-test-files/a.txt")
		So(err, ShouldBeNil)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		So(fi.IsDir, ShouldEqual, false)
		if fi.IsDir {
			So(fi.Size, ShouldEqual, 0)
		} else {
			So(fi.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}

		_objectId = fi.ObjectId

		fc, err = FileExists(dev, sid, []FileProp{{_objectId, "/mtp-test-files/a.txt"}})
		So(err, ShouldBeNil)
		fi = fc[0].FileInfo

		So(fc[0].Exists, ShouldEqual, true)
		So(fi.ObjectId, ShouldEqual, _objectId)
		So(fi.IsDir, ShouldEqual, false)
		if fi.IsDir {
			So(fi.Size, ShouldEqual, 0)
		} else {
			So(fi.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}
	})

	Convey("Testing multiple valid files | objectIds | FileExists", t, func() {
		fi1, err := GetObjectFromPath(dev, sid, "/mtp-test-files/mock_dir1/a.txt")
		So(err, ShouldBeNil)

		fi2, err := GetObjectFromPath(dev, sid, "/mtp-test-files/a.txt")
		So(err, ShouldBeNil)

		_objectId1 := fi1.ObjectId
		_objectId2 := fi2.ObjectId

		fc, err := FileExists(dev, sid, []FileProp{{ObjectId: _objectId1}, {ObjectId: _objectId2}})
		So(err, ShouldBeNil)
		_fi1 := fc[0].FileInfo

		So(fc[0].Exists, ShouldEqual, true)
		if fi1.IsDir {
			So(fi1.Size, ShouldEqual, 0)
		} else {
			So(fi1.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}
		if fi2.IsDir {
			So(fi2.Size, ShouldEqual, 0)
		} else {
			So(fi2.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}

		So(err, ShouldBeNil)
		_fi2 := fc[1].FileInfo

		So(fc[1].Exists, ShouldEqual, true)
		if fi1.IsDir {
			So(fi1.Size, ShouldEqual, 0)
		} else {
			So(fi1.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}
		if fi2.IsDir {
			So(fi2.Size, ShouldEqual, 0)
		} else {
			So(fi2.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}

		So(_fi1.ObjectId, ShouldEqual, _objectId1)
		So(_fi1.IsDir, ShouldEqual, false)

		So(_fi2.ObjectId, ShouldEqual, _objectId2)
		So(_fi2.IsDir, ShouldEqual, false)
	})

	Convey("Testing multiple valid files | fullPaths | FileExists", t, func() {
		fi1, err := GetObjectFromPath(dev, sid, "/mtp-test-files/mock_dir1/a.txt")
		So(err, ShouldBeNil)

		fi2, err := GetObjectFromPath(dev, sid, "/mtp-test-files/A.TXT")
		So(err, ShouldBeNil)

		_objectId1 := fi1.ObjectId
		_objectId2 := fi2.ObjectId

		fc, err := FileExists(dev, sid, []FileProp{{FullPath: fi1.FullPath}, {FullPath: fi2.FullPath}})
		So(err, ShouldBeNil)
		_fi1 := fc[0].FileInfo

		So(fc[0].Exists, ShouldEqual, true)
		if fi1.IsDir {
			So(fi1.Size, ShouldEqual, 0)
		} else {
			So(fi1.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}
		if fi2.IsDir {
			So(fi2.Size, ShouldEqual, 0)
		} else {
			So(fi2.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}

		_fi2 := fc[1].FileInfo

		So(fc[1].Exists, ShouldEqual, true)
		if fi1.IsDir {
			So(fi1.Size, ShouldEqual, 0)
		} else {
			So(fi1.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}
		if fi2.IsDir {
			So(fi2.Size, ShouldEqual, 0)
		} else {
			So(fi2.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}

		So(_fi1.ObjectId, ShouldEqual, _objectId1)
		So(_fi1.IsDir, ShouldEqual, false)

		So(_fi2.ObjectId, ShouldEqual, _objectId2)
		So(_fi2.IsDir, ShouldEqual, false)
	})

	Convey("Testing non existing file | FileExists | Should throw error", t, func() {
		// test the directory '/fake'
		fc, err := FileExists(dev, sid, []FileProp{{0, "/fake/"}})
		So(err, ShouldBeNil)

		fi := fc[0].FileInfo
		So(fc[0].Exists, ShouldEqual, false)
		So(fi, ShouldBeNil)

		// test the file '/mtp-test-files/fake.txt'
		fc, err = FileExists(dev, sid, []FileProp{{0, "/mtp-test-files/fake.txt"}})
		So(err, ShouldBeNil)

		fi = fc[0].FileInfo
		So(fc[0].Exists, ShouldEqual, false)
		So(fi, ShouldBeNil)

		// test a fake objectId
		fc, err = FileExists(dev, sid, []FileProp{{ObjectId: uint32(987654754)}})
		So(err, ShouldBeNil)

		fi = fc[0].FileInfo
		So(fc[0].Exists, ShouldEqual, false)
		So(fi, ShouldBeNil)
	})

	Convey("Testing non existent and existent files | FileExists | Should throw error", t, func() {
		// test the a valid file and fake file '/fake' and "/mtp-test-files/mock_dir1/a.txt"
		fc, err := FileExists(dev, sid, []FileProp{{FullPath: "/fake/"}, {FullPath: "/mtp-test-files/mock_dir1/a.txt"}})
		So(err, ShouldBeNil)

		fi := fc[0].FileInfo
		So(fc[0].Exists, ShouldEqual, false)
		So(fi, ShouldBeNil)

		fi = fc[1].FileInfo
		So(fc[1].Exists, ShouldEqual, true)
		So(fi, ShouldNotBeNil)

		// test the a valid file and fake file  "/mtp-test-files/mock_dir1" && '/fake'
		fc, err = FileExists(dev, sid, []FileProp{{FullPath: "/mtp-test-files/mock_dir1/a.txt"}, {FullPath: "/fake/"}})
		So(err, ShouldBeNil)

		fi = fc[0].FileInfo
		So(fc[0].Exists, ShouldEqual, true)
		So(fi, ShouldNotBeNil)

		fi = fc[1].FileInfo
		So(fc[1].Exists, ShouldEqual, false)
		So(fi, ShouldBeNil)

		// test the a valid file and fake object id: "/mtp-test-files/mock_dir1" and uint32('987654754')
		fc, err = FileExists(dev, sid, []FileProp{{FullPath: "/mtp-test-files/mock_dir1/a.txt"}, {ObjectId: uint32(987654754)}})
		So(err, ShouldBeNil)

		fi = fc[0].FileInfo
		So(fc[0].Exists, ShouldEqual, true)
		So(fi, ShouldNotBeNil)

		fi = fc[1].FileInfo
		So(fc[1].Exists, ShouldEqual, false)
		So(fi, ShouldBeNil)
	})

	Convey("Testing multiple non existent files | FileExists | Should throw error", t, func() {
		// test the files '/fake' and "/mtp-test-files/fake.txt"
		fc, err := FileExists(dev, sid, []FileProp{{FullPath: "/fake/"}, {FullPath: "/mtp-test-files/fake.txt"}})
		So(err, ShouldBeNil)

		fi := fc[0].FileInfo
		So(fc[0].Exists, ShouldEqual, false)
		So(fi, ShouldBeNil)

		fi = fc[1].FileInfo
		So(fc[1].Exists, ShouldEqual, false)
		So(fi, ShouldBeNil)

		//// test the file '/mtp-test-files/fake.txt' and a fake ObjectId
		fc, err = FileExists(dev, sid, []FileProp{{ObjectId: uint32(987654754)}, {FullPath: "/mtp-test-files/fake.txt"}})
		So(err, ShouldBeNil)

		fi = fc[0].FileInfo
		So(fc[0].Exists, ShouldEqual, false)
		So(fi, ShouldBeNil)

		fi = fc[1].FileInfo
		So(fc[1].Exists, ShouldEqual, false)
		So(fi, ShouldBeNil)

		//// test the multiple fake ObjectIds
		fc, err = FileExists(dev, sid, []FileProp{{ObjectId: uint32(987654754)}, {ObjectId: uint32(987654700)}})
		So(err, ShouldBeNil)

		fi = fc[0].FileInfo
		So(fc[0].Exists, ShouldEqual, false)
		So(fi, ShouldBeNil)

		fi = fc[1].FileInfo
		So(fc[1].Exists, ShouldEqual, false)
		So(fi, ShouldBeNil)
	})

	Dispose(dev)
}

func
TestGetObjectFromObjectIdOrPath(t *testing.T) {
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
		fi, err := GetObjectFromObjectIdOrPath(dev, sid, FileProp{0, "/mtp-test-files/"})
		So(err, ShouldBeNil)
		So(fi.IsDir, ShouldEqual, true)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		if fi.IsDir {
			So(fi.Size, ShouldEqual, 0)
		} else {
			So(fi.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}

		// objectId=0 && fullPath="mtp-test-files/"
		fi, err = GetObjectFromObjectIdOrPath(dev, sid, FileProp{0, "mtp-test-files/"})
		So(err, ShouldBeNil)
		So(fi.IsDir, ShouldEqual, true)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		if fi.IsDir {
			So(fi.Size, ShouldEqual, 0)
		} else {
			So(fi.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}

		// objectId=0 && fullPath="mtp-test-files"
		fi, err = GetObjectFromObjectIdOrPath(dev, sid, FileProp{0, "mtp-test-files"})
		So(err, ShouldBeNil)
		So(fi.IsDir, ShouldEqual, true)
		So(fi.ObjectId, ShouldBeGreaterThan, 0)
		if fi.IsDir {
			So(fi.Size, ShouldEqual, 0)
		} else {
			So(fi.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}

		// objectId=parentId && fullPath="mtp-test-files"
		fi, err = GetObjectFromObjectIdOrPath(dev, sid, FileProp{ParentObjectId, "mtp-test-files"})

		So(err, ShouldBeNil)
		So(fi.IsDir, ShouldEqual, true)
		So(fi.ObjectId, ShouldEqual, ParentObjectId)
		if fi.IsDir {
			So(fi.Size, ShouldEqual, 0)
		} else {
			So(fi.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}

		// objectId=parentId && fullPath=""
		fi, err = GetObjectFromObjectIdOrPath(dev, sid, FileProp{ParentObjectId, ""})
		So(err, ShouldBeNil)
		So(fi.IsDir, ShouldEqual, true)
		So(fi.ObjectId, ShouldEqual, ParentObjectId)
		if fi.IsDir {
			So(fi.Size, ShouldEqual, 0)
		} else {
			So(fi.Size, ShouldBeGreaterThanOrEqualTo, 0)
		}
	})

	Convey("Testing invalid files | GetObjectFromObjectIdOrPath", t, func() {
		// objectId=0 && fullPath=""
		fi, err := GetObjectFromObjectIdOrPath(dev, sid, FileProp{0, ""})
		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(fi, ShouldBeNil)

		// objectId=fake && fullPath=""
		fi, err = GetObjectFromObjectIdOrPath(dev, sid, FileProp{1234567, ""})
		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, FileObjectError{})
		So(fi, ShouldBeNil)

		// objectId=0 && fullPath="/fake"
		fi, err = GetObjectFromObjectIdOrPath(dev, sid, FileProp{0, "/fake"})
		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(fi, ShouldBeNil)
	})

	Dispose(dev)
}

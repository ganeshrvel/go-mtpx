package mtpx

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"math/rand"
	"testing"
	"time"
)

func TestRenameFile(t *testing.T) {
	dev, err := Initialize(Init{})
	if err != nil {
		log.Panic(err)
	}

	storages, err := FetchStorages(dev)
	if err != nil {
		log.Panic(err)
	}

	sid := storages[0].Sid

	Convey("Rename an existing object | using objectId | RenameFile", t, func() {
		// create a random directory
		// test the directory '/mtp-test-files/temp_dir/test-RenameFile/{random}'
		fileName := fmt.Sprintf("/mtp-test-files/temp_dir/test-RenameFile/%x", rand.Int31())
		renameRandFileName := fmt.Sprintf("renamed-%x", rand.Int31())

		objectId, err := MakeDirectory(dev, sid, fileName)

		So(err, ShouldBeNil)
		So(objectId, ShouldBeGreaterThan, 0)

		//rename the object using objectId
		objId, err := RenameFile(dev, sid, FileProp{objectId, ""}, renameRandFileName)

		So(err, ShouldBeNil)
		So(objId, ShouldEqual, objectId)

		//try renaming the object using using the same [newFileName]
		objId, err = RenameFile(dev, sid, FileProp{objectId, ""}, renameRandFileName)
		So(err, ShouldBeNil)
		So(objId, ShouldEqual, objectId)
	})

	Convey("Rename an existing object | using fullPath | RenameFile", t, func() {
		// create a random directory
		// test the directory '/mtp-test-files/temp_dir/test-RenameFile/{random}'
		fileName := fmt.Sprintf("/mtp-test-files/temp_dir/test-RenameFile/%x", rand.Int31())
		renameRandFileName := fmt.Sprintf("renamed-%x", rand.Int31())

		objectId, err := MakeDirectory(dev, sid, fileName)

		So(err, ShouldBeNil)
		So(objectId, ShouldBeGreaterThan, 0)

		//rename the object using objectId
		objId, err := RenameFile(dev, sid, FileProp{0, fileName}, renameRandFileName)

		So(err, ShouldBeNil)
		So(objId, ShouldEqual, objectId)

		time.Sleep(10000)

		//try renaming the object using using the same [newFileName]
		objId, err = RenameFile(dev, sid, FileProp{0, getFullPath("/mtp-test-files/temp_dir/test-RenameFile/", renameRandFileName)}, renameRandFileName)

		So(err, ShouldBeNil)
		So(objId, ShouldEqual, objectId)
	})

	Convey("Rename an non existing object | using objectId | RenameFile | Should throw an error", t, func() {
		//rename the object using objectId
		objId, err := RenameFile(dev, sid, FileProp{1234567, ""}, "fake name")

		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objId, ShouldEqual, 0)
	})

	Convey("Rename an non existing object | using fullPath | RenameFile | Should throw an error", t, func() {
		// test the directory '/mtp-test-files/temp_dir/test-RenameFile/{random}'
		fileName := fmt.Sprintf("/mtp-test-files/temp_dir/test-RenameFile/%x", rand.Int31())
		objId, err := RenameFile(dev, sid, FileProp{0, fileName}, "fake name")

		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objId, ShouldEqual, 0)
	})

	Dispose(dev)
}

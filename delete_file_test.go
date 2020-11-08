package mtpx

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"math/rand"
	"testing"
)

func TestDeleteFile(t *testing.T) {
	dev, err := Initialize(Init{})
	if err != nil {
		log.Panic(err)
	}

	storages, err := FetchStorages(dev)
	if err != nil {
		log.Panic(err)
	}

	sid := storages[0].Sid

	Convey("Delete an existing object | using objectId | DeleteFile", t, func() {
		// create a random directory
		// test the directory '/mtp-test-files/temp_dir/test-DeleteFile/{random}'
		directoryName := fmt.Sprintf("/mtp-test-files/temp_dir/test-DeleteFile/%x", rand.Int31())

		objectId, err := MakeDirectory(dev, sid, directoryName)

		So(err, ShouldBeNil)
		So(objectId, ShouldBeGreaterThan, 0)

		//delete the object using objectId
		err = DeleteFile(dev, sid, []FileProp{{objectId, ""}})

		So(err, ShouldBeNil)
	})

	Convey("Delete multiple existing objects | using objectId | DeleteFile", t, func() {
		// create a random directory
		// test the directory '/mtp-test-files/temp_dir/test-DeleteFile/{random}'
		directoryName1 := fmt.Sprintf("/mtp-test-files/temp_dir/test-DeleteFile/%x", rand.Int31())
		directoryName2 := fmt.Sprintf("/mtp-test-files/temp_dir/test-DeleteFile/%x", rand.Int31())

		objectId1, err := MakeDirectory(dev, sid, directoryName1)
		So(err, ShouldBeNil)

		objectId2, err := MakeDirectory(dev, sid, directoryName2)
		So(err, ShouldBeNil)

		So(objectId1, ShouldBeGreaterThan, 0)

		//delete the object using objectId
		err = DeleteFile(dev, sid, []FileProp{
			{objectId1, ""},
			{objectId2, ""},
		})
		So(err, ShouldBeNil)
	})

	Convey("Delete an existing object | using fullPath | DeleteFile", t, func() {
		// create a random directory
		// test the directory '/mtp-test-files/temp_dir/test-DeleteFile/{random}'
		directoryName := fmt.Sprintf("/mtp-test-files/temp_dir/test-DeleteFile/%x", rand.Int31())

		objectId, err := MakeDirectory(dev, sid, directoryName)

		So(err, ShouldBeNil)
		So(objectId, ShouldBeGreaterThan, 0)

		//delete the object using objectId
		err = DeleteFile(dev, sid, []FileProp{{0, directoryName}})

		So(err, ShouldBeNil)
	})

	Convey("Delete multiple existing objects | using fullPath | DeleteFile", t, func() {
		// create a random directory
		// test the directory '/mtp-test-files/temp_dir/test-DeleteFile/{random}'
		directoryName1 := fmt.Sprintf("/mtp-test-files/temp_dir/test-DeleteFile/%x", rand.Int31())
		directoryName2 := fmt.Sprintf("/mtp-test-files/temp_dir/test-DeleteFile/%x", rand.Int31())

		objectId1, err := MakeDirectory(dev, sid, directoryName1)
		So(err, ShouldBeNil)

		objectId2, err := MakeDirectory(dev, sid, directoryName2)
		So(err, ShouldBeNil)

		So(objectId1, ShouldBeGreaterThan, 0)
		So(objectId2, ShouldBeGreaterThan, 0)

		//delete the object using objectId
		err = DeleteFile(dev, sid, []FileProp{
			{0, directoryName1},
			{0, directoryName2},
		})

		So(err, ShouldBeNil)
	})

	Convey("Delete an non existing object | using objectId | DeleteFile", t, func() {
		//delete the object using objectId
		err = DeleteFile(dev, sid, []FileProp{{1234567, ""}})

		So(err, ShouldBeNil)
	})

	Convey("Delete an non existing object | using fullPath | DeleteFile", t, func() {
		// test the directory '/mtp-test-files/temp_dir/test-DeleteFile/{random}'
		directoryName := fmt.Sprintf("/mtp-test-files/temp_dir/test-DeleteFile/%x", rand.Int31())

		//delete the object using objectId
		err = DeleteFile(dev, sid, []FileProp{{0, directoryName}})

		So(err, ShouldBeNil)
	})

	Dispose(dev)
}

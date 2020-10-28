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

	sid := storages[0].sid

	Convey("Delete an existing object | using objectId | DeleteFile", t, func() {
		// create a random directory
		// test the directory '/mtp-test-files/temp_dir/test-DeleteFile/{random}'
		directoryName := fmt.Sprintf("/mtp-test-files/temp_dir/test-DeleteFile/%x", rand.Int31())

		objectId, err := MakeDirectory(dev, sid, directoryName)

		So(err, ShouldBeNil)
		So(objectId, ShouldBeGreaterThan, 0)

		//delete the object using objectId
		err = DeleteFile(dev, sid, objectId, "")

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
		err = DeleteFile(dev, sid, 0, directoryName)

		So(err, ShouldBeNil)
	})

	Convey("Delete an non existing object | using objectId | DeleteFile", t, func() {
		//delete the object using objectId
		err = DeleteFile(dev, sid, 1234567, "")

		So(err, ShouldBeNil)
	})

	Convey("Delete an non existing object | using fullPath | DeleteFile", t, func() {
		// test the directory '/mtp-test-files/temp_dir/test-DeleteFile/{random}'
		directoryName := fmt.Sprintf("/mtp-test-files/temp_dir/test-DeleteFile/%x", rand.Int31())

		//delete the object using objectId
		err = DeleteFile(dev, sid, 0, directoryName)

		So(err, ShouldBeNil)
	})

	Dispose(dev)
}

package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"testing"
)

func TestUploadFiles(t *testing.T) {
	dev, err := Initialize(Init{})
	if err != nil {
		log.Panic(err)
	}

	storages, err := FetchStorages(dev)
	if err != nil {
		log.Panic(err)
	}

	sid := storages[0].sid

	Convey("General | UploadFiles", t, func() {
		uploadFile1 := getTestMocksAsset("mock_dir1")
		sources := []string{uploadFile1}
		destination := "/mtp-test-files/temp_dir/test_UploadFiles"

		var currentSentTime int64
		var currentSentFiles int
		objectIdDest, totalFiles, totalSize, err := UploadFiles(dev, sid,
			sources,
			destination,
			func(fi *UploadFileInfo) {
				So(fi, ShouldNotBeNil)
				So(fi.StartTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)
				So(fi.LatestSentTime.UnixNano(), ShouldBeGreaterThan, currentSentTime)
				So(fi.Speed, ShouldBeGreaterThanOrEqualTo, 0)
				So(fi.FilesSent, ShouldBeGreaterThanOrEqualTo, 0)

				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)
				So(fi.FileInfo.FullPath, ShouldStartWith, destination)

				currentSentTime = fi.LatestSentTime.UnixNano()
				currentSentFiles = fi.FilesSent
			},
		)

		So(err, ShouldBeNil)
		So(currentSentFiles, ShouldEqual, 5)
		So(totalFiles, ShouldEqual, 5)
		So(totalFiles, ShouldEqual, currentSentFiles)
		So(totalSize, ShouldEqual, 35)

		fi, err := GetObjectFromPath(dev, sid, destination)
		So(err, ShouldBeNil)

		So(objectIdDest, ShouldEqual, fi.ObjectId)
	})

	Dispose(dev)
}

package mtpx

import (
	"fmt"
	"math/rand"
	"strings"

	//"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	//"time"

	//"math/rand"
	"os"
	//"strings"
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

	sid := storages[0].Sid

	Convey("General | UploadFiles", t, func() {
		// destination directories: '/mtp-test-files/temp_dir/test_UploadFiles'
		// source directories: 'mock_dir1'
		uploadFile1 := getTestMocksAsset("mock_dir1")
		sources := []string{uploadFile1}
		destination := "/mtp-test-files/temp_dir/test_UploadFiles"

		var prevLatestSentTime int64
		var prevFilesSent int64
		var prevFilesSentProgress float32
		var prevCurrentSentProgress float32
		var prevBulkSentProgress float32
		var prevBulkSent int64
		objectIdDest, totalFiles, totalSize, err := UploadFiles(dev, sid,
			sources,
			destination,
			false,
			func(fi *os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				// this function should not be called
				count := 0
				So(count, ShouldNotEqual, count)

				return nil
			},
			func(fi *ProgressInfo, err error) error {
				So(err, ShouldBeNil)
				So(fi, ShouldNotBeNil)

				So(fi.StartTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)
				So(fi.LatestSentTime.UnixNano(), ShouldBeGreaterThan, prevLatestSentTime)
				prevLatestSentTime = fi.LatestSentTime.UnixNano()

				So(fi.Speed, ShouldEqual, 0)
				So(fi.FilesSent, ShouldEqual, prevFilesSent)

				if fi.Status == InProgress {
					prevFilesSent += 1
				}

				So(fi.TotalDirectories, ShouldEqual, 0)

				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)
				So(fi.FileInfo.FullPath, ShouldStartWith, destination)

				So(fi.FilesSentProgress, ShouldBeGreaterThanOrEqualTo, prevFilesSentProgress)
				prevFilesSentProgress = fi.FilesSentProgress

				// current progress tests
				So(fi.Current.Total, ShouldBeGreaterThan, 0)
				So(fi.Current.Sent, ShouldEqual, fi.Current.Total)

				So(fi.Current.Progress, ShouldBeGreaterThanOrEqualTo, prevCurrentSentProgress)
				prevCurrentSentProgress = fi.Current.Progress

				// bulk progress tests
				So(fi.Bulk.Total, ShouldEqual, 0)
				So(fi.Bulk.Sent, ShouldBeGreaterThanOrEqualTo, prevBulkSent)
				prevBulkSent = fi.Bulk.Sent

				So(fi.Bulk.Progress, ShouldBeGreaterThanOrEqualTo, prevBulkSentProgress)
				prevBulkSentProgress = fi.Bulk.Progress

				return nil
			},
		)

		So(err, ShouldBeNil)
		So(prevFilesSent, ShouldEqual, 5)
		So(totalFiles, ShouldEqual, 5)
		So(totalFiles, ShouldEqual, prevFilesSent)
		So(totalSize, ShouldEqual, 35)

		fi, err := GetObjectFromPath(dev, sid, destination)
		So(err, ShouldBeNil)

		So(objectIdDest, ShouldEqual, fi.ObjectId)
	})

	Convey("Single directory | Random destination | UploadFiles", t, func() {
		// destination directories: '/mtp-test-files/temp_dir/test_UploadFiles'
		// source directories: 'mock_dir1'
		uploadFile1 := getTestMocksAsset("mock_dir1")
		sources := []string{uploadFile1}

		randFName := fmt.Sprintf("%x", rand.Int31())
		destination := getFullPath("/mtp-test-files/temp_dir/test_UploadFiles", randFName)

		dirList := []string{"/mock_dir1/1/a.txt", "/mock_dir1/2/b.txt", "/mock_dir1/3/2/b.txt", "/mock_dir1/3/b.txt", "/mock_dir1/a.txt"}

		var prevLatestSentTime int64
		var prevFilesSent int64
		var prevFilesSentProgress float32
		var prevCurrentSentProgress float32
		var prevBulkSentProgress float32
		var prevBulkSent int64
		var status TransferStatus
		objectIdDest, totalFiles, totalSize, err := UploadFiles(dev, sid,
			sources,
			destination,
			false,
			func(fi *os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				// this function should not be called
				count := 0
				So(count, ShouldNotEqual, count)

				return nil
			},
			func(fi *ProgressInfo, err error) error {
				So(err, ShouldBeNil)
				So(fi, ShouldNotBeNil)

				So(fi.StartTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)
				So(fi.LatestSentTime.UnixNano(), ShouldBeGreaterThan, prevLatestSentTime)
				prevLatestSentTime = fi.LatestSentTime.UnixNano()

				So(fi.Speed, ShouldEqual, 0)
				So(fi.FilesSent, ShouldEqual, prevFilesSent)
				if fi.Status == InProgress {
					prevFilesSent += 1
				}

				So(fi.TotalDirectories, ShouldEqual, 0)

				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)
				So(fi.FileInfo.FullPath, ShouldStartWith, destination)

				if fi.Status == InProgress {
					So(fi.FileInfo.FullPath, ShouldEndWith, dirList[fi.FilesSent])
				}

				So(fi.FilesSentProgress, ShouldBeGreaterThanOrEqualTo, prevFilesSentProgress)
				prevFilesSentProgress = fi.FilesSentProgress

				// current progress tests
				So(fi.Current.Total, ShouldBeGreaterThan, 0)
				So(fi.Current.Sent, ShouldEqual, fi.Current.Total)

				So(fi.Current.Progress, ShouldBeGreaterThanOrEqualTo, prevCurrentSentProgress)
				prevCurrentSentProgress = fi.Current.Progress

				// bulk progress tests
				So(fi.Bulk.Total, ShouldEqual, 0)
				So(fi.Bulk.Sent, ShouldBeGreaterThanOrEqualTo, prevBulkSent)
				prevBulkSent = fi.Bulk.Sent

				So(fi.Bulk.Progress, ShouldBeGreaterThanOrEqualTo, prevBulkSentProgress)
				prevBulkSentProgress = fi.Bulk.Progress

				status = fi.Status
				return nil
			},
		)

		So(err, ShouldBeNil)
		So(prevFilesSent, ShouldEqual, 5)
		So(totalFiles, ShouldEqual, 5)
		So(totalFiles, ShouldEqual, prevFilesSent)
		So(totalSize, ShouldEqual, 35)
		So(status, ShouldEqual, Completed)

		fi, err := GetObjectFromPath(dev, sid, destination)
		So(err, ShouldBeNil)

		So(objectIdDest, ShouldEqual, fi.ObjectId)

		//walk the directory on device and verify
		dirList1 := []string{
			"/mock_dir1",
			"/mock_dir1/1",
			"/mock_dir1/1/a.txt",
			"/mock_dir1/2",
			"/mock_dir1/2/b.txt",
			"/mock_dir1/3",
			"/mock_dir1/3/2",
			"/mock_dir1/3/2/b.txt",
			"/mock_dir1/3/b.txt",
			"/mock_dir1/a.txt"}

		objectId, totalListFiles, err := Walk(dev, sid, destination, true, true, func(objectId uint32, fi *FileInfo, err error) error {
			So(err, ShouldBeNil)

			contains, index := StringContains(dirList1, strings.TrimPrefix(fi.FullPath, destination))
			So(contains, ShouldEqual, true)
			dirList1 = RemoveIndex(dirList1, index)

			return nil
		})

		So(err, ShouldBeNil)
		So(objectIdDest, ShouldEqual, objectId)
		So(totalListFiles, ShouldEqual, 10)
	})

	Convey("Multiple directories | Random destination | UploadFiles", t, func() {
		// destination directories: '/mtp-test-files/temp_dir/test_UploadFiles/{random}'
		// source directories: 'mock_dir1' and 'mock_dir2'
		uploadFile1 := getTestMocksAsset("mock_dir1")
		uploadFile2 := getTestMocksAsset("mock_dir2")
		sources := []string{uploadFile1, uploadFile2}

		randFName := fmt.Sprintf("%x", rand.Int31())
		destination := getFullPath("/mtp-test-files/temp_dir/test_UploadFiles", randFName)

		dirList := []string{
			"/mock_dir1/1/a.txt", "/mock_dir1/2/b.txt", "/mock_dir1/3/2/b.txt", "/mock_dir1/3/b.txt", "/mock_dir1/a.txt",
			"/mock_dir2/1/a.txt", "/mock_dir2/2/b.txt", "/mock_dir2/3/2/b.txt", "/mock_dir2/3/b.txt", "/mock_dir2/a.txt",
		}

		var prevLatestSentTime int64
		var prevFilesSent int64
		var prevFilesSentProgress float32
		var prevCurrentSentProgress float32
		var prevBulkSentProgress float32
		var prevBulkSent int64
		var status TransferStatus
		objectIdDest, totalFiles, totalSize, err := UploadFiles(dev, sid,
			sources,
			destination,
			false,
			func(fi *os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				// this function should not be called
				count := 0
				So(count, ShouldNotEqual, count)

				return nil
			},
			func(fi *ProgressInfo, err error) error {
				So(err, ShouldBeNil)
				So(fi, ShouldNotBeNil)

				So(fi.StartTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)
				So(fi.LatestSentTime.UnixNano(), ShouldBeGreaterThan, prevLatestSentTime)
				prevLatestSentTime = fi.LatestSentTime.UnixNano()

				So(fi.Speed, ShouldEqual, 0)
				So(fi.FilesSent, ShouldEqual, prevFilesSent)
				if fi.Status == InProgress {
					prevFilesSent += 1
				}

				So(fi.TotalDirectories, ShouldEqual, 0)

				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)
				So(fi.FileInfo.FullPath, ShouldStartWith, destination)

				if fi.Status == InProgress {
					So(fi.FileInfo.FullPath, ShouldEndWith, dirList[fi.FilesSent])
				}

				So(fi.FilesSentProgress, ShouldBeGreaterThanOrEqualTo, prevFilesSentProgress)
				prevFilesSentProgress = fi.FilesSentProgress

				// current progress tests
				So(fi.Current.Total, ShouldBeGreaterThan, 0)
				So(fi.Current.Sent, ShouldEqual, fi.Current.Total)

				So(fi.Current.Progress, ShouldBeGreaterThanOrEqualTo, prevCurrentSentProgress)
				prevCurrentSentProgress = fi.Current.Progress

				// bulk progress tests
				So(fi.Bulk.Total, ShouldEqual, 0)
				So(fi.Bulk.Sent, ShouldBeGreaterThanOrEqualTo, prevBulkSent)
				prevBulkSent = fi.Bulk.Sent

				So(fi.Bulk.Progress, ShouldBeGreaterThanOrEqualTo, prevBulkSentProgress)
				prevBulkSentProgress = fi.Bulk.Progress

				status = fi.Status
				return nil
			},
		)
		So(err, ShouldBeNil)
		So(err, ShouldBeNil)
		So(prevFilesSent, ShouldEqual, 5*2)
		So(totalFiles, ShouldEqual, 5*2)
		So(totalFiles, ShouldEqual, prevFilesSent)
		So(totalSize, ShouldEqual, 35*2)
		So(status, ShouldEqual, Completed)

		////walk the directory on device and verify
		dirList1 := []string{
			"/mock_dir1",
			"/mock_dir1/1",
			"/mock_dir1/1/a.txt",
			"/mock_dir1/2",
			"/mock_dir1/2/b.txt",
			"/mock_dir1/3",
			"/mock_dir1/3/2",
			"/mock_dir1/3/2/b.txt",
			"/mock_dir1/3/b.txt",
			"/mock_dir1/a.txt",
			"/mock_dir2",
			"/mock_dir2/1",
			"/mock_dir2/1/a.txt",
			"/mock_dir2/2",
			"/mock_dir2/2/b.txt",
			"/mock_dir2/3",
			"/mock_dir2/3/2",
			"/mock_dir2/3/2/b.txt",
			"/mock_dir2/3/b.txt",
			"/mock_dir2/a.txt",
		}

		objectId, totalListFiles, err := Walk(dev, sid, destination, true, true, func(objectId uint32, fi *FileInfo, err error) error {
			So(err, ShouldBeNil)

			contains, index := StringContains(dirList1, strings.TrimPrefix(fi.FullPath, destination))
			So(contains, ShouldEqual, true)
			dirList1 = RemoveIndex(dirList1, index)

			return nil
		})

		fi, err := GetObjectFromPath(dev, sid, destination)
		So(err, ShouldBeNil)

		So(objectIdDest, ShouldEqual, fi.ObjectId)
		So(objectIdDest, ShouldEqual, objectId)
		So(totalListFiles, ShouldEqual, 10*2)
	})

	Convey("Multiple directories | same name | Random destination | UploadFiles", t, func() {
		// destination directories: '/mtp-test-files/temp_dir/test_UploadFiles/{random}'
		// source directories: 'mock_dir1'
		// source directories: 'mock_dir1'
		uploadFile1 := getTestMocksAsset("mock_dir1")
		uploadFile2 := getTestMocksAsset("mock_dir1")
		sources := []string{uploadFile1, uploadFile2}

		randFName := fmt.Sprintf("%x", rand.Int31())
		destination := getFullPath("/mtp-test-files/temp_dir/test_UploadFiles", randFName)

		dirList := []string{
			"/mock_dir1/1/a.txt", "/mock_dir1/2/b.txt", "/mock_dir1/3/2/b.txt", "/mock_dir1/3/b.txt", "/mock_dir1/a.txt",
			"/mock_dir1/1/a.txt", "/mock_dir1/2/b.txt", "/mock_dir1/3/2/b.txt", "/mock_dir1/3/b.txt", "/mock_dir1/a.txt",
		}

		var prevLatestSentTime int64
		var prevFilesSent int64
		var prevFilesSentProgress float32
		var prevCurrentSentProgress float32
		var prevBulkSentProgress float32
		var prevBulkSent int64
		var status TransferStatus
		objectIdDest, totalFiles, totalSize, err := UploadFiles(dev, sid,
			sources,
			destination,
			false,
			func(fi *os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				// this function should not be called
				count := 0
				So(count, ShouldNotEqual, count)

				return nil
			},
			func(fi *ProgressInfo, err error) error {
				So(err, ShouldBeNil)
				So(fi, ShouldNotBeNil)

				So(fi.StartTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)
				So(fi.LatestSentTime.UnixNano(), ShouldBeGreaterThan, prevLatestSentTime)
				prevLatestSentTime = fi.LatestSentTime.UnixNano()

				So(fi.Speed, ShouldEqual, 0)
				So(fi.FilesSent, ShouldEqual, prevFilesSent)
				if fi.Status == InProgress {
					prevFilesSent += 1
				}

				So(fi.TotalDirectories, ShouldEqual, 0)

				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)
				So(fi.FileInfo.FullPath, ShouldStartWith, destination)

				if fi.Status == InProgress {
					So(fi.FileInfo.FullPath, ShouldEndWith, dirList[fi.FilesSent])
				}

				So(fi.FilesSentProgress, ShouldBeGreaterThanOrEqualTo, prevFilesSentProgress)
				prevFilesSentProgress = fi.FilesSentProgress

				// current progress tests
				So(fi.Current.Total, ShouldBeGreaterThan, 0)
				So(fi.Current.Sent, ShouldEqual, fi.Current.Total)

				So(fi.Current.Progress, ShouldBeGreaterThanOrEqualTo, prevCurrentSentProgress)
				prevCurrentSentProgress = fi.Current.Progress

				// bulk progress tests
				So(fi.Bulk.Total, ShouldEqual, 0)
				So(fi.Bulk.Sent, ShouldBeGreaterThanOrEqualTo, prevBulkSent)
				prevBulkSent = fi.Bulk.Sent

				So(fi.Bulk.Progress, ShouldBeGreaterThanOrEqualTo, prevBulkSentProgress)
				prevBulkSentProgress = fi.Bulk.Progress

				status = fi.Status
				return nil
			},
		)

		So(err, ShouldBeNil)
		So(prevFilesSent, ShouldEqual, 5*2)
		So(totalFiles, ShouldEqual, 5*2)
		So(totalSize, ShouldEqual, 35*2)
		So(status, ShouldEqual, Completed)

		////walk the directory on device and verify
		dirList1 := []string{
			"/mock_dir1",
			"/mock_dir1/1",
			"/mock_dir1/1/a.txt",
			"/mock_dir1/2",
			"/mock_dir1/2/b.txt",
			"/mock_dir1/3",
			"/mock_dir1/3/2",
			"/mock_dir1/3/2/b.txt",
			"/mock_dir1/3/b.txt",
			"/mock_dir1/a.txt",
		}

		objectId, totalListFiles, err := Walk(dev, sid, destination, true, true, func(objectId uint32, fi *FileInfo, err error) error {
			So(err, ShouldBeNil)

			contains, index := StringContains(dirList1, strings.TrimPrefix(fi.FullPath, destination))
			So(contains, ShouldEqual, true)
			dirList1 = RemoveIndex(dirList1, index)

			return nil
		})
		fi, err := GetObjectFromPath(dev, sid, destination)
		So(err, ShouldBeNil)

		So(objectIdDest, ShouldEqual, fi.ObjectId)
		So(objectIdDest, ShouldEqual, objectId)
		So(totalListFiles, ShouldEqual, 10*1)
	})

	Convey("Single File | Random destination | UploadFiles", t, func() {
		// test the directories: '/mtp-test-files/temp_dir/test_UploadFiles/{random}'
		// source file: 'a.txt'
		uploadFile1 := getTestMocksAsset("mock_dir1/a.txt")
		sources := []string{uploadFile1}

		randFName := fmt.Sprintf("%x", rand.Int31())
		destination := getFullPath("/mtp-test-files/temp_dir/test_UploadFiles", randFName)

		dirList := []string{"/a.txt"}

		var prevLatestSentTime int64
		var prevFilesSent int64
		var prevFilesSentProgress float32
		var prevCurrentSentProgress float32
		var prevBulkSentProgress float32
		var prevBulkSent int64
		var status TransferStatus
		objectIdDest, totalFiles, totalSize, err := UploadFiles(dev, sid,
			sources,
			destination,
			false,
			func(fi *os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				// this function should not be called
				count := 0
				So(count, ShouldNotEqual, count)

				return nil
			},
			func(fi *ProgressInfo, err error) error {
				So(err, ShouldBeNil)
				So(fi, ShouldNotBeNil)

				So(fi.StartTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)
				So(fi.LatestSentTime.UnixNano(), ShouldBeGreaterThan, prevLatestSentTime)
				prevLatestSentTime = fi.LatestSentTime.UnixNano()

				So(fi.Speed, ShouldEqual, 0)
				So(fi.FilesSent, ShouldEqual, prevFilesSent)
				if fi.Status == InProgress {
					prevFilesSent += 1
				}

				So(fi.TotalDirectories, ShouldEqual, 0)

				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)
				So(fi.FileInfo.FullPath, ShouldStartWith, destination)

				if fi.Status == InProgress {
					So(fi.FileInfo.FullPath, ShouldEndWith, dirList[fi.FilesSent])
				}

				So(fi.FilesSentProgress, ShouldBeGreaterThanOrEqualTo, prevFilesSentProgress)
				prevFilesSentProgress = fi.FilesSentProgress

				// current progress tests
				So(fi.Current.Total, ShouldBeGreaterThan, 0)
				So(fi.Current.Sent, ShouldEqual, fi.Current.Total)

				So(fi.Current.Progress, ShouldBeGreaterThanOrEqualTo, prevCurrentSentProgress)
				prevCurrentSentProgress = fi.Current.Progress

				// bulk progress tests
				So(fi.Bulk.Total, ShouldEqual, 0)
				So(fi.Bulk.Sent, ShouldBeGreaterThanOrEqualTo, prevBulkSent)
				prevBulkSent = fi.Bulk.Sent

				So(fi.Bulk.Progress, ShouldBeGreaterThanOrEqualTo, prevBulkSentProgress)
				prevBulkSentProgress = fi.Bulk.Progress

				status = fi.Status
				return nil
			},
		)

		So(err, ShouldBeNil)
		So(prevFilesSent, ShouldEqual, 1)
		So(totalFiles, ShouldEqual, 1)
		So(totalSize, ShouldEqual, 9)
		So(status, ShouldEqual, Completed)

		fi, err := GetObjectFromPath(dev, sid, destination)
		So(err, ShouldBeNil)

		So(objectIdDest, ShouldEqual, fi.ObjectId)
	})

	Convey("Multiple Files | Random destination | UploadFiles", t, func() {
		// test the directories: '/mtp-test-files/temp_dir/test_UploadFiles/{random}'
		// source file: 'a.txt'
		// source file: 'b.txt'

		uploadFile1 := getTestMocksAsset("mock_dir1/a.txt")
		uploadFile2 := getTestMocksAsset("mock_dir1/2/b.txt")
		sources := []string{uploadFile1, uploadFile2}

		randFName := fmt.Sprintf("%x", rand.Int31())
		destination := getFullPath("/mtp-test-files/temp_dir/test_UploadFiles", randFName)

		dirList := []string{"/a.txt", "/b.txt"}

		var prevLatestSentTime int64
		var prevFilesSent int64
		var prevFilesSentProgress float32
		var prevCurrentSentProgress float32
		var prevBulkSentProgress float32
		var prevBulkSent int64
		var status TransferStatus
		objectIdDest, totalFiles, totalSize, err := UploadFiles(dev, sid,
			sources,
			destination,
			false,
			func(fi *os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				// this function should not be called
				count := 0
				So(count, ShouldNotEqual, count)

				return nil
			},
			func(fi *ProgressInfo, err error) error {
				So(err, ShouldBeNil)
				So(fi, ShouldNotBeNil)

				So(fi.StartTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)
				So(fi.LatestSentTime.UnixNano(), ShouldBeGreaterThan, prevLatestSentTime)
				prevLatestSentTime = fi.LatestSentTime.UnixNano()

				So(fi.Speed, ShouldEqual, 0)
				So(fi.FilesSent, ShouldEqual, prevFilesSent)
				if fi.Status == InProgress {
					prevFilesSent += 1
				}

				So(fi.TotalDirectories, ShouldEqual, 0)

				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)
				So(fi.FileInfo.FullPath, ShouldStartWith, destination)

				if fi.Status == InProgress {
					So(fi.FileInfo.FullPath, ShouldEndWith, dirList[fi.FilesSent])
				}

				So(fi.FilesSentProgress, ShouldBeGreaterThanOrEqualTo, prevFilesSentProgress)
				prevFilesSentProgress = fi.FilesSentProgress

				// current progress tests
				So(fi.Current.Total, ShouldBeGreaterThan, 0)
				So(fi.Current.Sent, ShouldEqual, fi.Current.Total)

				So(fi.Current.Progress, ShouldBeGreaterThanOrEqualTo, prevCurrentSentProgress)
				prevCurrentSentProgress = fi.Current.Progress

				// bulk progress tests
				So(fi.Bulk.Total, ShouldEqual, 0)
				So(fi.Bulk.Sent, ShouldBeGreaterThanOrEqualTo, prevBulkSent)
				prevBulkSent = fi.Bulk.Sent

				So(fi.Bulk.Progress, ShouldBeGreaterThanOrEqualTo, prevBulkSentProgress)
				prevBulkSentProgress = fi.Bulk.Progress

				status = fi.Status
				return nil
			},
		)
		So(err, ShouldBeNil)
		So(prevFilesSent, ShouldEqual, 1*2)
		So(totalFiles, ShouldEqual, 1*2)
		So(totalSize, ShouldEqual, 15)
		So(status, ShouldEqual, Completed)

		fi, err := GetObjectFromPath(dev, sid, destination)
		So(err, ShouldBeNil)

		So(objectIdDest, ShouldEqual, fi.ObjectId)
	})

	Convey("Multiple Files | same name | Random destination | UploadFiles", t, func() {
		// test the directories: '/mtp-test-files/temp_dir/test_UploadFiles/{random}'
		// source "mock_dir1/a.txt"
		// source "mock_dir1/1/a.txt"

		uploadFile1 := getTestMocksAsset("mock_dir1/a.txt")
		uploadFile2 := getTestMocksAsset("mock_dir1/1/a.txt")
		sources := []string{uploadFile1, uploadFile2}

		randFName := fmt.Sprintf("%x", rand.Int31())
		destination := getFullPath("/mtp-test-files/temp_dir/test_UploadFiles", randFName)

		dirList := []string{"/a.txt", "/a.txt"}

		var prevLatestSentTime int64
		var prevFilesSent int64
		var prevFilesSentProgress float32
		var prevCurrentSentProgress float32
		var prevBulkSentProgress float32
		var prevBulkSent int64
		var status TransferStatus
		objectIdDest, totalFiles, totalSize, err := UploadFiles(dev, sid,
			sources,
			destination,
			false,
			func(fi *os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				// this function should not be called
				count := 0
				So(count, ShouldNotEqual, count)

				return nil
			},
			func(fi *ProgressInfo, err error) error {
				So(err, ShouldBeNil)
				So(fi, ShouldNotBeNil)

				So(fi.StartTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)
				So(fi.LatestSentTime.UnixNano(), ShouldBeGreaterThan, prevLatestSentTime)
				prevLatestSentTime = fi.LatestSentTime.UnixNano()

				So(fi.Speed, ShouldEqual, 0)
				So(fi.FilesSent, ShouldEqual, prevFilesSent)
				if fi.Status == InProgress {
					prevFilesSent += 1
				}

				So(fi.TotalDirectories, ShouldEqual, 0)

				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)
				So(fi.FileInfo.FullPath, ShouldStartWith, destination)

				if fi.Status == InProgress {
					So(fi.FileInfo.FullPath, ShouldEndWith, dirList[fi.FilesSent])
				}

				So(fi.FilesSentProgress, ShouldBeGreaterThanOrEqualTo, prevFilesSentProgress)
				prevFilesSentProgress = fi.FilesSentProgress

				// current progress tests
				So(fi.Current.Total, ShouldBeGreaterThan, 0)
				So(fi.Current.Sent, ShouldEqual, fi.Current.Total)

				So(fi.Current.Progress, ShouldBeGreaterThanOrEqualTo, prevCurrentSentProgress)
				prevCurrentSentProgress = fi.Current.Progress

				// bulk progress tests
				So(fi.Bulk.Total, ShouldEqual, 0)
				So(fi.Bulk.Sent, ShouldBeGreaterThanOrEqualTo, prevBulkSent)
				prevBulkSent = fi.Bulk.Sent

				So(fi.Bulk.Progress, ShouldBeGreaterThanOrEqualTo, prevBulkSentProgress)
				prevBulkSentProgress = fi.Bulk.Progress

				status = fi.Status

				return nil
			},
		)

		So(err, ShouldBeNil)
		So(prevFilesSent, ShouldEqual, 2)
		So(totalFiles, ShouldEqual, 2)
		So(totalFiles, ShouldEqual, prevFilesSent)
		So(totalSize, ShouldEqual, 17)
		So(status, ShouldEqual, Completed)

		fi, err := GetObjectFromPath(dev, sid, destination)
		So(err, ShouldBeNil)

		So(objectIdDest, ShouldEqual, fi.ObjectId)
	})

	var _destination string
	Convey("Directories and Files | Random destination | UploadFiles", t, func() {
		// test the directories: '/mtp-test-files/temp_dir/test_UploadFiles/{random}'
		// source "mock_dir1/a.txt"
		// source "mock_dir1/1/a.txt"

		uploadFile1 := getTestMocksAsset("mock_dir1/")
		uploadFile2 := getTestMocksAsset("mock_dir1/1/a.txt")
		sources := []string{uploadFile1, uploadFile2}

		randFName := fmt.Sprintf("%x", rand.Int31())
		_destination = getFullPath("/mtp-test-files/temp_dir/test_UploadFiles", randFName)

		dirList := []string{
			"/mock_dir1/1/a.txt", "/mock_dir1/2/b.txt", "/mock_dir1/3/2/b.txt", "/mock_dir1/3/b.txt", "/mock_dir1/a.txt", "a.txt",
		}

		var prevLatestSentTime int64
		var prevFilesSent int64
		var prevFilesSentProgress float32
		var prevCurrentSentProgress float32
		var prevBulkSentProgress float32
		var prevBulkSent int64
		var status TransferStatus
		objectIdDest, totalFiles, totalSize, err := UploadFiles(dev, sid,
			sources,
			_destination,
			false,
			func(fi *os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				// this function should not be called
				count := 0
				So(count, ShouldNotEqual, count)

				return nil
			},
			func(fi *ProgressInfo, err error) error {
				So(err, ShouldBeNil)
				So(fi, ShouldNotBeNil)

				So(fi.StartTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)
				So(fi.LatestSentTime.UnixNano(), ShouldBeGreaterThan, prevLatestSentTime)
				prevLatestSentTime = fi.LatestSentTime.UnixNano()

				So(fi.Speed, ShouldEqual, 0)
				So(fi.FilesSent, ShouldEqual, prevFilesSent)
				if fi.Status == InProgress {
					prevFilesSent += 1
				}

				So(fi.TotalDirectories, ShouldEqual, 0)

				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)
				So(fi.FileInfo.FullPath, ShouldStartWith, _destination)

				if fi.Status == InProgress {
					So(fi.FileInfo.FullPath, ShouldEndWith, dirList[fi.FilesSent])
				}

				So(fi.FilesSentProgress, ShouldBeGreaterThanOrEqualTo, prevFilesSentProgress)
				prevFilesSentProgress = fi.FilesSentProgress

				// current progress tests
				So(fi.Current.Total, ShouldBeGreaterThan, 0)
				So(fi.Current.Sent, ShouldEqual, fi.Current.Total)

				So(fi.Current.Progress, ShouldBeGreaterThanOrEqualTo, prevCurrentSentProgress)
				prevCurrentSentProgress = fi.Current.Progress

				// bulk progress tests
				So(fi.Bulk.Total, ShouldEqual, 0)
				So(fi.Bulk.Sent, ShouldBeGreaterThanOrEqualTo, prevBulkSent)
				prevBulkSent = fi.Bulk.Sent

				So(fi.Bulk.Progress, ShouldBeGreaterThanOrEqualTo, prevBulkSentProgress)
				prevBulkSentProgress = fi.Bulk.Progress

				status = fi.Status

				return nil
			},
		)

		So(err, ShouldBeNil)
		So(prevFilesSent, ShouldEqual, 6)
		So(totalFiles, ShouldEqual, 6)
		So(totalFiles, ShouldEqual, prevFilesSent)
		So(totalSize, ShouldEqual, 43)
		So(status, ShouldEqual, Completed)

		fi, err := GetObjectFromPath(dev, sid, _destination)
		So(err, ShouldBeNil)

		So(objectIdDest, ShouldEqual, fi.ObjectId)
	})

	Convey("Directories and Files | Previous destination | UploadFiles", t, func() {
		// test the directories: '/mtp-test-files/temp_dir/test_UploadFiles/{random}'
		// source "mock_dir1/a.txt"
		// source "mock_dir1/1/a.txt"

		uploadFile1 := getTestMocksAsset("mock_dir1/")
		uploadFile2 := getTestMocksAsset("mock_dir1/1/a.txt")
		sources := []string{uploadFile1, uploadFile2}

		dirList := []string{
			"/mock_dir1/1/a.txt", "/mock_dir1/2/b.txt", "/mock_dir1/3/2/b.txt", "/mock_dir1/3/b.txt", "/mock_dir1/a.txt", "a.txt",
		}

		var prevLatestSentTime int64
		var prevFilesSent int64
		var prevFilesSentProgress float32
		var prevCurrentSentProgress float32
		var prevBulkSentProgress float32
		var prevBulkSent int64
		var status TransferStatus
		objectIdDest, totalFiles, totalSize, err := UploadFiles(dev, sid,
			sources,
			_destination,
			false,
			func(fi *os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				// this function should not be called
				count := 0
				So(count, ShouldNotEqual, count)

				return nil
			},
			func(fi *ProgressInfo, err error) error {
				So(err, ShouldBeNil)
				So(fi, ShouldNotBeNil)

				So(fi.StartTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)
				So(fi.LatestSentTime.UnixNano(), ShouldBeGreaterThan, prevLatestSentTime)
				prevLatestSentTime = fi.LatestSentTime.UnixNano()

				So(fi.Speed, ShouldEqual, 0)
				So(fi.FilesSent, ShouldEqual, prevFilesSent)
				if fi.Status == InProgress {
					prevFilesSent += 1
				}

				So(fi.TotalDirectories, ShouldEqual, 0)

				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)
				So(fi.FileInfo.FullPath, ShouldStartWith, _destination)

				if fi.Status == InProgress {
					So(fi.FileInfo.FullPath, ShouldEndWith, dirList[fi.FilesSent])
				}

				So(fi.FilesSentProgress, ShouldBeGreaterThanOrEqualTo, prevFilesSentProgress)
				prevFilesSentProgress = fi.FilesSentProgress

				// current progress tests
				So(fi.Current.Total, ShouldBeGreaterThan, 0)
				So(fi.Current.Sent, ShouldEqual, fi.Current.Total)

				So(fi.Current.Progress, ShouldBeGreaterThanOrEqualTo, prevCurrentSentProgress)
				prevCurrentSentProgress = fi.Current.Progress

				// bulk progress tests
				So(fi.Bulk.Total, ShouldEqual, 0)
				So(fi.Bulk.Sent, ShouldBeGreaterThanOrEqualTo, prevBulkSent)
				prevBulkSent = fi.Bulk.Sent

				So(fi.Bulk.Progress, ShouldBeGreaterThanOrEqualTo, prevBulkSentProgress)
				prevBulkSentProgress = fi.Bulk.Progress

				status = fi.Status

				return nil
			},
		)

		So(err, ShouldBeNil)
		So(prevFilesSent, ShouldEqual, 6)
		So(totalFiles, ShouldEqual, 6)
		So(totalFiles, ShouldEqual, prevFilesSent)
		So(totalSize, ShouldEqual, 43)
		So(status, ShouldEqual, Completed)

		fi, err := GetObjectFromPath(dev, sid, _destination)
		So(err, ShouldBeNil)

		So(objectIdDest, ShouldEqual, fi.ObjectId)
	})

	//todo huge file, speed test (single huge file, multiple huge file, along with smaller files (single, multiple), along with directories  (single and multiple), mixed
	//todo preprocessing
	//todo preprocessing error

	//
	//	Convey("Invalid source | Random destination | UploadFiles | should throw an error ", t, func() {
	//		// test the directories: '/mtp-test-files/temp_dir/test_UploadFiles/{random}'
	//
	//		uploadFile1 := getTestMocksAsset("mock_dir1/")
	//		uploadFile2 := "fake/1/111a.txt"
	//		sources := []string{uploadFile1, uploadFile2}
	//
	//		randFName := fmt.Sprintf("%x", rand.Int31())
	//		_destination = getFullPath("/mtp-test-files/temp_dir/test_UploadFiles", randFName)
	//
	//		dirList := []string{
	//			"/mock_dir1/1/a.txt", "/mock_dir1/2/b.txt", "/mock_dir1/3/2/b.txt", "/mock_dir1/3/b.txt", "/mock_dir1/a.txt", "a.txt",
	//		}
	//
	//		var currentSentTime int64
	//		var currentSentFiles int
	//		objectIdDest, totalFiles, totalSize, err := UploadFiles(dev, sid,
	//			sources,
	//			_destination,
	//			func(fi *ProgressInfo, err error) error {
	//				So(err, ShouldBeNil)
	//				So(fi, ShouldNotBeNil)
	//				So(fi.StartTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)
	//				So(fi.LatestSentTime.UnixNano(), ShouldBeGreaterThan, currentSentTime)
	//				So(fi.Speed, ShouldBeGreaterThanOrEqualTo, 0)
	//				So(fi.FilesSent, ShouldBeGreaterThanOrEqualTo, 0)
	//
	//				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)
	//				So(fi.FileInfo.FullPath, ShouldStartWith, _destination)
	//
	//				So(fi.FileInfo.FullPath, ShouldEndWith, dirList[fi.FilesSent-1])
	//
	//				currentSentTime = fi.LatestSentTime.UnixNano()
	//				currentSentFiles = fi.FilesSent
	//
	//				return nil
	//			},
	//		)
	//
	//		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
	//		So(currentSentFiles, ShouldEqual, 5)
	//		So(totalFiles, ShouldEqual, 5)
	//		So(totalFiles, ShouldEqual, currentSentFiles)
	//		So(totalSize, ShouldEqual, 35)
	//
	//		fi, err := GetObjectFromPath(dev, sid, _destination)
	//		So(err, ShouldBeNil)
	//
	//		So(objectIdDest, ShouldEqual, fi.ObjectId)
	//	})
	//
	//	Convey("Callback return error | UploadFiles | should throw an error ", t, func() {
	//		// test the directories: '/mtp-test-files/temp_dir/test_UploadFiles/{random}'
	//
	//		uploadFile1 := getTestMocksAsset("mock_dir1/")
	//		sources := []string{uploadFile1}
	//
	//		randFName := fmt.Sprintf("%x", rand.Int31())
	//		_destination = getFullPath("/mtp-test-files/temp_dir/test_UploadFiles", randFName)
	//
	//		dirList := []string{
	//			"/mock_dir1/1/a.txt", "/mock_dir1/2/b.txt", "/mock_dir1/3/2/b.txt", "/mock_dir1/3/b.txt", "/mock_dir1/a.txt",
	//		}
	//
	//		var currentSentTime int64
	//		var currentSentFiles int
	//		objectIdDest, totalFiles, totalSize, err := UploadFiles(dev, sid,
	//			sources, _destination,
	//			func(fi *ProgressInfo, err error) error {
	//				So(err, ShouldBeNil)
	//				So(fi, ShouldNotBeNil)
	//				So(fi.StartTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)
	//				So(fi.LatestSentTime.UnixNano(), ShouldBeGreaterThan, currentSentTime)
	//				So(fi.Speed, ShouldBeGreaterThanOrEqualTo, 0)
	//				So(fi.FilesSent, ShouldBeGreaterThanOrEqualTo, 0)
	//
	//				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)
	//				So(fi.FileInfo.FullPath, ShouldStartWith, _destination)
	//
	//				So(fi.FileInfo.FullPath, ShouldEndWith, dirList[fi.FilesSent-1])
	//
	//				currentSentTime = fi.LatestSentTime.UnixNano()
	//				currentSentFiles = fi.FilesSent
	//
	//				return FilePermissionError{error: fmt.Errorf("some random error")}
	//			},
	//		)
	//
	//		So(err, ShouldHaveSameTypeAs, FileTransferError{})
	//		So(currentSentFiles, ShouldEqual, 1)
	//		So(totalFiles, ShouldEqual, 1)
	//		So(totalFiles, ShouldEqual, currentSentFiles)
	//		So(totalSize, ShouldEqual, 8)
	//
	//		fi, err := GetObjectFromPath(dev, sid, _destination)
	//		So(err, ShouldBeNil)
	//
	//		So(objectIdDest, ShouldEqual, fi.ObjectId)
	//	})
	//
	Dispose(dev)
}

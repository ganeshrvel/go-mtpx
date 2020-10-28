package mtpx

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"math/rand"
	"strings"
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
		// test the directories: '/mtp-test-files/temp_dir/test_UploadFiles'

		uploadFile1 := getTestMocksAsset("mock_dir1")
		sources := []string{uploadFile1}
		destination := "/mtp-test-files/temp_dir/test_UploadFiles"

		var currentSentTime int64
		var currentSentFiles int
		objectIdDest, totalFiles, totalSize, err := UploadFiles(dev, sid,
			sources,
			destination,
			func(fi *TransferredFileInfo, err error) error {
				So(err, ShouldBeNil)
				So(fi, ShouldNotBeNil)
				So(fi.StartTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)
				So(fi.LatestSentTime.UnixNano(), ShouldBeGreaterThan, currentSentTime)
				So(fi.Speed, ShouldBeGreaterThanOrEqualTo, 0)
				So(fi.FilesSent, ShouldBeGreaterThanOrEqualTo, 0)

				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)
				So(fi.FileInfo.FullPath, ShouldStartWith, destination)

				currentSentTime = fi.LatestSentTime.UnixNano()
				currentSentFiles = fi.FilesSent

				return nil
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

	Convey("Single directory | Random destination | UploadFiles", t, func() {
		// test the directories: '/mtp-test-files/temp_dir/test_UploadFiles/{random}'

		uploadFile1 := getTestMocksAsset("mock_dir1")
		sources := []string{uploadFile1}

		randFName := fmt.Sprintf("%x", rand.Int31())
		destination := getFullPath("/mtp-test-files/temp_dir/test_UploadFiles", randFName)

		dirList := []string{"/mock_dir1/1/a.txt", "/mock_dir1/2/b.txt", "/mock_dir1/3/2/b.txt", "/mock_dir1/3/b.txt", "/mock_dir1/a.txt"}

		var currentSentTime int64
		var currentSentFiles int
		objectIdDest, totalFiles, totalSize, err := UploadFiles(dev, sid,
			sources,
			destination,
			func(fi *TransferredFileInfo, err error) error {
				So(err, ShouldBeNil)
				So(fi, ShouldNotBeNil)
				So(fi.StartTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)
				So(fi.LatestSentTime.UnixNano(), ShouldBeGreaterThan, currentSentTime)
				So(fi.Speed, ShouldBeGreaterThanOrEqualTo, 0)
				So(fi.FilesSent, ShouldBeGreaterThanOrEqualTo, 0)

				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)
				So(fi.FileInfo.FullPath, ShouldStartWith, destination)

				So(fi.FileInfo.FullPath, ShouldEndWith, dirList[fi.FilesSent-1])

				currentSentTime = fi.LatestSentTime.UnixNano()
				currentSentFiles = fi.FilesSent

				return nil
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
		// test the directories: '/mtp-test-files/temp_dir/test_UploadFiles/{random}'

		uploadFile1 := getTestMocksAsset("mock_dir1")
		uploadFile2 := getTestMocksAsset("mock_dir2")
		sources := []string{uploadFile1, uploadFile2}

		randFName := fmt.Sprintf("%x", rand.Int31())
		destination := getFullPath("/mtp-test-files/temp_dir/test_UploadFiles", randFName)

		dirList := []string{
			"/mock_dir1/1/a.txt", "/mock_dir1/2/b.txt", "/mock_dir1/3/2/b.txt", "/mock_dir1/3/b.txt", "/mock_dir1/a.txt",
			"/mock_dir2/1/a.txt", "/mock_dir2/2/b.txt", "/mock_dir2/3/2/b.txt", "/mock_dir2/3/b.txt", "/mock_dir2/a.txt",
		}

		var currentSentTime int64
		var currentSentFiles int
		objectIdDest, totalFiles, totalSize, err := UploadFiles(dev, sid,
			sources,
			destination,
			func(fi *TransferredFileInfo, err error) error {
				So(err, ShouldBeNil)
				So(fi, ShouldNotBeNil)
				So(fi.StartTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)
				So(fi.LatestSentTime.UnixNano(), ShouldBeGreaterThan, currentSentTime)
				So(fi.Speed, ShouldBeGreaterThanOrEqualTo, 0)
				So(fi.FilesSent, ShouldBeGreaterThanOrEqualTo, 0)

				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)
				So(fi.FileInfo.FullPath, ShouldStartWith, destination)

				So(fi.FileInfo.FullPath, ShouldEndWith, dirList[fi.FilesSent-1])

				currentSentTime = fi.LatestSentTime.UnixNano()
				currentSentFiles = fi.FilesSent

				return nil
			},
		)

		So(err, ShouldBeNil)
		So(currentSentFiles, ShouldEqual, 10)
		So(totalFiles, ShouldEqual, 10)
		So(totalFiles, ShouldEqual, currentSentFiles)
		So(totalSize, ShouldEqual, 70)

		fi, err := GetObjectFromPath(dev, sid, destination)
		So(err, ShouldBeNil)

		So(objectIdDest, ShouldEqual, fi.ObjectId)
	})

	Convey("Multiple directories | same name | Random destination | UploadFiles", t, func() {
		// test the directories: '/mtp-test-files/temp_dir/test_UploadFiles/{random}'

		uploadFile1 := getTestMocksAsset("mock_dir1")
		uploadFile2 := getTestMocksAsset("mock_dir1")
		sources := []string{uploadFile1, uploadFile2}

		randFName := fmt.Sprintf("%x", rand.Int31())
		destination := getFullPath("/mtp-test-files/temp_dir/test_UploadFiles", randFName)

		dirList := []string{
			"/mock_dir1/1/a.txt", "/mock_dir1/2/b.txt", "/mock_dir1/3/2/b.txt", "/mock_dir1/3/b.txt", "/mock_dir1/a.txt",
			"/mock_dir1/1/a.txt", "/mock_dir1/2/b.txt", "/mock_dir1/3/2/b.txt", "/mock_dir1/3/b.txt", "/mock_dir1/a.txt",
		}

		var currentSentTime int64
		var currentSentFiles int
		objectIdDest, totalFiles, totalSize, err := UploadFiles(dev, sid,
			sources,
			destination,
			func(fi *TransferredFileInfo, err error) error {
				So(err, ShouldBeNil)
				So(fi, ShouldNotBeNil)
				So(fi.StartTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)
				So(fi.LatestSentTime.UnixNano(), ShouldBeGreaterThan, currentSentTime)
				So(fi.Speed, ShouldBeGreaterThanOrEqualTo, 0)
				So(fi.FilesSent, ShouldBeGreaterThanOrEqualTo, 0)

				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)
				So(fi.FileInfo.FullPath, ShouldStartWith, destination)

				So(fi.FileInfo.FullPath, ShouldEndWith, dirList[fi.FilesSent-1])

				currentSentTime = fi.LatestSentTime.UnixNano()
				currentSentFiles = fi.FilesSent

				return nil
			},
		)

		So(err, ShouldBeNil)
		So(currentSentFiles, ShouldEqual, 10)
		So(totalFiles, ShouldEqual, 10)
		So(totalFiles, ShouldEqual, currentSentFiles)
		So(totalSize, ShouldEqual, 70)

		fi, err := GetObjectFromPath(dev, sid, destination)
		So(err, ShouldBeNil)

		So(objectIdDest, ShouldEqual, fi.ObjectId)
	})

	Convey("Single File | Random destination | UploadFiles", t, func() {
		// test the directories: '/mtp-test-files/temp_dir/test_UploadFiles/{random}'

		uploadFile1 := getTestMocksAsset("mock_dir1/a.txt")
		sources := []string{uploadFile1}

		randFName := fmt.Sprintf("%x", rand.Int31())
		destination := getFullPath("/mtp-test-files/temp_dir/test_UploadFiles", randFName)

		dirList := []string{"/a.txt"}

		var currentSentTime int64
		var currentSentFiles int
		objectIdDest, totalFiles, totalSize, err := UploadFiles(dev, sid,
			sources,
			destination,
			func(fi *TransferredFileInfo, err error) error {
				So(err, ShouldBeNil)
				So(fi, ShouldNotBeNil)
				So(fi.StartTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)
				So(fi.LatestSentTime.UnixNano(), ShouldBeGreaterThan, currentSentTime)
				So(fi.Speed, ShouldBeGreaterThanOrEqualTo, 0)
				So(fi.FilesSent, ShouldBeGreaterThanOrEqualTo, 0)

				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)
				So(fi.FileInfo.FullPath, ShouldStartWith, destination)

				So(fi.FileInfo.FullPath, ShouldEndWith, dirList[fi.FilesSent-1])

				currentSentTime = fi.LatestSentTime.UnixNano()
				currentSentFiles = fi.FilesSent

				return nil
			},
		)

		So(err, ShouldBeNil)
		So(currentSentFiles, ShouldEqual, 1)
		So(totalFiles, ShouldEqual, 1)
		So(totalFiles, ShouldEqual, currentSentFiles)
		So(totalSize, ShouldEqual, 9)

		fi, err := GetObjectFromPath(dev, sid, destination)
		So(err, ShouldBeNil)

		So(objectIdDest, ShouldEqual, fi.ObjectId)
	})

	Convey("Multiple Files | Random destination | UploadFiles", t, func() {
		// test the directories: '/mtp-test-files/temp_dir/test_UploadFiles/{random}'

		uploadFile1 := getTestMocksAsset("mock_dir1/a.txt")
		uploadFile2 := getTestMocksAsset("mock_dir1/2/b.txt")
		sources := []string{uploadFile1, uploadFile2}

		randFName := fmt.Sprintf("%x", rand.Int31())
		destination := getFullPath("/mtp-test-files/temp_dir/test_UploadFiles", randFName)

		dirList := []string{"/a.txt", "/b.txt"}

		var currentSentTime int64
		var currentSentFiles int
		objectIdDest, totalFiles, totalSize, err := UploadFiles(dev, sid,
			sources,
			destination,
			func(fi *TransferredFileInfo, err error) error {
				So(err, ShouldBeNil)
				So(fi, ShouldNotBeNil)
				So(fi.StartTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)
				So(fi.LatestSentTime.UnixNano(), ShouldBeGreaterThan, currentSentTime)
				So(fi.Speed, ShouldBeGreaterThanOrEqualTo, 0)
				So(fi.FilesSent, ShouldBeGreaterThanOrEqualTo, 0)

				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)
				So(fi.FileInfo.FullPath, ShouldStartWith, destination)

				So(fi.FileInfo.FullPath, ShouldEndWith, dirList[fi.FilesSent-1])

				currentSentTime = fi.LatestSentTime.UnixNano()
				currentSentFiles = fi.FilesSent

				return nil
			},
		)

		So(err, ShouldBeNil)
		So(currentSentFiles, ShouldEqual, 2)
		So(totalFiles, ShouldEqual, 2)
		So(totalFiles, ShouldEqual, currentSentFiles)
		So(totalSize, ShouldEqual, 15)

		fi, err := GetObjectFromPath(dev, sid, destination)
		So(err, ShouldBeNil)

		So(objectIdDest, ShouldEqual, fi.ObjectId)
	})

	Convey("Multiple Files | same name | Random destination | UploadFiles", t, func() {
		// test the directories: '/mtp-test-files/temp_dir/test_UploadFiles/{random}'

		uploadFile1 := getTestMocksAsset("mock_dir1/a.txt")
		uploadFile2 := getTestMocksAsset("mock_dir1/1/a.txt")
		sources := []string{uploadFile1, uploadFile2}

		randFName := fmt.Sprintf("%x", rand.Int31())
		destination := getFullPath("/mtp-test-files/temp_dir/test_UploadFiles", randFName)

		dirList := []string{"/a.txt", "/a.txt"}

		var currentSentTime int64
		var currentSentFiles int
		objectIdDest, totalFiles, totalSize, err := UploadFiles(dev, sid,
			sources,
			destination,
			func(fi *TransferredFileInfo, err error) error {
				So(err, ShouldBeNil)
				So(fi, ShouldNotBeNil)
				So(fi.StartTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)
				So(fi.LatestSentTime.UnixNano(), ShouldBeGreaterThan, currentSentTime)
				So(fi.Speed, ShouldBeGreaterThanOrEqualTo, 0)
				So(fi.FilesSent, ShouldBeGreaterThanOrEqualTo, 0)

				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)
				So(fi.FileInfo.FullPath, ShouldStartWith, destination)

				So(fi.FileInfo.FullPath, ShouldEndWith, dirList[fi.FilesSent-1])

				currentSentTime = fi.LatestSentTime.UnixNano()
				currentSentFiles = fi.FilesSent

				return nil
			},
		)

		So(err, ShouldBeNil)
		So(currentSentFiles, ShouldEqual, 2)
		So(totalFiles, ShouldEqual, 2)
		So(totalFiles, ShouldEqual, currentSentFiles)
		So(totalSize, ShouldEqual, 17)

		fi, err := GetObjectFromPath(dev, sid, destination)
		So(err, ShouldBeNil)

		So(objectIdDest, ShouldEqual, fi.ObjectId)
	})

	var _destination string
	Convey("Directories and Files | Random destination | UploadFiles", t, func() {
		// test the directories: '/mtp-test-files/temp_dir/test_UploadFiles/{random}'

		uploadFile1 := getTestMocksAsset("mock_dir1/")
		uploadFile2 := getTestMocksAsset("mock_dir1/1/a.txt")
		sources := []string{uploadFile1, uploadFile2}

		randFName := fmt.Sprintf("%x", rand.Int31())
		_destination = getFullPath("/mtp-test-files/temp_dir/test_UploadFiles", randFName)

		dirList := []string{
			"/mock_dir1/1/a.txt", "/mock_dir1/2/b.txt", "/mock_dir1/3/2/b.txt", "/mock_dir1/3/b.txt", "/mock_dir1/a.txt", "a.txt",
		}

		var currentSentTime int64
		var currentSentFiles int
		objectIdDest, totalFiles, totalSize, err := UploadFiles(dev, sid,
			sources,
			_destination,
			func(fi *TransferredFileInfo, err error) error {
				So(err, ShouldBeNil)
				So(fi, ShouldNotBeNil)
				So(fi.StartTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)
				So(fi.LatestSentTime.UnixNano(), ShouldBeGreaterThan, currentSentTime)
				So(fi.Speed, ShouldBeGreaterThanOrEqualTo, 0)
				So(fi.FilesSent, ShouldBeGreaterThanOrEqualTo, 0)

				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)
				So(fi.FileInfo.FullPath, ShouldStartWith, _destination)

				So(fi.FileInfo.FullPath, ShouldEndWith, dirList[fi.FilesSent-1])

				currentSentTime = fi.LatestSentTime.UnixNano()
				currentSentFiles = fi.FilesSent

				return nil
			},
		)

		So(err, ShouldBeNil)
		So(currentSentFiles, ShouldEqual, 6)
		So(totalFiles, ShouldEqual, 6)
		So(totalFiles, ShouldEqual, currentSentFiles)
		So(totalSize, ShouldEqual, 43)

		fi, err := GetObjectFromPath(dev, sid, _destination)
		So(err, ShouldBeNil)

		So(objectIdDest, ShouldEqual, fi.ObjectId)
	})

	Convey("Directories and Files | Previous destination | UploadFiles", t, func() {
		uploadFile1 := getTestMocksAsset("mock_dir1/")
		uploadFile2 := getTestMocksAsset("mock_dir1/1/a.txt")
		sources := []string{uploadFile1, uploadFile2}

		dirList := []string{
			"/mock_dir1/1/a.txt", "/mock_dir1/2/b.txt", "/mock_dir1/3/2/b.txt", "/mock_dir1/3/b.txt", "/mock_dir1/a.txt", "a.txt",
		}

		var currentSentTime int64
		var currentSentFiles int
		objectIdDest, totalFiles, totalSize, err := UploadFiles(dev, sid,
			sources,
			_destination,
			func(fi *TransferredFileInfo, err error) error {
				So(err, ShouldBeNil)
				So(fi, ShouldNotBeNil)
				So(fi.StartTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)
				So(fi.LatestSentTime.UnixNano(), ShouldBeGreaterThan, currentSentTime)
				So(fi.Speed, ShouldBeGreaterThanOrEqualTo, 0)
				So(fi.FilesSent, ShouldBeGreaterThanOrEqualTo, 0)

				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)
				So(fi.FileInfo.FullPath, ShouldStartWith, _destination)

				So(fi.FileInfo.FullPath, ShouldEndWith, dirList[fi.FilesSent-1])

				currentSentTime = fi.LatestSentTime.UnixNano()
				currentSentFiles = fi.FilesSent

				return nil
			},
		)

		So(err, ShouldBeNil)
		So(currentSentFiles, ShouldEqual, 6)
		So(totalFiles, ShouldEqual, 6)
		So(totalFiles, ShouldEqual, currentSentFiles)
		So(totalSize, ShouldEqual, 43)

		fi, err := GetObjectFromPath(dev, sid, _destination)
		So(err, ShouldBeNil)

		So(objectIdDest, ShouldEqual, fi.ObjectId)
	})

	Convey("Invalid source | Random destination | UploadFiles | should throw an error ", t, func() {
		// test the directories: '/mtp-test-files/temp_dir/test_UploadFiles/{random}'

		uploadFile1 := getTestMocksAsset("mock_dir1/")
		uploadFile2 := "fake/1/111a.txt"
		sources := []string{uploadFile1, uploadFile2}

		randFName := fmt.Sprintf("%x", rand.Int31())
		_destination = getFullPath("/mtp-test-files/temp_dir/test_UploadFiles", randFName)

		dirList := []string{
			"/mock_dir1/1/a.txt", "/mock_dir1/2/b.txt", "/mock_dir1/3/2/b.txt", "/mock_dir1/3/b.txt", "/mock_dir1/a.txt", "a.txt",
		}

		var currentSentTime int64
		var currentSentFiles int
		objectIdDest, totalFiles, totalSize, err := UploadFiles(dev, sid,
			sources,
			_destination,
			func(fi *TransferredFileInfo, err error) error {
				So(err, ShouldBeNil)
				So(fi, ShouldNotBeNil)
				So(fi.StartTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)
				So(fi.LatestSentTime.UnixNano(), ShouldBeGreaterThan, currentSentTime)
				So(fi.Speed, ShouldBeGreaterThanOrEqualTo, 0)
				So(fi.FilesSent, ShouldBeGreaterThanOrEqualTo, 0)

				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)
				So(fi.FileInfo.FullPath, ShouldStartWith, _destination)

				So(fi.FileInfo.FullPath, ShouldEndWith, dirList[fi.FilesSent-1])

				currentSentTime = fi.LatestSentTime.UnixNano()
				currentSentFiles = fi.FilesSent

				return nil
			},
		)

		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(currentSentFiles, ShouldEqual, 5)
		So(totalFiles, ShouldEqual, 5)
		So(totalFiles, ShouldEqual, currentSentFiles)
		So(totalSize, ShouldEqual, 35)

		fi, err := GetObjectFromPath(dev, sid, _destination)
		So(err, ShouldBeNil)

		So(objectIdDest, ShouldEqual, fi.ObjectId)
	})

	Convey("Callback return error | UploadFiles | should throw an error ", t, func() {
		// test the directories: '/mtp-test-files/temp_dir/test_UploadFiles/{random}'

		uploadFile1 := getTestMocksAsset("mock_dir1/")
		sources := []string{uploadFile1}

		randFName := fmt.Sprintf("%x", rand.Int31())
		_destination = getFullPath("/mtp-test-files/temp_dir/test_UploadFiles", randFName)

		dirList := []string{
			"/mock_dir1/1/a.txt", "/mock_dir1/2/b.txt", "/mock_dir1/3/2/b.txt", "/mock_dir1/3/b.txt", "/mock_dir1/a.txt",
		}

		var currentSentTime int64
		var currentSentFiles int
		objectIdDest, totalFiles, totalSize, err := UploadFiles(dev, sid,
			sources, _destination,
			func(fi *TransferredFileInfo, err error) error {
				So(err, ShouldBeNil)
				So(fi, ShouldNotBeNil)
				So(fi.StartTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)
				So(fi.LatestSentTime.UnixNano(), ShouldBeGreaterThan, currentSentTime)
				So(fi.Speed, ShouldBeGreaterThanOrEqualTo, 0)
				So(fi.FilesSent, ShouldBeGreaterThanOrEqualTo, 0)

				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)
				So(fi.FileInfo.FullPath, ShouldStartWith, _destination)

				So(fi.FileInfo.FullPath, ShouldEndWith, dirList[fi.FilesSent-1])

				currentSentTime = fi.LatestSentTime.UnixNano()
				currentSentFiles = fi.FilesSent

				return FilePermissionError{error: fmt.Errorf("some random error")}
			},
		)

		So(err, ShouldHaveSameTypeAs, FileTransferError{})
		So(currentSentFiles, ShouldEqual, 1)
		So(totalFiles, ShouldEqual, 1)
		So(totalFiles, ShouldEqual, currentSentFiles)
		So(totalSize, ShouldEqual, 8)

		fi, err := GetObjectFromPath(dev, sid, _destination)
		So(err, ShouldBeNil)

		So(objectIdDest, ShouldEqual, fi.ObjectId)
	})

	Dispose(dev)
}

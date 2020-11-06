package mtpx

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestDownloadFiles(t *testing.T) {
	dev, err := Initialize(Init{})
	if err != nil {
		log.Panic(err)
	}

	storages, err := FetchStorages(dev)
	if err != nil {
		log.Panic(err)
	}

	sid := storages[0].Sid

	Convey("General | DownloadFiles", t, func() {
		// test the directories: '/mtp-test-files/mock_dir1/'

		destination := newTempMocksDir("test_DownloadTest", true)
		sourceFile1 := "/mtp-test-files/mock_dir1/"
		sources := []string{sourceFile1}

		var currentSentTime int64
		var currentSentFiles int
		totalFiles, totalSize, err := DownloadFiles(dev, sid,
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
				So(fi.FileInfo.FullPath, ShouldStartWith, sources[0])

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
	})

	Convey("Single directory | DownloadFiles", t, func() {
		// test the directories: '/mtp-test-files/mock_dir1/'

		destination := newTempMocksDir("test_DownloadTest", true)
		sourceFile1 := "/mtp-test-files/mock_dir1/"
		sources := []string{sourceFile1}

		dirList := []string{"/mtp-test-files/mock_dir1/1/a.txt", "/mtp-test-files/mock_dir1/a.txt", "/mtp-test-files/mock_dir1/3/b.txt", "/mtp-test-files/mock_dir1/3/2/b.txt", "/mtp-test-files/mock_dir1/2/b.txt"}

		var currentSentTime int64
		var currentSentFiles int
		totalFiles, totalSize, err := DownloadFiles(dev, sid,
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
				So(fi.FileInfo.FullPath, ShouldStartWith, sources[0])

				contains, index := StringContains(dirList, fi.FileInfo.FullPath)
				So(contains, ShouldEqual, true)
				dirList = RemoveIndex(dirList, index)

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

		// walk the destination directory on device and verify
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

		_count := 0
		err = filepath.Walk(getFullPath(destination, "mock_dir1"), func(path string, info os.FileInfo, err error) error {
			So(path, ShouldEndWith, dirList1[_count])

			_count += 1
			return nil
		})

		So(err, ShouldBeNil)

	})

	Convey("Multiple directories | DownloadFiles", t, func() {
		// test the directories: '/mtp-test-files/mock_dir1' and '/mtp-test-files/mock_dir2'

		destination := newTempMocksDir("test_DownloadTest", true)
		sourceFile1 := "/mtp-test-files/mock_dir1/"
		sourceFile2 := "/mtp-test-files/mock_dir2/"
		sources := []string{sourceFile1, sourceFile2}

		dirList := []string{
			"/mtp-test-files/mock_dir1/1/a.txt", "/mtp-test-files/mock_dir1/a.txt", "/mtp-test-files/mock_dir1/3/b.txt", "/mtp-test-files/mock_dir1/3/2/b.txt", "/mtp-test-files/mock_dir1/2/b.txt",
			"/mtp-test-files/mock_dir2/1/a.txt", "/mtp-test-files/mock_dir2/a.txt", "/mtp-test-files/mock_dir2/3/b.txt", "/mtp-test-files/mock_dir2/3/2/b.txt", "/mtp-test-files/mock_dir2/2/b.txt",
		}

		var currentSentTime int64
		var currentSentFiles int
		count := 0
		totalFiles, totalSize, err := DownloadFiles(dev, sid,
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

				sIndex := 0
				if count > 4 {
					sIndex += 1
				}

				So(fi.FileInfo.FullPath, ShouldStartWith, sources[sIndex])
				contains, index := StringContains(dirList, fi.FileInfo.FullPath)
				So(contains, ShouldEqual, true)
				dirList = RemoveIndex(dirList, index)

				currentSentTime = fi.LatestSentTime.UnixNano()
				currentSentFiles = fi.FilesSent

				count += 1

				return nil
			},
		)

		So(err, ShouldBeNil)
		So(currentSentFiles, ShouldEqual, 10)
		So(totalFiles, ShouldEqual, 10)
		So(totalFiles, ShouldEqual, currentSentFiles)
		So(totalSize, ShouldEqual, 70)
	})

	Convey("Multiple directories | same name | DownloadFiles", t, func() {
		// test the directories: '/mtp-test-files/mock_dir1' and '/mtp-test-files/mock_dir1'

		destination := newTempMocksDir("test_DownloadTest", true)
		sourceFile1 := "/mtp-test-files/mock_dir1/"
		sourceFile2 := "/mtp-test-files/mock_dir1/"
		sources := []string{sourceFile1, sourceFile2}

		dirList := []string{
			"/mtp-test-files/mock_dir1/1/a.txt", "/mtp-test-files/mock_dir1/a.txt", "/mtp-test-files/mock_dir1/3/b.txt", "/mtp-test-files/mock_dir1/3/2/b.txt", "/mtp-test-files/mock_dir1/2/b.txt",
			"/mtp-test-files/mock_dir1/1/a.txt", "/mtp-test-files/mock_dir1/a.txt", "/mtp-test-files/mock_dir1/3/b.txt", "/mtp-test-files/mock_dir1/3/2/b.txt", "/mtp-test-files/mock_dir1/2/b.txt",
		}

		var currentSentTime int64
		var currentSentFiles int
		totalFiles, totalSize, err := DownloadFiles(dev, sid,
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

				So(fi.FileInfo.FullPath, ShouldStartWith, sources[0])
				contains, index := StringContains(dirList, fi.FileInfo.FullPath)
				So(contains, ShouldEqual, true)
				dirList = RemoveIndex(dirList, index)

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
	})

	Convey("Single File | DownloadFiles", t, func() {
		// test the directories: '/mtp-test-files/mock_dir1/a.txt'

		destination := newTempMocksDir("test_DownloadTest", true)
		sourceFile1 := "/mtp-test-files/mock_dir1/a.txt"
		sources := []string{sourceFile1}

		dirList := []string{"/a.txt"}

		var currentSentTime int64
		var currentSentFiles int
		totalFiles, totalSize, err := DownloadFiles(dev, sid,
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

				So(fi.FileInfo.FullPath, ShouldStartWith, sources[0])
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
	})

	Convey("Multiple Files | DownloadFiles", t, func() {
		// test the directories: '/mtp-test-files/mock_dir1/a.txt' and '/mtp-test-files/mock_dir1/2/b.txt'

		destination := newTempMocksDir("test_DownloadTest", true)
		sourceFile1 := "/mtp-test-files/mock_dir1/a.txt"
		sourceFile2 := "/mtp-test-files/mock_dir1/2/b.txt"
		sources := []string{sourceFile1, sourceFile2}

		dirList := []string{"/a.txt", "/b.txt"}

		var currentSentTime int64
		var currentSentFiles int
		count := 0
		totalFiles, totalSize, err := DownloadFiles(dev, sid,
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

				sIndex := 0
				if count > 0 {
					sIndex += 1
				}

				So(fi.FileInfo.FullPath, ShouldStartWith, sources[sIndex])
				So(fi.FileInfo.FullPath, ShouldEndWith, dirList[fi.FilesSent-1])

				currentSentTime = fi.LatestSentTime.UnixNano()
				currentSentFiles = fi.FilesSent

				count += 1

				return nil
			},
		)

		So(err, ShouldBeNil)
		So(currentSentFiles, ShouldEqual, 2)
		So(totalFiles, ShouldEqual, 2)
		So(totalFiles, ShouldEqual, currentSentFiles)
		So(totalSize, ShouldEqual, 15)
	})

	Convey("Multiple Files | same name | DownloadFiles", t, func() {
		// test the directories: '/mtp-test-files/mock_dir1/a.txt' and '/mtp-test-files/mock_dir1/1/a.txt'

		destination := newTempMocksDir("test_DownloadTest", true)
		sourceFile1 := "/mtp-test-files/mock_dir1/a.txt"
		sourceFile2 := "/mtp-test-files/mock_dir1/1/a.txt"
		sources := []string{sourceFile1, sourceFile2}

		dirList := []string{"/a.txt", "/a.txt"}

		var currentSentTime int64
		var currentSentFiles int
		count := 0
		totalFiles, totalSize, err := DownloadFiles(dev, sid,
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

				sIndex := 0
				if count > 0 {
					sIndex += 1
				}
				So(fi.FileInfo.FullPath, ShouldStartWith, sources[sIndex])
				So(fi.FileInfo.FullPath, ShouldEndWith, dirList[fi.FilesSent-1])

				currentSentTime = fi.LatestSentTime.UnixNano()
				currentSentFiles = fi.FilesSent
				count += 1

				return nil
			},
		)

		So(err, ShouldBeNil)
		So(currentSentFiles, ShouldEqual, 2)
		So(totalFiles, ShouldEqual, 2)
		So(totalFiles, ShouldEqual, currentSentFiles)
		So(totalSize, ShouldEqual, 17)
	})

	var _destination string
	Convey("Directories and Files | DownloadFiles", t, func() {
		// test the directories: '/mtp-test-files/mock_dir1/' and '/mtp-test-files/mock_dir1/1/a.txt'

		_destination = newTempMocksDir("test_DownloadTest", true)
		sourceFile1 := "/mtp-test-files/mock_dir1/"
		sourceFile2 := "/mtp-test-files/mock_dir1/1/a.txt"
		sources := []string{sourceFile1, sourceFile2}

		dirList := []string{
			"/mtp-test-files/mock_dir1/1/a.txt", "/mtp-test-files/mock_dir1/a.txt", "/mtp-test-files/mock_dir1/3/b.txt", "/mtp-test-files/mock_dir1/3/2/b.txt", "/mtp-test-files/mock_dir1/2/b.txt",

			"/mtp-test-files/mock_dir1/1/a.txt",
		}

		var currentSentTime int64
		var currentSentFiles int
		count := 0
		totalFiles, totalSize, err := DownloadFiles(dev, sid,
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
				sIndex := 0
				if count > 4 {
					sIndex += 1
				}
				So(fi.FileInfo.FullPath, ShouldStartWith, sources[sIndex])

				contains, index := StringContains(dirList, fi.FileInfo.FullPath)
				So(contains, ShouldEqual, true)
				dirList = RemoveIndex(dirList, index)


				currentSentTime = fi.LatestSentTime.UnixNano()
				currentSentFiles = fi.FilesSent
				count += 1

				return nil
			},
		)

		So(err, ShouldBeNil)
		So(currentSentFiles, ShouldEqual, 6)
		So(totalFiles, ShouldEqual, 6)
		So(totalFiles, ShouldEqual, currentSentFiles)
		So(totalSize, ShouldEqual, 43)
	})

	Convey("Directories and Files | overwritting the same destination | DownloadFiles", t, func() {
		sourceFile1 := "/mtp-test-files/mock_dir1/"
		sourceFile2 := "/mtp-test-files/mock_dir1/1/a.txt"
		sources := []string{sourceFile1, sourceFile2}

		dirList := []string{
			"/mtp-test-files/mock_dir1/1/a.txt", "/mtp-test-files/mock_dir1/a.txt", "/mtp-test-files/mock_dir1/3/b.txt", "/mtp-test-files/mock_dir1/3/2/b.txt", "/mtp-test-files/mock_dir1/2/b.txt",

			"/mtp-test-files/mock_dir1/1/a.txt",
		}

		var currentSentTime int64
		var currentSentFiles int
		count := 0
		totalFiles, totalSize, err := DownloadFiles(dev, sid,
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
				sIndex := 0
				if count > 4 {
					sIndex += 1
				}
				So(fi.FileInfo.FullPath, ShouldStartWith, sources[sIndex])

				contains, index := StringContains(dirList, fi.FileInfo.FullPath)
				So(contains, ShouldEqual, true)
				dirList = RemoveIndex(dirList, index)

				currentSentTime = fi.LatestSentTime.UnixNano()
				currentSentFiles = fi.FilesSent
				count += 1

				return nil
			},
		)

		So(err, ShouldBeNil)
		So(currentSentFiles, ShouldEqual, 6)
		So(totalFiles, ShouldEqual, 6)
		So(totalFiles, ShouldEqual, currentSentFiles)
		So(totalSize, ShouldEqual, 43)
	})

	Convey("Invalid source | DownloadFiles | should throw an error ", t, func() {
		// test the directories: '/mtp-test-files/mock_dir1/' and 'fake/1/111a.txt'

		destination := newTempMocksDir("test_DownloadTest", true)

		sourceFile1 := "/mtp-test-files/mock_dir1/"
		sourceFile2 := "fake/1/111a.txt"
		sources := []string{sourceFile1, sourceFile2}

		dirList := []string{
			"/mtp-test-files/mock_dir1/1/a.txt", "/mtp-test-files/mock_dir1/a.txt", "/mtp-test-files/mock_dir1/3/b.txt", "/mtp-test-files/mock_dir1/3/2/b.txt", "/mtp-test-files/mock_dir1/2/b.txt",
		}

		var currentSentTime int64
		var currentSentFiles int
		count := 0
		totalFiles, totalSize, err := DownloadFiles(dev, sid,
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
				sIndex := 0
				if count > 4 {
					sIndex += 1
				}
				So(fi.FileInfo.FullPath, ShouldStartWith, sources[sIndex])

				contains, index := StringContains(dirList, fi.FileInfo.FullPath)
				So(contains, ShouldEqual, true)
				dirList = RemoveIndex(dirList, index)

				currentSentTime = fi.LatestSentTime.UnixNano()
				currentSentFiles = fi.FilesSent
				count += 1

				return nil
			},
		)

		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(currentSentFiles, ShouldEqual, 5)
		So(totalFiles, ShouldEqual, 5)
		So(totalFiles, ShouldEqual, currentSentFiles)
		So(totalSize, ShouldEqual, 35)
	})

	Convey("Callback return error | DownloadFiles | should throw an error ", t, func() {
		// test the directories: '/mtp-test-files/mock_dir1/'

		destination := newTempMocksDir("test_DownloadTest", true)

		sourceFile1 := "/mtp-test-files/mock_dir1/"
		sources := []string{sourceFile1}

		dirList := []string{
			"/mock_dir1/1/a.txt", "/mock_dir1/a.txt", "/mock_dir1/3/b.txt", "/mock_dir1/3/2/b.txt", "/mock_dir1/2/b.txt",
			"a.txt",
		}

		var currentSentTime int64
		var currentSentFiles int
		count := 0
		totalFiles, totalSize, err := DownloadFiles(dev, sid,
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
				sIndex := 0
				if count > 4 {
					sIndex += 1
				}
				So(fi.FileInfo.FullPath, ShouldStartWith, sources[sIndex])

				So(fi.FileInfo.FullPath, ShouldEndWith, dirList[fi.FilesSent-1])

				currentSentTime = fi.LatestSentTime.UnixNano()
				currentSentFiles = fi.FilesSent
				count += 1

				return FileObjectError{error: fmt.Errorf("some random error")}
			},
		)

		So(err, ShouldHaveSameTypeAs, FileTransferError{})
		So(currentSentFiles, ShouldEqual, 1)
		So(totalFiles, ShouldEqual, 1)
		So(totalFiles, ShouldEqual, currentSentFiles)
		So(totalSize, ShouldEqual, 8)
	})

	Dispose(dev)
}

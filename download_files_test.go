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

	Convey("Single directory | DownloadFiles", t, func() {
		// test directories: '/mtp-test-files/mock_dir1/'
		destination := newTempMocksDir("test_DownloadTest", true)
		sourceFile1 := "/mtp-test-files/mock_dir1/"
		sources := []string{sourceFile1}

		dirList := []string{"/mtp-test-files/mock_dir1/1/a.txt", "/mtp-test-files/mock_dir1/a.txt", "/mtp-test-files/mock_dir1/3/b.txt", "/mtp-test-files/mock_dir1/3/2/b.txt", "/mtp-test-files/mock_dir1/2/b.txt"}

		var prevLatestSentTime int64
		var prevFilesSent int64
		var prevFilesSentProgress float32
		var prevCurrentSentProgress float32
		var prevBulkSentProgress float32
		var prevBulkSent int64
		var status TransferStatus
		totalFiles, totalSize, err := DownloadFiles(dev, sid,
			sources,
			destination,
			false,
			func(fi *FileInfo, err error) error {
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

				So(fi.Speed, ShouldBeGreaterThanOrEqualTo, 0)
				So(fi.FilesSent, ShouldEqual, prevFilesSent)

				if fi.Status == InProgress {
					prevFilesSent += 1
				}

				So(fi.TotalDirectories, ShouldEqual, 0)

				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)
				So(fi.FileInfo.FullPath, ShouldStartWith, sources[0])

				if fi.Status == InProgress {
					contains, index := StringContains(dirList, fi.FileInfo.FullPath)
					So(contains, ShouldEqual, true)
					dirList = RemoveIndex(dirList, index)
				}

				So(fi.FilesSentProgress, ShouldBeGreaterThanOrEqualTo, prevFilesSentProgress)
				prevFilesSentProgress = fi.FilesSentProgress

				// current progress tests
				So(fi.ActiveFileSize.Total, ShouldBeGreaterThan, 0)
				So(fi.ActiveFileSize.Sent, ShouldEqual, fi.ActiveFileSize.Total)

				So(fi.ActiveFileSize.Progress, ShouldBeGreaterThanOrEqualTo, prevCurrentSentProgress)
				prevCurrentSentProgress = fi.ActiveFileSize.Progress

				// bulk progress tests
				So(fi.BulkFileSize.Total, ShouldEqual, 0)
				So(fi.BulkFileSize.Sent, ShouldBeGreaterThanOrEqualTo, prevBulkSent)
				prevBulkSent = fi.BulkFileSize.Sent

				So(fi.BulkFileSize.Progress, ShouldBeGreaterThanOrEqualTo, prevBulkSentProgress)
				prevBulkSentProgress = fi.BulkFileSize.Progress

				status = fi.Status

				return nil
			},
		)

		So(status, ShouldEqual, Completed)
		So(err, ShouldBeNil)
		So(prevFilesSent, ShouldEqual, 5)
		So(totalFiles, ShouldEqual, 5)
		So(totalFiles, ShouldEqual, prevFilesSent)
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

	Convey("Single Large file | DownloadFiles", t, func() {
		// test directories: '4mb_txt_file'
		destination := newTempMocksDir("test_DownloadTest", true)
		sourceFile1 := "/mtp-test-files/4mb_txt_file"
		sources := []string{sourceFile1}

		dirList := []string{"/mtp-test-files/4mb_txt_file"}

		var prevLatestSentTime int64
		var prevFilesSent int64
		var prevFilesSentProgress float32
		var prevCurrentSentProgress float32
		var prevBulkSentProgress float32
		var prevBulkSent int64
		var status TransferStatus
		var prevSentFileFullPath string
		totalFiles, totalSize, err := DownloadFiles(dev, sid,
			sources,
			destination,
			false,
			func(fi *FileInfo, err error) error {
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

				So(fi.Speed, ShouldBeGreaterThanOrEqualTo, 0)
				if prevSentFileFullPath != fi.FileInfo.FullPath {
					So(fi.FilesSent, ShouldEqual, prevFilesSent)

					if fi.Status == InProgress {
						contains, index := StringContains(dirList, fi.FileInfo.FullPath)
						So(contains, ShouldEqual, true)
						dirList = RemoveIndex(dirList, index)

						prevFilesSent += 1
					}

				}
				prevSentFileFullPath = fi.FileInfo.FullPath

				So(fi.TotalDirectories, ShouldEqual, 0)

				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)
				So(fi.FileInfo.FullPath, ShouldStartWith, sources[0])

				So(fi.FilesSentProgress, ShouldBeGreaterThanOrEqualTo, prevFilesSentProgress)
				prevFilesSentProgress = fi.FilesSentProgress

				// current progress tests
				So(fi.ActiveFileSize.Total, ShouldBeGreaterThan, 0)
				So(fi.ActiveFileSize.Sent, ShouldBeLessThanOrEqualTo, fi.ActiveFileSize.Total)

				So(fi.ActiveFileSize.Progress, ShouldBeGreaterThanOrEqualTo, prevCurrentSentProgress)
				prevCurrentSentProgress = fi.ActiveFileSize.Progress

				// bulk progress tests
				So(fi.BulkFileSize.Total, ShouldEqual, 0)
				So(fi.BulkFileSize.Sent, ShouldBeGreaterThanOrEqualTo, prevBulkSent)
				prevBulkSent = fi.BulkFileSize.Sent

				So(fi.BulkFileSize.Progress, ShouldBeGreaterThanOrEqualTo, prevBulkSentProgress)
				prevBulkSentProgress = fi.BulkFileSize.Progress

				status = fi.Status

				return nil
			},
		)

		So(status, ShouldEqual, Completed)
		So(err, ShouldBeNil)
		So(prevFilesSent, ShouldEqual, 1)
		So(totalFiles, ShouldEqual, 1)
		So(totalFiles, ShouldEqual, prevFilesSent)
		So(totalSize, ShouldEqual, 4194304)

		// walk the destination directory on device and verify
		dirList1 := []string{
			"4mb_txt_file",
		}

		_count := 0
		err = filepath.Walk(getFullPath(destination, "4mb_txt_file"), func(path string, info os.FileInfo, err error) error {
			So(path, ShouldEndWith, dirList1[_count])

			_count += 1
			return nil
		})

		So(err, ShouldBeNil)
	})

	Convey("Multiple Large files | DownloadFiles", t, func() {
		// test directories: '4mb_txt_file'
		destination := newTempMocksDir("test_DownloadTest", true)
		sourceFile1 := "/mtp-test-files/4mb_txt_file"
		sourceFile2 := "/mtp-test-files/4mb_txt_file_2"
		sources := []string{sourceFile1, sourceFile2}

		dirList := []string{
			"/mtp-test-files/4mb_txt_file",
			"/mtp-test-files/4mb_txt_file_2",
		}

		var prevLatestSentTime int64
		var prevFilesSent int64
		var prevFilesSentProgress float32
		var prevCurrentSentProgress float32
		var prevBulkSentProgress float32
		var prevBulkSent int64
		var status TransferStatus
		var prevSentFileFullPath string
		totalFiles, totalSize, err := DownloadFiles(dev, sid,
			sources,
			destination,
			false,
			func(fi *FileInfo, err error) error {
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

				So(fi.Speed, ShouldBeGreaterThanOrEqualTo, 0)
				if prevSentFileFullPath != fi.FileInfo.FullPath {
					So(fi.FilesSent, ShouldEqual, prevFilesSent)

					if fi.Status == InProgress {
						contains, index := StringContains(dirList, fi.FileInfo.FullPath)
						So(contains, ShouldEqual, true)
						dirList = RemoveIndex(dirList, index)

						prevFilesSent += 1
					}

				}
				prevSentFileFullPath = fi.FileInfo.FullPath

				So(fi.TotalDirectories, ShouldEqual, 0)

				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)
				So(fi.FileInfo.FullPath, ShouldStartWith, sources[0])

				So(fi.FilesSentProgress, ShouldBeGreaterThanOrEqualTo, prevFilesSentProgress)
				prevFilesSentProgress = fi.FilesSentProgress

				// current progress tests
				So(fi.ActiveFileSize.Total, ShouldBeGreaterThan, 0)
				So(fi.ActiveFileSize.Sent, ShouldBeLessThanOrEqualTo, fi.ActiveFileSize.Total)

				if prevSentFileFullPath != fi.FileInfo.FullPath {
					So(fi.ActiveFileSize.Progress, ShouldBeGreaterThanOrEqualTo, prevCurrentSentProgress)
				}
				prevCurrentSentProgress = fi.ActiveFileSize.Progress

				// bulk progress tests
				So(fi.BulkFileSize.Total, ShouldEqual, 0)
				So(fi.BulkFileSize.Sent, ShouldBeGreaterThanOrEqualTo, prevBulkSent)
				prevBulkSent = fi.BulkFileSize.Sent

				So(fi.BulkFileSize.Progress, ShouldBeGreaterThanOrEqualTo, prevBulkSentProgress)
				prevBulkSentProgress = fi.BulkFileSize.Progress

				status = fi.Status

				return nil
			},
		)

		So(status, ShouldEqual, Completed)
		So(err, ShouldBeNil)
		So(prevFilesSent, ShouldEqual, 2)
		So(totalFiles, ShouldEqual, 2)
		So(totalFiles, ShouldEqual, prevFilesSent)
		So(totalSize, ShouldEqual, 4194304*2)

		// walk the destination directory on device and verify
		dirList1 := []string{
			"4mb_txt_file",
			"4mb_txt_file_2",
		}

		_count := 0
		err = filepath.Walk(getFullPath(destination, "4mb_txt_file"), func(path string, info os.FileInfo, err error) error {
			So(path, ShouldEndWith, dirList1[_count])

			_count += 1
			return nil
		})

		So(err, ShouldBeNil)
	})

	Convey("Multiple directories | Random destination | DownloadFiles", t, func() {
		// test directories: '/mtp-test-files/mock_dir1' and '/mtp-test-files/mock_dir2'

		destination := newTempMocksDir("test_DownloadTest", true)
		sourceFile1 := "/mtp-test-files/mock_dir1/"
		sourceFile2 := "/mtp-test-files/mock_dir2/"
		sources := []string{sourceFile1, sourceFile2}

		dirList := []string{
			"/mtp-test-files/mock_dir1/1/a.txt", "/mtp-test-files/mock_dir1/a.txt", "/mtp-test-files/mock_dir1/3/b.txt", "/mtp-test-files/mock_dir1/3/2/b.txt", "/mtp-test-files/mock_dir1/2/b.txt",
			"/mtp-test-files/mock_dir2/1/a.txt", "/mtp-test-files/mock_dir2/a.txt", "/mtp-test-files/mock_dir2/3/b.txt", "/mtp-test-files/mock_dir2/3/2/b.txt", "/mtp-test-files/mock_dir2/2/b.txt",
		}

		var prevLatestSentTime int64
		var prevFilesSent int64
		var prevFilesSentProgress float32
		var prevCurrentSentProgress float32
		var prevBulkSentProgress float32
		var prevBulkSent int64
		var status TransferStatus
		totalFiles, totalSize, err := DownloadFiles(dev, sid,
			sources,
			destination,
			false,
			func(fi *FileInfo, err error) error {
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

				So(fi.Speed, ShouldBeGreaterThanOrEqualTo, 0)
				So(fi.FilesSent, ShouldEqual, prevFilesSent)

				if fi.Status == InProgress {
					prevFilesSent += 1
				}

				So(fi.TotalDirectories, ShouldEqual, 0)
				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)

				if fi.Status == InProgress {
					contains, index := StringContains(dirList, fi.FileInfo.FullPath)
					So(contains, ShouldEqual, true)
					dirList = RemoveIndex(dirList, index)
				}

				So(fi.FilesSentProgress, ShouldBeGreaterThanOrEqualTo, prevFilesSentProgress)
				prevFilesSentProgress = fi.FilesSentProgress

				// current progress tests
				So(fi.ActiveFileSize.Total, ShouldBeGreaterThan, 0)
				So(fi.ActiveFileSize.Sent, ShouldEqual, fi.ActiveFileSize.Total)

				So(fi.ActiveFileSize.Progress, ShouldBeGreaterThanOrEqualTo, prevCurrentSentProgress)
				prevCurrentSentProgress = fi.ActiveFileSize.Progress

				// bulk progress tests
				So(fi.BulkFileSize.Total, ShouldEqual, 0)
				So(fi.BulkFileSize.Sent, ShouldBeGreaterThanOrEqualTo, prevBulkSent)
				prevBulkSent = fi.BulkFileSize.Sent

				So(fi.BulkFileSize.Progress, ShouldBeGreaterThanOrEqualTo, prevBulkSentProgress)
				prevBulkSentProgress = fi.BulkFileSize.Progress
				status = fi.Status

				return nil
			},
		)

		So(err, ShouldBeNil)
		So(prevFilesSent, ShouldEqual, 5*2)
		So(totalFiles, ShouldEqual, 5*2)
		So(totalFiles, ShouldEqual, prevFilesSent)
		So(totalSize, ShouldEqual, 35*2)
		So(status, ShouldEqual, Completed)

		//walk the destination directory on device and verify
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

		_count := 0
		err = filepath.Walk(destination, func(path string, info os.FileInfo, err error) error {
			if _count < 1 {
				return nil
			}

			So(dirList1[_count], ShouldEndWith, path)

			_count += 1
			return nil
		})

		So(err, ShouldBeNil)
	})

	Convey("Multiple directories | same name | DownloadFiles", t, func() {
		// test directories: '/mtp-test-files/mock_dir1' and '/mtp-test-files/mock_dir1'

		destination := newTempMocksDir("test_DownloadTest", true)
		sourceFile1 := "/mtp-test-files/mock_dir1/"
		sourceFile2 := "/mtp-test-files/mock_dir1/"
		sources := []string{sourceFile1, sourceFile2}

		dirList := []string{
			"/mtp-test-files/mock_dir1/1/a.txt", "/mtp-test-files/mock_dir1/a.txt", "/mtp-test-files/mock_dir1/3/b.txt", "/mtp-test-files/mock_dir1/3/2/b.txt", "/mtp-test-files/mock_dir1/2/b.txt",
			"/mtp-test-files/mock_dir1/1/a.txt", "/mtp-test-files/mock_dir1/a.txt", "/mtp-test-files/mock_dir1/3/b.txt", "/mtp-test-files/mock_dir1/3/2/b.txt", "/mtp-test-files/mock_dir1/2/b.txt",
		}

		var prevLatestSentTime int64
		var prevFilesSent int64
		var prevFilesSentProgress float32
		var prevCurrentSentProgress float32
		var prevBulkSentProgress float32
		var prevBulkSent int64
		var status TransferStatus
		totalFiles, totalSize, err := DownloadFiles(dev, sid,
			sources,
			destination,
			false,
			func(fi *FileInfo, err error) error {
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

				So(fi.Speed, ShouldBeGreaterThanOrEqualTo, 0)
				So(fi.FilesSent, ShouldEqual, prevFilesSent)

				if fi.Status == InProgress {
					prevFilesSent += 1
				}

				So(fi.TotalDirectories, ShouldEqual, 0)
				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)

				if fi.Status == InProgress {
					contains, index := StringContains(dirList, fi.FileInfo.FullPath)
					So(contains, ShouldEqual, true)
					dirList = RemoveIndex(dirList, index)
				}

				So(fi.FilesSentProgress, ShouldBeGreaterThanOrEqualTo, prevFilesSentProgress)
				prevFilesSentProgress = fi.FilesSentProgress

				// current progress tests
				So(fi.ActiveFileSize.Total, ShouldBeGreaterThan, 0)
				So(fi.ActiveFileSize.Sent, ShouldEqual, fi.ActiveFileSize.Total)

				So(fi.ActiveFileSize.Progress, ShouldBeGreaterThanOrEqualTo, prevCurrentSentProgress)
				prevCurrentSentProgress = fi.ActiveFileSize.Progress

				// bulk progress tests
				So(fi.BulkFileSize.Total, ShouldEqual, 0)
				So(fi.BulkFileSize.Sent, ShouldBeGreaterThanOrEqualTo, prevBulkSent)
				prevBulkSent = fi.BulkFileSize.Sent

				So(fi.BulkFileSize.Progress, ShouldBeGreaterThanOrEqualTo, prevBulkSentProgress)
				prevBulkSentProgress = fi.BulkFileSize.Progress

				status = fi.Status
				return nil
			},
		)

		So(status, ShouldEqual, Completed)
		So(err, ShouldBeNil)
		So(prevFilesSent, ShouldEqual, 5*2)
		So(totalFiles, ShouldEqual, 5*2)
		So(totalFiles, ShouldEqual, prevFilesSent)
		So(totalSize, ShouldEqual, 35*2)

		//walk the destination directory on device and verify
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

		_count := 0
		err = filepath.Walk(destination, func(path string, info os.FileInfo, err error) error {
			if _count < 1 {
				return nil
			}

			So(dirList1[_count], ShouldEndWith, path)

			_count += 1
			return nil
		})

		So(err, ShouldBeNil)
	})

	Convey("Single File | DownloadFiles", t, func() {
		// test directories: '/mtp-test-files/mock_dir1/a.txt'

		destination := newTempMocksDir("test_DownloadTest", true)
		sourceFile1 := "/mtp-test-files/mock_dir1/a.txt"
		sources := []string{sourceFile1}

		dirList := []string{
			"/mtp-test-files/mock_dir1/a.txt",
		}

		var prevLatestSentTime int64
		var prevFilesSent int64
		var prevFilesSentProgress float32
		var prevCurrentSentProgress float32
		var prevBulkSentProgress float32
		var prevBulkSent int64
		var status TransferStatus
		totalFiles, totalSize, err := DownloadFiles(dev, sid,
			sources,
			destination,
			false,
			func(fi *FileInfo, err error) error {
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

				So(fi.Speed, ShouldBeGreaterThanOrEqualTo, 0)
				So(fi.FilesSent, ShouldEqual, prevFilesSent)

				if fi.Status == InProgress {
					prevFilesSent += 1
				}

				So(fi.TotalDirectories, ShouldEqual, 0)
				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)

				if fi.Status == InProgress {
					contains, index := StringContains(dirList, fi.FileInfo.FullPath)
					So(contains, ShouldEqual, true)
					dirList = RemoveIndex(dirList, index)
				}

				So(fi.FilesSentProgress, ShouldBeGreaterThanOrEqualTo, prevFilesSentProgress)
				prevFilesSentProgress = fi.FilesSentProgress

				// current progress tests
				So(fi.ActiveFileSize.Total, ShouldBeGreaterThan, 0)
				So(fi.ActiveFileSize.Sent, ShouldEqual, fi.ActiveFileSize.Total)

				So(fi.ActiveFileSize.Progress, ShouldBeGreaterThanOrEqualTo, prevCurrentSentProgress)
				prevCurrentSentProgress = fi.ActiveFileSize.Progress

				// bulk progress tests
				So(fi.BulkFileSize.Total, ShouldEqual, 0)
				So(fi.BulkFileSize.Sent, ShouldBeGreaterThanOrEqualTo, prevBulkSent)
				prevBulkSent = fi.BulkFileSize.Sent

				So(fi.BulkFileSize.Progress, ShouldBeGreaterThanOrEqualTo, prevBulkSentProgress)
				prevBulkSentProgress = fi.BulkFileSize.Progress

				status = fi.Status
				return nil
			},
		)

		So(status, ShouldEqual, Completed)
		So(err, ShouldBeNil)
		So(prevFilesSent, ShouldEqual, 1)
		So(totalFiles, ShouldEqual, 1)
		So(totalFiles, ShouldEqual, prevFilesSent)
		So(totalSize, ShouldEqual, 9)

		//walk the destination directory on device and verify
		dirList1 := []string{
			"/mock_dir1/a.txt",
		}

		_count := 0
		err = filepath.Walk(destination, func(path string, info os.FileInfo, err error) error {
			if _count < 1 {
				return nil
			}

			So(dirList1[_count], ShouldEndWith, path)

			_count += 1
			return nil
		})

		So(err, ShouldBeNil)
	})

	Convey("Multiple Files | DownloadFiles", t, func() {
		// test directories: '/mtp-test-files/mock_dir1/a.txt' and '/mtp-test-files/mock_dir1/2/b.txt'

		destination := newTempMocksDir("test_DownloadTest", true)
		sourceFile1 := "/mtp-test-files/mock_dir1/a.txt"
		sourceFile2 := "/mtp-test-files/mock_dir1/2/b.txt"
		sources := []string{sourceFile1, sourceFile2}

		dirList := []string{
			"/mtp-test-files/mock_dir1/a.txt",
			"/mtp-test-files/mock_dir1/2/b.txt",
		}

		var prevLatestSentTime int64
		var prevFilesSent int64
		var prevFilesSentProgress float32
		var prevCurrentSentProgress float32
		var prevBulkSentProgress float32
		var prevBulkSent int64
		var status TransferStatus
		totalFiles, totalSize, err := DownloadFiles(dev, sid,
			sources,
			destination,
			false,
			func(fi *FileInfo, err error) error {
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

				So(fi.Speed, ShouldBeGreaterThanOrEqualTo, 0)
				So(fi.FilesSent, ShouldEqual, prevFilesSent)

				if fi.Status == InProgress {
					prevFilesSent += 1
				}

				So(fi.TotalDirectories, ShouldEqual, 0)
				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)

				if fi.Status == InProgress {
					contains, index := StringContains(dirList, fi.FileInfo.FullPath)
					So(contains, ShouldEqual, true)
					dirList = RemoveIndex(dirList, index)
				}

				So(fi.FilesSentProgress, ShouldBeGreaterThanOrEqualTo, prevFilesSentProgress)
				prevFilesSentProgress = fi.FilesSentProgress

				// current progress tests
				So(fi.ActiveFileSize.Total, ShouldBeGreaterThan, 0)
				So(fi.ActiveFileSize.Sent, ShouldEqual, fi.ActiveFileSize.Total)

				So(fi.ActiveFileSize.Progress, ShouldBeGreaterThanOrEqualTo, prevCurrentSentProgress)
				prevCurrentSentProgress = fi.ActiveFileSize.Progress

				// bulk progress tests
				So(fi.BulkFileSize.Total, ShouldEqual, 0)
				So(fi.BulkFileSize.Sent, ShouldBeGreaterThanOrEqualTo, prevBulkSent)
				prevBulkSent = fi.BulkFileSize.Sent

				So(fi.BulkFileSize.Progress, ShouldBeGreaterThanOrEqualTo, prevBulkSentProgress)
				prevBulkSentProgress = fi.BulkFileSize.Progress

				status = fi.Status
				return nil
			},
		)

		So(status, ShouldEqual, Completed)
		So(err, ShouldBeNil)
		So(prevFilesSent, ShouldEqual, 1*2)
		So(totalFiles, ShouldEqual, 1*2)
		So(totalFiles, ShouldEqual, prevFilesSent)
		So(totalSize, ShouldEqual, 15)

		//walk the destination directory on device and verify
		dirList1 := []string{
			"/mock_dir1/a.txt",
			"/mock_dir1/2/b.txt",
		}

		_count := 0
		err = filepath.Walk(destination, func(path string, info os.FileInfo, err error) error {
			if _count < 1 {
				return nil
			}

			So(dirList1[_count], ShouldEndWith, path)

			_count += 1
			return nil
		})

		So(err, ShouldBeNil)
	})

	Convey("Multiple Files | same name | DownloadFiles", t, func() {
		// test directories: '/mtp-test-files/mock_dir1/a.txt' and '/mtp-test-files/mock_dir1/1/a.txt'

		destination := newTempMocksDir("test_DownloadTest", true)
		sourceFile1 := "/mtp-test-files/mock_dir1/a.txt"
		sourceFile2 := "/mtp-test-files/mock_dir1/1/a.txt"
		sources := []string{sourceFile1, sourceFile2}

		dirList := []string{
			"/mtp-test-files/mock_dir1/a.txt",
			"/mtp-test-files/mock_dir1/1/a.txt",
		}

		var prevLatestSentTime int64
		var prevFilesSent int64
		var prevFilesSentProgress float32
		var prevCurrentSentProgress float32
		var prevBulkSentProgress float32
		var prevBulkSent int64
		var status TransferStatus
		totalFiles, totalSize, err := DownloadFiles(dev, sid,
			sources,
			destination,
			false,
			func(fi *FileInfo, err error) error {
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

				So(fi.Speed, ShouldBeGreaterThanOrEqualTo, 0)
				So(fi.FilesSent, ShouldEqual, prevFilesSent)

				if fi.Status == InProgress {
					prevFilesSent += 1
				}

				So(fi.TotalDirectories, ShouldEqual, 0)
				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)

				if fi.Status == InProgress {
					contains, index := StringContains(dirList, fi.FileInfo.FullPath)
					So(contains, ShouldEqual, true)
					dirList = RemoveIndex(dirList, index)
				}

				So(fi.FilesSentProgress, ShouldBeGreaterThanOrEqualTo, prevFilesSentProgress)
				prevFilesSentProgress = fi.FilesSentProgress

				// current progress tests
				So(fi.ActiveFileSize.Total, ShouldBeGreaterThan, 0)
				So(fi.ActiveFileSize.Sent, ShouldEqual, fi.ActiveFileSize.Total)

				So(fi.ActiveFileSize.Progress, ShouldBeGreaterThanOrEqualTo, prevCurrentSentProgress)
				prevCurrentSentProgress = fi.ActiveFileSize.Progress

				// bulk progress tests
				So(fi.BulkFileSize.Total, ShouldEqual, 0)
				So(fi.BulkFileSize.Sent, ShouldBeGreaterThanOrEqualTo, prevBulkSent)
				prevBulkSent = fi.BulkFileSize.Sent

				So(fi.BulkFileSize.Progress, ShouldBeGreaterThanOrEqualTo, prevBulkSentProgress)
				prevBulkSentProgress = fi.BulkFileSize.Progress

				status = fi.Status
				return nil
			},
		)

		So(status, ShouldEqual, Completed)
		So(err, ShouldBeNil)
		So(prevFilesSent, ShouldEqual, 1*2)
		So(totalFiles, ShouldEqual, 1*2)
		So(totalFiles, ShouldEqual, prevFilesSent)
		So(totalSize, ShouldEqual, 17)

		//walk the destination directory on device and verify
		dirList1 := []string{
			"/mock_dir1/a.txt",
			"/mock_dir1/1/a.txt",
		}

		_count := 0
		err = filepath.Walk(destination, func(path string, info os.FileInfo, err error) error {
			if _count < 1 {
				return nil
			}

			So(dirList1[_count], ShouldEndWith, path)

			_count += 1
			return nil
		})

		So(err, ShouldBeNil)
	})

	var _destination string
	Convey("Directories and Files | DownloadFiles", t, func() {
		// test directories: '/mtp-test-files/mock_dir1/' and '/mtp-test-files/mock_dir1/1/a.txt'

		_destination = newTempMocksDir("test_DownloadTest", true)
		sourceFile1 := "/mtp-test-files/mock_dir1/"
		sourceFile2 := "/mtp-test-files/mock_dir1/1/a.txt"
		sources := []string{sourceFile1, sourceFile2}

		dirList := []string{
			"/mtp-test-files/mock_dir1/1/a.txt", "/mtp-test-files/mock_dir1/a.txt", "/mtp-test-files/mock_dir1/3/b.txt", "/mtp-test-files/mock_dir1/3/2/b.txt", "/mtp-test-files/mock_dir1/2/b.txt",
			"/mtp-test-files/mock_dir1/1/a.txt",
		}

		var prevLatestSentTime int64
		var prevFilesSent int64
		var prevFilesSentProgress float32
		var prevCurrentSentProgress float32
		var prevBulkSentProgress float32
		var prevBulkSent int64
		var status TransferStatus
		totalFiles, totalSize, err := DownloadFiles(dev, sid,
			sources,
			_destination,
			false,
			func(fi *FileInfo, err error) error {
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

				So(fi.Speed, ShouldBeGreaterThanOrEqualTo, 0)
				So(fi.FilesSent, ShouldEqual, prevFilesSent)

				if fi.Status == InProgress {
					prevFilesSent += 1
				}

				So(fi.TotalDirectories, ShouldEqual, 0)
				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)

				if fi.Status == InProgress {
					contains, index := StringContains(dirList, fi.FileInfo.FullPath)
					So(contains, ShouldEqual, true)
					dirList = RemoveIndex(dirList, index)
				}

				So(fi.FilesSentProgress, ShouldBeGreaterThanOrEqualTo, prevFilesSentProgress)
				prevFilesSentProgress = fi.FilesSentProgress

				// current progress tests
				So(fi.ActiveFileSize.Total, ShouldBeGreaterThan, 0)
				So(fi.ActiveFileSize.Sent, ShouldEqual, fi.ActiveFileSize.Total)

				So(fi.ActiveFileSize.Progress, ShouldBeGreaterThanOrEqualTo, prevCurrentSentProgress)
				prevCurrentSentProgress = fi.ActiveFileSize.Progress

				// bulk progress tests
				So(fi.BulkFileSize.Total, ShouldEqual, 0)
				So(fi.BulkFileSize.Sent, ShouldBeGreaterThanOrEqualTo, prevBulkSent)
				prevBulkSent = fi.BulkFileSize.Sent

				So(fi.BulkFileSize.Progress, ShouldBeGreaterThanOrEqualTo, prevBulkSentProgress)
				prevBulkSentProgress = fi.BulkFileSize.Progress

				status = fi.Status
				return nil
			},
		)

		So(status, ShouldEqual, Completed)
		So(err, ShouldBeNil)
		So(prevFilesSent, ShouldEqual, 6)
		So(totalFiles, ShouldEqual, 6)
		So(totalFiles, ShouldEqual, prevFilesSent)
		So(totalSize, ShouldEqual, 43)

		//walk the destination directory on device and verify
		dirList1 := []string{
			"/mock_dir1/1/a.txt", "/mock_dir1/a.txt", "/mock_dir1/3/b.txt", "/mock_dir1/3/2/b.txt", "/mock_dir1/2/b.txt",
			"/mock_dir1/1/a.txt",
		}

		_count := 0
		err = filepath.Walk(_destination, func(path string, info os.FileInfo, err error) error {
			if _count < 1 {
				return nil
			}

			So(dirList1[_count], ShouldEndWith, path)

			_count += 1
			return nil
		})

		So(err, ShouldBeNil)
	})

	Convey("Directories and Files | overwritting the same destination | DownloadFiles", t, func() {
		sourceFile1 := "/mtp-test-files/mock_dir1/"
		sourceFile2 := "/mtp-test-files/mock_dir1/1/a.txt"
		sources := []string{sourceFile1, sourceFile2}

		dirList := []string{
			"/mtp-test-files/mock_dir1/1/a.txt", "/mtp-test-files/mock_dir1/a.txt", "/mtp-test-files/mock_dir1/3/b.txt", "/mtp-test-files/mock_dir1/3/2/b.txt", "/mtp-test-files/mock_dir1/2/b.txt",

			"/mtp-test-files/mock_dir1/1/a.txt",
		}

		var prevLatestSentTime int64
		var prevFilesSent int64
		var prevFilesSentProgress float32
		var prevCurrentSentProgress float32
		var prevBulkSentProgress float32
		var prevBulkSent int64
		var status TransferStatus
		totalFiles, totalSize, err := DownloadFiles(dev, sid,
			sources,
			_destination,
			false,
			func(fi *FileInfo, err error) error {
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

				So(fi.Speed, ShouldBeGreaterThanOrEqualTo, 0)
				So(fi.FilesSent, ShouldEqual, prevFilesSent)

				if fi.Status == InProgress {
					prevFilesSent += 1
				}

				So(fi.TotalDirectories, ShouldEqual, 0)
				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)

				if fi.Status == InProgress {
					contains, index := StringContains(dirList, fi.FileInfo.FullPath)
					So(contains, ShouldEqual, true)
					dirList = RemoveIndex(dirList, index)
				}

				So(fi.FilesSentProgress, ShouldBeGreaterThanOrEqualTo, prevFilesSentProgress)
				prevFilesSentProgress = fi.FilesSentProgress

				// current progress tests
				So(fi.ActiveFileSize.Total, ShouldBeGreaterThan, 0)
				So(fi.ActiveFileSize.Sent, ShouldEqual, fi.ActiveFileSize.Total)

				So(fi.ActiveFileSize.Progress, ShouldBeGreaterThanOrEqualTo, prevCurrentSentProgress)
				prevCurrentSentProgress = fi.ActiveFileSize.Progress

				// bulk progress tests
				So(fi.BulkFileSize.Total, ShouldEqual, 0)
				So(fi.BulkFileSize.Sent, ShouldBeGreaterThanOrEqualTo, prevBulkSent)
				prevBulkSent = fi.BulkFileSize.Sent

				So(fi.BulkFileSize.Progress, ShouldBeGreaterThanOrEqualTo, prevBulkSentProgress)
				prevBulkSentProgress = fi.BulkFileSize.Progress

				status = fi.Status
				return nil
			},
		)

		So(status, ShouldEqual, Completed)
		So(err, ShouldBeNil)
		So(prevFilesSent, ShouldEqual, 6)
		So(totalFiles, ShouldEqual, 6)
		So(totalFiles, ShouldEqual, prevFilesSent)
		So(totalSize, ShouldEqual, 43)

		//walk the destination directory on device and verify
		dirList1 := []string{
			"/mock_dir1/1/a.txt", "/mock_dir1/a.txt", "/mock_dir1/3/b.txt", "/mock_dir1/3/2/b.txt", "/mock_dir1/2/b.txt",
			"/mock_dir1/1/a.txt",
		}

		_count := 0
		err = filepath.Walk(_destination, func(path string, info os.FileInfo, err error) error {
			if _count < 1 {
				return nil
			}

			So(dirList1[_count], ShouldEndWith, path)

			_count += 1
			return nil
		})

		So(err, ShouldBeNil)
	})

	Convey("Multiple Large files and Multiple assorted files and directories | DownloadFiles", t, func() {
		// test directories: '4mb_txt_file'
		destination := newTempMocksDir("test_DownloadTest", true)
		sourceFile1 := "/mtp-test-files/4mb_txt_file"
		sourceFile2 := "/mtp-test-files/4mb_txt_file_2"
		sourceFile3 := "/mtp-test-files/mock_dir1/a.txt"
		sourceFile4 := "/mtp-test-files/mock_dir1"
		sources := []string{sourceFile1, sourceFile2, sourceFile3, sourceFile4}

		dirList := []string{
			"/mtp-test-files/4mb_txt_file",
			"/mtp-test-files/4mb_txt_file_2",
			"/mtp-test-files/mock_dir1/a.txt",
			"/mtp-test-files/mock_dir1/1/a.txt",
			"/mtp-test-files/mock_dir1/2/b.txt",
			"/mtp-test-files/mock_dir1/3/2/b.txt",
			"/mtp-test-files/mock_dir1/3/b.txt",
			"/mtp-test-files/mock_dir1/a.txt",
		}

		var prevLatestSentTime int64
		var prevFilesSent int64
		var prevFilesSentProgress float32
		var prevCurrentSentProgress float32
		var prevBulkSentProgress float32
		var prevBulkSent int64
		var status TransferStatus
		var prevSentFileFullPath string
		totalFiles, totalSize, err := DownloadFiles(dev, sid,
			sources,
			destination,
			false,
			func(fi *FileInfo, err error) error {
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

				So(fi.Speed, ShouldBeGreaterThanOrEqualTo, 0)
				if prevSentFileFullPath != fi.FileInfo.FullPath {
					So(fi.FilesSent, ShouldEqual, prevFilesSent)

					if fi.Status == InProgress {
						contains, index := StringContains(dirList, fi.FileInfo.FullPath)
						So(contains, ShouldEqual, true)
						dirList = RemoveIndex(dirList, index)

						prevFilesSent += 1
					}

				}
				prevSentFileFullPath = fi.FileInfo.FullPath

				So(fi.TotalDirectories, ShouldEqual, 0)
				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)

				So(fi.FilesSentProgress, ShouldBeGreaterThanOrEqualTo, prevFilesSentProgress)
				prevFilesSentProgress = fi.FilesSentProgress

				// current progress tests
				So(fi.ActiveFileSize.Total, ShouldBeGreaterThan, 0)
				So(fi.ActiveFileSize.Sent, ShouldBeLessThanOrEqualTo, fi.ActiveFileSize.Total)

				if prevSentFileFullPath != fi.FileInfo.FullPath {
					So(fi.ActiveFileSize.Progress, ShouldBeGreaterThanOrEqualTo, prevCurrentSentProgress)
				}
				prevCurrentSentProgress = fi.ActiveFileSize.Progress

				// bulk progress tests
				So(fi.BulkFileSize.Total, ShouldEqual, 0)
				So(fi.BulkFileSize.Sent, ShouldBeGreaterThanOrEqualTo, prevBulkSent)
				prevBulkSent = fi.BulkFileSize.Sent

				So(fi.BulkFileSize.Progress, ShouldBeGreaterThanOrEqualTo, prevBulkSentProgress)
				prevBulkSentProgress = fi.BulkFileSize.Progress

				status = fi.Status

				return nil
			},
		)

		So(err, ShouldBeNil)
		So(status, ShouldEqual, Completed)
		So(prevFilesSent, ShouldEqual, 2*4)
		So(totalFiles, ShouldEqual, 2*4)
		So(totalFiles, ShouldEqual, prevFilesSent)
		So(totalSize, ShouldEqual, 8388652)

		// walk the destination directory on device and verify
		dirList1 := []string{
			"4mb_txt_file",
			"4mb_txt_file_2",
			"a.txt",
			"mock_dir1/1/a.txt",
			"mock_dir1/a.txt",
			"mock_dir1/3/b.txt",
			"mock_dir1/3/2/b.txt",
			"mock_dir1/2/b.txt",
			"mock_dir1/1/a.txt",
			"mock_dir1/a.txt",
		}

		_count := 0
		err = filepath.Walk(getFullPath(destination, "4mb_txt_file"), func(path string, info os.FileInfo, err error) error {
			So(path, ShouldEndWith, dirList1[_count])

			_count += 1
			return nil
		})

		So(err, ShouldBeNil)
	})

	Convey("preprocessing=true | Single directory | DownloadFiles", t, func() {
		// test directories: '/mtp-test-files/mock_dir1'

		destination := newTempMocksDir("test_DownloadTest", true)
		sourceFile1 := "/mtp-test-files/mock_dir1/"
		sources := []string{sourceFile1}

		dirList := []string{
			"/mtp-test-files/mock_dir1/1/a.txt", "/mtp-test-files/mock_dir1/a.txt", "/mtp-test-files/mock_dir1/3/b.txt", "/mtp-test-files/mock_dir1/3/2/b.txt", "/mtp-test-files/mock_dir1/2/b.txt",
		}

		preprocessingDirList := []string{
			"/mtp-test-files/mock_dir1/1/a.txt",
			"/mtp-test-files/mock_dir1/a.txt",
			"/mtp-test-files/mock_dir1/3/b.txt",
			"/mtp-test-files/mock_dir1/3/2/b.txt",
			"/mtp-test-files/mock_dir1/2/b.txt",
		}

		var prevLatestSentTime int64
		var prevFilesSent int64
		var prevFilesSentProgress float32
		var prevCurrentSentProgress float32
		var prevBulkSentProgress float32
		var prevBulkSent int64
		var status TransferStatus
		preprocessingIndex := 0
		totalFiles, totalSize, err := DownloadFiles(dev, sid,
			sources,
			destination,
			true,
			func(fi *FileInfo, err error) error {
				if err != nil {
					return err
				}

				So((*fi).FullPath, ShouldEndWith, preprocessingDirList[preprocessingIndex])

				preprocessingIndex += 1

				return nil
			},
			func(fi *ProgressInfo, err error) error {
				So(err, ShouldBeNil)
				So(fi, ShouldNotBeNil)

				So(fi.StartTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)
				So(fi.LatestSentTime.UnixNano(), ShouldBeGreaterThan, prevLatestSentTime)
				prevLatestSentTime = fi.LatestSentTime.UnixNano()

				So(fi.Speed, ShouldBeGreaterThanOrEqualTo, 0)
				So(fi.FilesSent, ShouldEqual, prevFilesSent)

				if fi.Status == InProgress {
					prevFilesSent += 1
				}

				So(fi.TotalDirectories, ShouldEqual, 4)
				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)

				if fi.Status == InProgress {
					contains, index := StringContains(dirList, fi.FileInfo.FullPath)
					So(contains, ShouldEqual, true)
					dirList = RemoveIndex(dirList, index)
				}

				So(fi.FilesSentProgress, ShouldBeGreaterThanOrEqualTo, prevFilesSentProgress)
				prevFilesSentProgress = fi.FilesSentProgress

				// current progress tests
				So(fi.ActiveFileSize.Total, ShouldBeGreaterThan, 0)
				So(fi.ActiveFileSize.Sent, ShouldEqual, fi.ActiveFileSize.Total)

				So(fi.ActiveFileSize.Progress, ShouldBeGreaterThanOrEqualTo, prevCurrentSentProgress)
				prevCurrentSentProgress = fi.ActiveFileSize.Progress

				// bulk progress tests
				So(fi.BulkFileSize.Total, ShouldEqual, 35)
				So(fi.BulkFileSize.Sent, ShouldBeGreaterThanOrEqualTo, prevBulkSent)
				prevBulkSent = fi.BulkFileSize.Sent

				So(fi.BulkFileSize.Progress, ShouldBeGreaterThanOrEqualTo, prevBulkSentProgress)
				prevBulkSentProgress = fi.BulkFileSize.Progress
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

		//walk the destination directory on device and verify
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

		_count := 0
		err = filepath.Walk(destination, func(path string, info os.FileInfo, err error) error {
			if _count < 1 {
				return nil
			}

			So(dirList1[_count], ShouldEndWith, path)

			_count += 1
			return nil
		})

		So(err, ShouldBeNil)
	})

	Convey("preprocessing=true | Multiple Large files and Multiple assorted files and directories | DownloadFiles", t, func() {
		// test directories: '4mb_txt_file'
		destination := newTempMocksDir("test_DownloadTest", true)
		sourceFile1 := "/mtp-test-files/4mb_txt_file"
		sourceFile2 := "/mtp-test-files/4mb_txt_file_2"
		sourceFile3 := "/mtp-test-files/mock_dir1/a.txt"
		sourceFile4 := "/mtp-test-files/mock_dir1"
		sources := []string{sourceFile1, sourceFile2, sourceFile3, sourceFile4}

		dirList := []string{
			"/mtp-test-files/4mb_txt_file",
			"/mtp-test-files/4mb_txt_file_2",
			"/mtp-test-files/mock_dir1/a.txt",
			"/mtp-test-files/mock_dir1/1/a.txt",
			"/mtp-test-files/mock_dir1/2/b.txt",
			"/mtp-test-files/mock_dir1/3/2/b.txt",
			"/mtp-test-files/mock_dir1/3/b.txt",
			"/mtp-test-files/mock_dir1/a.txt",
		}

		preprocessingDirList := []string{
			"/mtp-test-files/4mb_txt_file",
			"/mtp-test-files/4mb_txt_file_2",
			"/mtp-test-files/mock_dir1/a.txt",
			"/mtp-test-files/mock_dir1/1/a.txt",
			"/mtp-test-files/mock_dir1/a.txt",
			"/mtp-test-files/mock_dir1/3/b.txt",
			"/mtp-test-files/mock_dir1/3/2/b.txt",
			"/mtp-test-files/mock_dir1/2/b.txt",
		}

		var prevLatestSentTime int64
		var prevFilesSent int64
		var prevFilesSentProgress float32
		var prevCurrentSentProgress float32
		var prevBulkSentProgress float32
		var prevBulkSent int64
		var status TransferStatus
		var prevSentFileFullPath string
		preprocessingIndex := 0
		totalFiles, totalSize, err := DownloadFiles(dev, sid,
			sources,
			destination,
			true,
			func(fi *FileInfo, err error) error {
				if err != nil {
					return err
				}

				So((*fi).FullPath, ShouldEndWith, preprocessingDirList[preprocessingIndex])

				preprocessingIndex += 1

				return nil
			},
			func(fi *ProgressInfo, err error) error {
				So(err, ShouldBeNil)
				So(fi, ShouldNotBeNil)

				So(fi.StartTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)
				So(fi.LatestSentTime.UnixNano(), ShouldBeGreaterThan, prevLatestSentTime)
				prevLatestSentTime = fi.LatestSentTime.UnixNano()

				So(fi.Speed, ShouldBeGreaterThanOrEqualTo, 0)
				if prevSentFileFullPath != fi.FileInfo.FullPath {
					So(fi.FilesSent, ShouldEqual, prevFilesSent)

					if fi.Status == InProgress {
						contains, index := StringContains(dirList, fi.FileInfo.FullPath)
						So(contains, ShouldEqual, true)
						dirList = RemoveIndex(dirList, index)

						prevFilesSent += 1
					}

				}
				prevSentFileFullPath = fi.FileInfo.FullPath

				So(fi.TotalDirectories, ShouldEqual, 4)
				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)

				So(fi.FilesSentProgress, ShouldBeGreaterThanOrEqualTo, prevFilesSentProgress)
				prevFilesSentProgress = fi.FilesSentProgress

				// current progress tests
				So(fi.ActiveFileSize.Total, ShouldBeGreaterThan, 0)
				So(fi.ActiveFileSize.Sent, ShouldBeLessThanOrEqualTo, fi.ActiveFileSize.Total)

				if prevSentFileFullPath != fi.FileInfo.FullPath {
					So(fi.ActiveFileSize.Progress, ShouldBeGreaterThanOrEqualTo, prevCurrentSentProgress)
				}
				prevCurrentSentProgress = fi.ActiveFileSize.Progress

				// bulk progress tests
				So(fi.BulkFileSize.Total, ShouldEqual, 8388652)
				So(fi.BulkFileSize.Sent, ShouldBeGreaterThanOrEqualTo, prevBulkSent)
				prevBulkSent = fi.BulkFileSize.Sent

				So(fi.BulkFileSize.Progress, ShouldBeGreaterThanOrEqualTo, prevBulkSentProgress)
				prevBulkSentProgress = fi.BulkFileSize.Progress

				status = fi.Status

				return nil
			},
		)

		So(err, ShouldBeNil)
		So(status, ShouldEqual, Completed)
		So(prevFilesSent, ShouldEqual, 2*4)
		So(totalFiles, ShouldEqual, 2*4)
		So(totalFiles, ShouldEqual, prevFilesSent)
		So(totalSize, ShouldEqual, 8388652)

		// walk the destination directory on device and verify
		dirList1 := []string{
			"4mb_txt_file",
			"4mb_txt_file_2",
			"a.txt",
			"mock_dir1/1/a.txt",
			"mock_dir1/a.txt",
			"mock_dir1/3/b.txt",
			"mock_dir1/3/2/b.txt",
			"mock_dir1/2/b.txt",
			"mock_dir1/1/a.txt",
			"mock_dir1/a.txt",
		}

		_count := 0
		err = filepath.Walk(getFullPath(destination, "4mb_txt_file"), func(path string, info os.FileInfo, err error) error {
			So(path, ShouldEndWith, dirList1[_count])

			_count += 1
			return nil
		})

		So(err, ShouldBeNil)
	})

	Convey("preprocessing=true | Invalid source | DownloadFiles | should throw an error ", t, func() {
		// test directories: '/mtp-test-files/mock_dir1/' and 'fake/1/111a.txt'
		destination := newTempMocksDir("test_DownloadTest", true)

		sourceFile1 := "/mtp-test-files/mock_dir1/"
		sourceFile2 := "fake/1/111a.txt"
		sources := []string{sourceFile1, sourceFile2}

		dirList := []string{
			"/mtp-test-files/mock_dir1/1/a.txt", "/mtp-test-files/mock_dir1/a.txt", "/mtp-test-files/mock_dir1/3/b.txt", "/mtp-test-files/mock_dir1/3/2/b.txt", "/mtp-test-files/mock_dir1/2/b.txt",
		}

		var prevLatestSentTime int64
		var prevFilesSent int64
		var prevFilesSentProgress float32
		var prevCurrentSentProgress float32
		var prevBulkSentProgress float32
		var prevBulkSent int64
		totalFiles, totalSize, err := DownloadFiles(dev, sid,
			sources,
			destination,
			true,
			func(fi *FileInfo, err error) error {
				if err != nil {
					return err
				}

				return InvalidPathError{error: fmt.Errorf("some error occured")}
			},
			func(fi *ProgressInfo, err error) error {
				So(err, ShouldBeNil)
				So(fi, ShouldNotBeNil)

				So(fi.StartTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)
				So(fi.LatestSentTime.UnixNano(), ShouldBeGreaterThan, prevLatestSentTime)
				prevLatestSentTime = fi.LatestSentTime.UnixNano()

				So(fi.Speed, ShouldBeGreaterThanOrEqualTo, 0)
				So(fi.FilesSent, ShouldEqual, prevFilesSent)

				if fi.Status == InProgress {
					prevFilesSent += 1
				}

				So(fi.TotalDirectories, ShouldEqual, 0)
				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)

				if fi.Status == InProgress {
					contains, index := StringContains(dirList, fi.FileInfo.FullPath)
					So(contains, ShouldEqual, true)
					dirList = RemoveIndex(dirList, index)
				}

				So(fi.FilesSentProgress, ShouldBeGreaterThanOrEqualTo, prevFilesSentProgress)
				prevFilesSentProgress = fi.FilesSentProgress

				// current progress tests
				So(fi.ActiveFileSize.Total, ShouldBeGreaterThan, 0)
				So(fi.ActiveFileSize.Sent, ShouldEqual, fi.ActiveFileSize.Total)

				So(fi.ActiveFileSize.Progress, ShouldBeGreaterThanOrEqualTo, prevCurrentSentProgress)
				prevCurrentSentProgress = fi.ActiveFileSize.Progress

				// bulk progress tests
				So(fi.BulkFileSize.Total, ShouldEqual, 0)
				So(fi.BulkFileSize.Sent, ShouldBeGreaterThanOrEqualTo, prevBulkSent)
				prevBulkSent = fi.BulkFileSize.Sent

				So(fi.BulkFileSize.Progress, ShouldBeGreaterThanOrEqualTo, prevBulkSentProgress)
				prevBulkSentProgress = fi.BulkFileSize.Progress

				return nil
			},
		)

		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(prevFilesSent, ShouldEqual, 0)
		So(totalFiles, ShouldEqual, 0)
		So(totalFiles, ShouldEqual, prevFilesSent)
		So(totalSize, ShouldEqual, 0)

		//walk the destination directory on device and verify
		_count := 0
		err = filepath.Walk(_destination, func(path string, info os.FileInfo, err error) error {
			if _count < 1 {
				return nil
			}

			// the code should not reach here
			So(true, ShouldNotEqual, true)

			_count += 1
			return nil
		})

		So(err, ShouldBeNil)
	})

	Convey("Invalid source | DownloadFiles | should throw an error ", t, func() {
		// test directories: '/mtp-test-files/mock_dir1/' and 'fake/1/111a.txt'
		destination := newTempMocksDir("test_DownloadTest", true)

		sourceFile1 := "/mtp-test-files/mock_dir1/"
		sourceFile2 := "fake/1/111a.txt"
		sources := []string{sourceFile1, sourceFile2}

		dirList := []string{
			"/mtp-test-files/mock_dir1/1/a.txt", "/mtp-test-files/mock_dir1/a.txt", "/mtp-test-files/mock_dir1/3/b.txt", "/mtp-test-files/mock_dir1/3/2/b.txt", "/mtp-test-files/mock_dir1/2/b.txt",
		}

		var prevLatestSentTime int64
		var prevFilesSent int64
		var prevFilesSentProgress float32
		var prevCurrentSentProgress float32
		var prevBulkSentProgress float32
		var prevBulkSent int64
		totalFiles, totalSize, err := DownloadFiles(dev, sid,
			sources,
			destination,
			false,
			func(fi *FileInfo, err error) error {
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

				So(fi.Speed, ShouldBeGreaterThanOrEqualTo, 0)
				So(fi.FilesSent, ShouldEqual, prevFilesSent)

				if fi.Status == InProgress {
					prevFilesSent += 1
				}

				So(fi.TotalDirectories, ShouldEqual, 0)
				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)

				if fi.Status == InProgress {
					contains, index := StringContains(dirList, fi.FileInfo.FullPath)
					So(contains, ShouldEqual, true)
					dirList = RemoveIndex(dirList, index)
				}

				So(fi.FilesSentProgress, ShouldBeGreaterThanOrEqualTo, prevFilesSentProgress)
				prevFilesSentProgress = fi.FilesSentProgress

				// current progress tests
				So(fi.ActiveFileSize.Total, ShouldBeGreaterThan, 0)
				So(fi.ActiveFileSize.Sent, ShouldEqual, fi.ActiveFileSize.Total)

				So(fi.ActiveFileSize.Progress, ShouldBeGreaterThanOrEqualTo, prevCurrentSentProgress)
				prevCurrentSentProgress = fi.ActiveFileSize.Progress

				// bulk progress tests
				So(fi.BulkFileSize.Total, ShouldEqual, 0)
				So(fi.BulkFileSize.Sent, ShouldBeGreaterThanOrEqualTo, prevBulkSent)
				prevBulkSent = fi.BulkFileSize.Sent

				So(fi.BulkFileSize.Progress, ShouldBeGreaterThanOrEqualTo, prevBulkSentProgress)
				prevBulkSentProgress = fi.BulkFileSize.Progress

				return nil
			},
		)

		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(prevFilesSent, ShouldEqual, 5)
		So(totalFiles, ShouldEqual, 5)
		So(totalFiles, ShouldEqual, prevFilesSent)
		So(totalSize, ShouldEqual, 35)

		//walk the destination directory on device and verify
		dirList1 := []string{
			"/mock_dir1/1/a.txt", "/mock_dir1/a.txt", "/mock_dir1/3/b.txt", "/mock_dir1/3/2/b.txt", "/mock_dir1/2/b.txt",
			"/mock_dir1/1/a.txt",
		}

		_count := 0
		err = filepath.Walk(_destination, func(path string, info os.FileInfo, err error) error {
			if _count < 1 {
				return nil
			}

			So(dirList1[_count], ShouldEndWith, path)

			_count += 1
			return nil
		})

		So(err, ShouldBeNil)
	})

	Convey("Callback return error | DownloadFiles | should throw an error ", t, func() {
		// test directories: '/mtp-test-files/mock_dir1/'

		destination := newTempMocksDir("test_DownloadTest", true)

		sourceFile1 := "/mtp-test-files/mock_dir1/"
		sources := []string{sourceFile1}

		dirList := []string{
			"/mtp-test-files/mock_dir1/1/a.txt",
		}

		var prevLatestSentTime int64
		var prevFilesSent int64
		var prevFilesSentProgress float32
		var prevCurrentSentProgress float32
		var prevBulkSentProgress float32
		var prevBulkSent int64
		totalFiles, totalSize, err := DownloadFiles(dev, sid,
			sources,
			destination,
			false,
			func(fi *FileInfo, err error) error {
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

				So(fi.Speed, ShouldBeGreaterThanOrEqualTo, 0)
				So(fi.FilesSent, ShouldEqual, prevFilesSent)

				if fi.Status == InProgress {
					prevFilesSent += 1
				}

				So(fi.TotalDirectories, ShouldEqual, 0)
				So(fi.FileInfo.ParentId, ShouldBeGreaterThan, 0)

				if fi.Status == InProgress {
					contains, index := StringContains(dirList, fi.FileInfo.FullPath)
					So(contains, ShouldEqual, true)
					dirList = RemoveIndex(dirList, index)
				}

				So(fi.FilesSentProgress, ShouldBeGreaterThanOrEqualTo, prevFilesSentProgress)
				prevFilesSentProgress = fi.FilesSentProgress

				// current progress tests
				So(fi.ActiveFileSize.Total, ShouldBeGreaterThan, 0)
				So(fi.ActiveFileSize.Sent, ShouldEqual, fi.ActiveFileSize.Total)

				So(fi.ActiveFileSize.Progress, ShouldBeGreaterThanOrEqualTo, prevCurrentSentProgress)
				prevCurrentSentProgress = fi.ActiveFileSize.Progress

				// bulk progress tests
				So(fi.BulkFileSize.Total, ShouldEqual, 0)
				So(fi.BulkFileSize.Sent, ShouldBeGreaterThanOrEqualTo, prevBulkSent)
				prevBulkSent = fi.BulkFileSize.Sent

				So(fi.BulkFileSize.Progress, ShouldBeGreaterThanOrEqualTo, prevBulkSentProgress)
				prevBulkSentProgress = fi.BulkFileSize.Progress

				return FileObjectError{error: fmt.Errorf("some random error")}
			},
		)

		So(err, ShouldHaveSameTypeAs, FileTransferError{})
		So(prevFilesSent, ShouldEqual, 1)
		So(totalFiles, ShouldEqual, 1)
		So(totalFiles, ShouldEqual, prevFilesSent)
		So(totalSize, ShouldEqual, 8)

		//walk the destination directory on device and verify
		dirList1 := []string{
			"/mock_dir1/1/a.txt",
		}

		_count := 0
		err = filepath.Walk(_destination, func(path string, info os.FileInfo, err error) error {
			if _count < 1 {
				return nil
			}

			So(dirList1[_count], ShouldEndWith, path)

			_count += 1
			return nil
		})

		So(err, ShouldBeNil)
	})

	Dispose(dev)
}

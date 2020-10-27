package main

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"testing"
)

func TestWalk(t *testing.T) {
	dev, err := Initialize(Init{})
	if err != nil {
		log.Panic(err)
	}

	storages, err := FetchStorages(dev)
	if err != nil {
		log.Panic(err)
	}

	sid := storages[0].sid

	Convey("Testing valid directory | Walk", t, func() {
		/////////////////
		// test the directory '/mtp-test-files'
		/////////////////
		fullPath := "/mtp-test-files"

		var children []*FileInfo
		objectId1, totalFiles1, err := Walk(dev, sid, fullPath, false,
			func(objectId uint32, fi *FileInfo, err error) error {
				So(err, ShouldBeNil)
				So(fi.FullPath, ShouldNotEqual, "/mtp-test-files")
				So(fi.FullPath, ShouldContainSubstring, "/mtp-test-files/")
				So(objectId, ShouldBeGreaterThan, 0)
				So(fi, ShouldNotBeNil)
				So(fi.ParentId, ShouldBeGreaterThan, 0)
				So(objectId, ShouldEqual, fi.ObjectId)

				children = append(children, fi)

				return nil
			})

		So(err, ShouldBeNil)
		So(totalFiles1, ShouldBeGreaterThanOrEqualTo, 4)

		// test if [objectId] == [objectId1] of '/mtp-test-files'
		fi, err := GetObjectFromPath(dev, sid, fullPath)
		So(err, ShouldBeNil)

		So(objectId1, ShouldEqual, fi.ObjectId)

		So(len(children), ShouldEqual, totalFiles1)

		/////////////////
		// test the directory '/mtp-test-files/'
		/////////////////
		fullPath = "/mtp-test-files/"
		children = []*FileInfo{}
		objectId2, totalFiles2, err := Walk(dev, sid, fullPath, false,
			func(objectId uint32, fi *FileInfo, err error) error {
				So(err, ShouldBeNil)
				// make sure that the first item is not the parent path itself
				So(fi.FullPath, ShouldNotEqual, "/mtp-test-files")
				So(fi.FullPath, ShouldContainSubstring, "/mtp-test-files/")
				So(objectId, ShouldBeGreaterThan, 0)
				So(fi, ShouldNotBeNil)
				So(fi.ParentId, ShouldBeGreaterThan, 0)
				So(objectId, ShouldEqual, fi.ObjectId)

				children = append(children, fi)

				return nil
			})

		So(err, ShouldBeNil)
		So(totalFiles2, ShouldBeGreaterThanOrEqualTo, totalFiles1)

		// test if [objectId2] == [objectId1] of [fullPath]
		fi, err = GetObjectFromPath(dev, sid, fullPath)
		So(err, ShouldBeNil)

		So(objectId1, ShouldEqual, fi.ObjectId)
		So(objectId1, ShouldEqual, objectId2)

		So(len(children), ShouldEqual, totalFiles2)

		/////////////////
		// test the directory 'mtp-test-files/'
		/////////////////
		fullPath = "mtp-test-files/"
		children = []*FileInfo{}
		objectId3, totalFiles3, err := Walk(dev, sid, fullPath, false,
			func(objectId uint32, fi *FileInfo, err error) error {
				So(err, ShouldBeNil)
				// make sure that the first item is not the parent path itself
				So(fi.FullPath, ShouldNotEqual, "/mtp-test-files")
				So(fi.FullPath, ShouldContainSubstring, "/mtp-test-files/")
				So(objectId, ShouldBeGreaterThan, 0)
				So(fi, ShouldNotBeNil)
				So(fi.ParentId, ShouldBeGreaterThan, 0)
				So(objectId, ShouldEqual, fi.ObjectId)

				children = append(children, fi)

				return nil
			})

		So(err, ShouldBeNil)
		So(totalFiles3, ShouldBeGreaterThanOrEqualTo, totalFiles1)

		// test if [objectId3] == [objectId] of [fullPath]
		fi, err = GetObjectFromPath(dev, sid, fullPath)
		So(err, ShouldBeNil)

		So(objectId3, ShouldEqual, fi.ObjectId)

		So(len(children), ShouldEqual, totalFiles3)

		/////////////////
		// test the directory 'mtp-test-files/mock_dir3/'
		/////////////////
		fullPath = "mtp-test-files/mock_dir3/"
		children = []*FileInfo{}

		objectId4, totalFiles4, err := Walk(dev, sid, fullPath, false,
			func(objectId uint32, fi *FileInfo, err error) error {
				So(err, ShouldBeNil)

				// make sure that the first item is not the parent path itself
				So(fi.FullPath, ShouldNotEqual, "/mtp-test-files/mock_dir3")
				So(fi.FullPath, ShouldContainSubstring, "/mtp-test-files/mock_dir3/")
				So(objectId, ShouldBeGreaterThan, 0)
				So(fi, ShouldNotBeNil)
				So(fi.ParentId, ShouldBeGreaterThan, 0)
				So(objectId, ShouldEqual, fi.ObjectId)

				children = append(children, fi)

				return nil
			})

		So(err, ShouldBeNil)
		So(totalFiles4, ShouldEqual, 5)

		// test if [objectId4] == [objectId] of [fullPath]
		fi, err = GetObjectFromPath(dev, sid, fullPath)
		So(err, ShouldBeNil)

		So(objectId4, ShouldEqual, fi.ObjectId)

		So(len(children), ShouldEqual, 5)
	})

	Convey("Testing valid directory | 1 | recursive=false | Walk", t, func() {
		//test the directory '/mtp-test-files/mock_dir1/1'
		fullPath := "/mtp-test-files/mock_dir1/1"

		var children []*FileInfo
		objectId, totalFiles, err := Walk(dev, sid, fullPath, false,
			func(objectId uint32, fi *FileInfo, err error) error {
				So(err, ShouldBeNil)

				// make sure that the first item is not the parent path itself
				So(fi.FullPath, ShouldNotEqual, "/mtp-test-files/mock_dir1/1")
				So(fi.FullPath, ShouldContainSubstring, "/mtp-test-files/mock_dir1/1/")
				So(objectId, ShouldBeGreaterThan, 0)
				So(fi, ShouldNotBeNil)
				So(fi.ParentId, ShouldBeGreaterThan, 0)
				So(objectId, ShouldEqual, fi.ObjectId)

				children = append(children, fi)

				return nil
			})

		So(err, ShouldBeNil)

		So(children, ShouldNotBeNil)
		So(len(children), ShouldEqual, totalFiles)
		So(len(children), ShouldEqual, 1)

		_file0 := children[0]

		So(_file0.ObjectId, ShouldBeGreaterThan, 0)
		So(_file0.Name, ShouldEqual, "a.txt")
		So(_file0.ParentId, ShouldEqual, objectId)
		So(_file0.Info.Filename, ShouldEqual, "a.txt")
		So(_file0.Extension, ShouldEqual, "txt")
		So(_file0.Size, ShouldEqual, 8)
		So(_file0.IsDir, ShouldEqual, false)
		So(_file0.FullPath, ShouldEqual, "/mtp-test-files/mock_dir1/1/a.txt")
		So(_file0.ModTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)
	})

	Convey("Testing valid directory | 2 | recursive=false | Walk", t, func() {
		//test the directory '/mtp-test-files/mock_dir1/'
		fullPath := "/mtp-test-files/mock_dir1/"

		var children []*FileInfo
		objectId, totalFiles, err := Walk(dev, sid, fullPath, false,
			func(objectId uint32, fi *FileInfo, err error) error {
				So(err, ShouldBeNil)

				// make sure that the first item is not the parent path itself
				So(fi.FullPath, ShouldNotEqual, "/mtp-test-files/mock_dir1")
				So(fi.FullPath, ShouldContainSubstring, "/mtp-test-files/mock_dir1/")
				So(objectId, ShouldBeGreaterThan, 0)
				So(fi, ShouldNotBeNil)
				So(fi.ParentId, ShouldBeGreaterThan, 0)
				So(objectId, ShouldEqual, fi.ObjectId)

				children = append(children, fi)

				return nil
			})

		So(err, ShouldBeNil)

		So(children, ShouldNotBeNil)
		So(len(children), ShouldEqual, totalFiles)
		So(len(children), ShouldEqual, 4)

		_file0 := children[0]

		So(_file0.ObjectId, ShouldBeGreaterThan, 0)
		So(_file0.Name, ShouldEqual, "1")
		So(_file0.ParentId, ShouldEqual, objectId)
		So(_file0.Info.Filename, ShouldEqual, "1")
		So(_file0.Extension, ShouldEqual, "")
		So(_file0.Size, ShouldBeGreaterThanOrEqualTo, 0)
		So(_file0.IsDir, ShouldEqual, true)
		So(_file0.FullPath, ShouldEqual, "/mtp-test-files/mock_dir1/1")
		So(_file0.ModTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)
	})

	Convey("Testing valid directory | 1 | recursive=true | Walk", t, func() {
		//test the directory '/mtp-test-files/mock_dir1/'
		fullPath := "/mtp-test-files/mock_dir1/"

		var children []*FileInfo
		objectId, totalFiles, err := Walk(dev, sid, fullPath, true,
			func(objectId uint32, fi *FileInfo, err error) error {
				So(err, ShouldBeNil)

				// make sure that the first item is not the parent path itself
				So(fi.FullPath, ShouldNotEqual, "/mtp-test-files/mock_dir1")
				So(fi.FullPath, ShouldContainSubstring, "/mtp-test-files/mock_dir1/")

				So(objectId, ShouldBeGreaterThan, 0)
				So(fi, ShouldNotBeNil)
				So(fi.ParentId, ShouldBeGreaterThan, 0)
				So(objectId, ShouldEqual, fi.ObjectId)

				children = append(children, fi)

				return nil
			})

		So(err, ShouldBeNil)

		const childrenLength = 9

		So(children, ShouldNotBeNil)
		So(len(children), ShouldEqual, childrenLength)
		So(totalFiles, ShouldEqual, 9)

		_file0 := children[0]

		So(_file0.ObjectId, ShouldBeGreaterThan, 0)
		So(_file0.Name, ShouldEqual, "1")
		So(_file0.ParentId, ShouldEqual, objectId)
		So(_file0.Info.Filename, ShouldEqual, "1")
		So(_file0.Extension, ShouldEqual, "")
		So(_file0.Size, ShouldBeGreaterThanOrEqualTo, 0)
		So(_file0.IsDir, ShouldEqual, true)
		So(_file0.FullPath, ShouldEqual, "/mtp-test-files/mock_dir1/1")
		So(_file0.ModTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)

		// test level 1 objects
		dirList1 := []string{"/mtp-test-files/mock_dir1/1", "/mtp-test-files/mock_dir1/1/a.txt", "/mtp-test-files/mock_dir1/a.txt", "/mtp-test-files/mock_dir1/3", "/mtp-test-files/mock_dir1/3/b.txt", "/mtp-test-files/mock_dir1/3/2", "/mtp-test-files/mock_dir1/3/2/b.txt", "/mtp-test-files/mock_dir1/2", "/mtp-test-files/mock_dir1/2/b.txt"}

		for i := range dirList1 {
			contains, index := StringContains(dirList1, children[i].FullPath)
			So(contains, ShouldEqual, true)
			dirList1 = RemoveIndex(dirList1, index)
		}
	})

	Convey("Testing valid file | recursive=true | Walk", t, func() {
		// test the directory '/mtp-test-files/a.txt'
		var children []*FileInfo
		objectId, totalFiles, err := Walk(dev, sid, "/mtp-test-files/a.txt", true,
			func(objectId uint32, fi *FileInfo, err error) error {
				So(err, ShouldBeNil)
				children = append(children, fi)
				So(fi.FullPath, ShouldEqual, "/mtp-test-files/a.txt")

				return nil
			})

		So(err, ShouldBeNil)
		So(objectId, ShouldBeGreaterThan, 0)
		So(totalFiles, ShouldEqual, 1)
		So(len(children), ShouldEqual, 1)
	})

	Convey("Testing non exisiting file | Walk | It should throw an error", t, func() {
		// test the directory '/fake' | recursive=true
		var children []*FileInfo
		objectId, totalFiles, err := Walk(dev, sid, "/fake", true,
			func(objectId uint32, fi *FileInfo, err error) error {
				So(err, ShouldBeNil)
				children = append(children, fi)

				return nil
			})

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objectId, ShouldEqual, 0)
		So(totalFiles, ShouldEqual, 0)
		So(len(children), ShouldEqual, 0)

		// test the directory '/fake' | recursive=false
		children = []*FileInfo{}
		objectId, totalFiles, err = Walk(dev, sid, "/fake", false,
			func(objectId uint32, fi *FileInfo, err error) error {
				So(err, ShouldBeNil)
				children = append(children, fi)

				return nil
			})

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objectId, ShouldEqual, 0)
		So(totalFiles, ShouldEqual, 0)
		So(len(children), ShouldEqual, 0)

		// test the directory '/mtp-test-files/fake' | recursive=true
		children = []*FileInfo{}
		objectId, totalFiles, err = Walk(dev, sid, "/mtp-test-files/fake", true,
			func(objectId uint32, fi *FileInfo, err error) error {
				So(err, ShouldBeNil)
				children = append(children, fi)

				return nil
			})

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objectId, ShouldEqual, 0)
		So(totalFiles, ShouldEqual, 0)
		So(len(children), ShouldEqual, 0)

		// test the directory '/mtp-test-files/fake' | recursive=false
		children = []*FileInfo{}
		objectId, totalFiles, err = Walk(dev, sid, "/mtp-test-files/fake", false,
			func(objectId uint32, fi *FileInfo, err error) error {
				So(err, ShouldBeNil)
				children = append(children, fi)

				return nil
			})

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objectId, ShouldEqual, 0)
		So(totalFiles, ShouldEqual, 0)
		So(len(children), ShouldEqual, 0)

		// test the directory=''
		objectId, totalFiles, err = Walk(dev, sid, "", true,
			func(objectId uint32, fi *FileInfo, err error) error {
				So(err, ShouldBeNil)
				children = append(children, fi)

				return nil
			})

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objectId, ShouldEqual, 0)
		So(totalFiles, ShouldEqual, 0)
		So(len(children), ShouldEqual, 0)
	})

	Convey("Testing callback error | Walk | It should throw an error", t, func() {
		// test the directory '/mtp-test-files' | recursive=true
		_, _, err := Walk(dev, sid, "/mtp-test-files", true,
			func(objectId uint32, fi *FileInfo, err error) error {
				So(err, ShouldBeNil)

				return InvalidPathError{error: fmt.Errorf("some error occured")}
			})

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
	})

	Dispose(dev)
}

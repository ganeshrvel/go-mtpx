package main

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"testing"
)

func TestWalkDirectory(t *testing.T) {
	dev, err := Initialize(Init{})
	if err != nil {
		log.Panic(err)
	}

	storages, err := FetchStorages(dev)
	if err != nil {
		log.Panic(err)
	}

	sid := storages[0].sid

	Convey("Testing valid directory | with objectId | objectId should be picked up instead of fullPath | WalkDirectory", t, func() {
		// test the root directory [ParentObjectId] | empty [fullPath]
		objectId, totalFiles, err := WalkDirectory(dev, sid, ParentObjectId, "", false, func(objectId uint32, fi *FileInfo) error {
			So(objectId, ShouldBeGreaterThan, 0)
			So(fi, ShouldNotBeNil)
			So(fi.ParentId, ShouldEqual, 0)

			return nil
		})

		So(err, ShouldBeNil)

		So(totalFiles, ShouldBeGreaterThan, 0)
		So(objectId, ShouldEqual, ParentObjectId)

		// test the root directory [ParentObjectId] | [fullPath]='/fake'
		objectId, totalFiles, err = WalkDirectory(dev, sid, ParentObjectId, "/fake", false, func(objectId uint32, fi *FileInfo) error {
			So(objectId, ShouldBeGreaterThan, 0)
			So(fi, ShouldNotBeNil)
			So(fi.ParentId, ShouldEqual, 0)

			return nil
		})

		So(err, ShouldBeNil)
		So(totalFiles, ShouldBeGreaterThan, 0)
		So(objectId, ShouldEqual, ParentObjectId)
	})

	Convey("Testing valid directory | without objectId | fullPath should be picked up instead of objectId | WalkDirectory", t, func() {
		/////////////////
		// test the directory '/mtp-test-files'
		/////////////////
		fullPath := "/mtp-test-files"

		var children []*FileInfo
		objectId1, totalFiles1, err := WalkDirectory(dev, sid, 0, fullPath, false, func(objectId uint32, fi *FileInfo) error {
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
		objectId2, totalFiles2, err := WalkDirectory(dev, sid, 0, fullPath, false, func(objectId uint32, fi *FileInfo) error {
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
		objectId3, totalFiles3, err := WalkDirectory(dev, sid, 0, fullPath, false, func(objectId uint32, fi *FileInfo) error {
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

		objectId4, totalFiles4, err := WalkDirectory(dev, sid, 0, fullPath, false, func(objectId uint32, fi *FileInfo) error {
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
		So(totalFiles4, ShouldBeGreaterThanOrEqualTo, totalFiles1)

		// test if [objectId4] == [objectId] of [fullPath]
		fi, err = GetObjectFromPath(dev, sid, fullPath)
		So(err, ShouldBeNil)

		So(objectId4, ShouldEqual, fi.ObjectId)

		So(len(children), ShouldEqual, 5)
	})

	Convey("Testing valid directory | 1 | recursive=false | WalkDirectory", t, func() {
		//test the directory '/mtp-test-files/mock_dir1/1'
		fullPath := "/mtp-test-files/mock_dir1/1"

		var children []*FileInfo
		objectId, totalFiles, err := WalkDirectory(dev, sid, 0, fullPath, false, func(objectId uint32, fi *FileInfo) error {
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

	Convey("Testing valid directory | 2 | recursive=false | WalkDirectory", t, func() {
		//test the directory '/mtp-test-files/mock_dir1/'
		fullPath := "/mtp-test-files/mock_dir1/"

		var children []*FileInfo
		objectId, totalFiles, err := WalkDirectory(dev, sid, 0, fullPath, false, func(objectId uint32, fi *FileInfo) error {
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
		So(_file0.Size, ShouldEqual, 4096)
		So(_file0.IsDir, ShouldEqual, true)
		So(_file0.FullPath, ShouldEqual, "/mtp-test-files/mock_dir1/1")
		So(_file0.ModTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)
	})

	Convey("Testing valid directory | 1 | recursive=true | WalkDirectory", t, func() {
		//test the directory '/mtp-test-files/mock_dir1/'
		fullPath := "/mtp-test-files/mock_dir1/"

		var children []*FileInfo
		objectId, totalFiles, err := WalkDirectory(dev, sid, 0, fullPath, true, func(objectId uint32, fi *FileInfo) error {
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
		So(_file0.Size, ShouldEqual, 4096)
		So(_file0.IsDir, ShouldEqual, true)
		So(_file0.FullPath, ShouldEqual, "/mtp-test-files/mock_dir1/1")
		So(_file0.ModTime.Year(), ShouldBeGreaterThanOrEqualTo, 2020)

		// test level 1 objects
		dirList1 := [childrenLength]string{"/mtp-test-files/mock_dir1/1", "/mtp-test-files/mock_dir1/1/a.txt", "/mtp-test-files/mock_dir1/a.txt", "/mtp-test-files/mock_dir1/3", "/mtp-test-files/mock_dir1/3/b.txt", "/mtp-test-files/mock_dir1/3/2", "/mtp-test-files/mock_dir1/3/2/b.txt", "/mtp-test-files/mock_dir1/2", "/mtp-test-files/mock_dir1/2/b.txt"}

		for i, _dir := range dirList1 {
			So(
				children[i].FullPath, ShouldEqual, _dir,
			)
		}
	})

	Convey("Testing non exisiting file | WalkDirectory | It should throw an error", t, func() {
		// test the directory '/fake' | recursive=true
		var children []*FileInfo
		objectId, totalFiles, err := WalkDirectory(dev, sid, 0, "/fake", true, func(objectId uint32, fi *FileInfo) error {
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
		objectId, totalFiles, err = WalkDirectory(dev, sid, 0, "/fake", false, func(objectId uint32, fi *FileInfo) error {
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
		objectId, totalFiles, err = WalkDirectory(dev, sid, 0, "/mtp-test-files/fake", true, func(objectId uint32, fi *FileInfo) error {
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
		objectId, totalFiles, err = WalkDirectory(dev, sid, 0, "/mtp-test-files/fake", false, func(objectId uint32, fi *FileInfo) error {
			children = append(children, fi)

			return nil
		})

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objectId, ShouldEqual, 0)
		So(totalFiles, ShouldEqual, 0)
		So(len(children), ShouldEqual, 0)

		// test the directory '/mtp-test-files/a.txt'
		children = []*FileInfo{}
		objectId, totalFiles, err = WalkDirectory(dev, sid, 0, "/mtp-test-files/a.txt", true, func(objectId uint32, fi *FileInfo) error {
			children = append(children, fi)

			return nil
		})

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objectId, ShouldEqual, 0)
		So(totalFiles, ShouldEqual, 0)
		So(len(children), ShouldEqual, 0)

		// test the directory=''
		objectId, totalFiles, err = WalkDirectory(dev, sid, 0, "", true, func(objectId uint32, fi *FileInfo) error {
			children = append(children, fi)

			return nil
		})

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objectId, ShouldEqual, 0)
		So(totalFiles, ShouldEqual, 0)
		So(len(children), ShouldEqual, 0)
	})

	Convey("Testing callback error | WalkDirectory | It should throw an error", t, func() {
		// test the directory '/mtp-test-files' | recursive=true
		_, _, err := WalkDirectory(dev, sid, 0, "/mtp-test-files", true, func(objectId uint32, fi *FileInfo) error {

			return InvalidPathError{error: fmt.Errorf("some error occured")}
		})

		So(err, ShouldBeError)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
	})

	Dispose(dev)
}

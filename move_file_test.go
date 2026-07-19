package mtpx

import (
	"fmt"
	"log"
	"math/rand"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// Note: like the other tests in this package, TestMoveFile talks to a real
// attached MTP device. Run these with a phone plugged in and the top-level
// test fixture directory (/mtp-test-files/temp_dir) present on it.
func TestMoveFile(t *testing.T) {
	dev, err := Initialize(Init{})
	if err != nil {
		log.Panic(err)
	}

	storages, err := FetchStorages(dev)
	if err != nil {
		log.Panic(err)
	}

	sid := storages[0].Sid

	Convey("Move an existing directory into a different parent | using fullPath | MoveFile", t, func() {
		srcPath := fmt.Sprintf("/mtp-test-files/temp_dir/test-MoveFile/src-%x", rand.Int31())
		destParent := fmt.Sprintf("/mtp-test-files/temp_dir/test-MoveFile/dst-%x", rand.Int31())

		srcObjectId, err := MakeDirectory(dev, sid, srcPath)
		So(err, ShouldBeNil)
		So(srcObjectId, ShouldBeGreaterThan, 0)

		// The destination parent will be created on demand by MoveFile via
		// MakeDirectory — no need to pre-create it here.
		movedObjectId, err := MoveFile(dev, sid, FileProp{FullPath: srcPath}, destParent)
		So(err, ShouldBeNil)
		So(movedObjectId, ShouldEqual, srcObjectId)

		// A second move into the same destination is a no-op — the source is
		// already parented there — and should return the same objectId.
		newSrcPath := getFullPath(destParent, filepath_Base(srcPath))
		sameObjectId, err := MoveFile(dev, sid, FileProp{FullPath: newSrcPath}, destParent)
		So(err, ShouldBeNil)
		So(sameObjectId, ShouldEqual, srcObjectId)
	})

	Convey("Move onto an existing destination name | MoveFile | Should throw FileAlreadyExistsError", t, func() {
		srcPath := fmt.Sprintf("/mtp-test-files/temp_dir/test-MoveFile/src-%x", rand.Int31())
		destParent := fmt.Sprintf("/mtp-test-files/temp_dir/test-MoveFile/dst-%x", rand.Int31())

		_, err := MakeDirectory(dev, sid, srcPath)
		So(err, ShouldBeNil)

		// Pre-seed the collision at destParent/<basename(srcPath)>.
		collidingPath := getFullPath(destParent, filepath_Base(srcPath))
		_, err = MakeDirectory(dev, sid, collidingPath)
		So(err, ShouldBeNil)

		objId, err := MoveFile(dev, sid, FileProp{FullPath: srcPath}, destParent)
		So(err, ShouldHaveSameTypeAs, FileAlreadyExistsError{})
		So(objId, ShouldEqual, 0)
	})

	Convey("Move a non-existent object | MoveFile | Should throw InvalidPathError", t, func() {
		srcPath := fmt.Sprintf("/mtp-test-files/temp_dir/test-MoveFile/missing-%x", rand.Int31())
		destParent := "/mtp-test-files/temp_dir/test-MoveFile"

		objId, err := MoveFile(dev, sid, FileProp{FullPath: srcPath}, destParent)
		So(err, ShouldHaveSameTypeAs, InvalidPathError{})
		So(objId, ShouldEqual, 0)
	})

	Dispose(dev)
}

// filepath_Base is imported inline to avoid pulling path/filepath into every
// test file. Kept local so the test reads independently.
func filepath_Base(p string) string {
	for i := len(p) - 1; i >= 0; i-- {
		if p[i] == '/' {
			return p[i+1:]
		}
	}
	return p
}

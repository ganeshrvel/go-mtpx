package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestUtils(t *testing.T) {
	Convey("Test extensions", t, func() {
		filenameList := []string{"abc.txt", "xyz.gz", "123", "123.tar.gz", ".ssh", ".gitignore", "github.com/ganeshrvel/one-archiver/e2e_list_test.go", "one-archiver/e2e_list_test.go", "e2e_list_test.go/.go.psd"}
		extList := []string{"txt", "gz", "", "tar.gz", "ssh", "gitignore", "go", "go", "go.psd"}

		for i, f := range filenameList {
			ext := extension(f, false)

			So(extList[i], ShouldEqual, ext)
		}
	})

	Convey("Test fixDirSlash", t, func() {
		filenameList := []string{"", "/", "//", "/xyz", "//xyz", "/xyz/", "/xyz//", "xyz/", "xyz", "xyz/124", "xyz/124/", "/xyz/124/"}
		dirList := []string{"/", "/", "/", "/xyz", "//xyz", "/xyz", "/xyz/", "/xyz", "/xyz", "/xyz/124", "/xyz/124", "/xyz/124"}

		for i, f := range filenameList {
			dir := fixDirSlash(f)

			So(dirList[i], ShouldEqual, dir)
		}
	})

	Convey("Test getFullPath", t, func() {
		type s struct {
			parentPath, filename, fullPath string
		}

		sl := []s{
			s{
				parentPath: "/",
				filename:   "abc",
				fullPath:   "/abc",
			},
			s{
				parentPath: "//",
				filename:   "abc",
				fullPath:   "///abc",
			},
			s{
				parentPath: "/",
				filename:   "abc/",
				fullPath:   "/abc/",
			},
			s{
				parentPath: "/xyz",
				filename:   "abc/",
				fullPath:   "/xyz/abc/",
			},
		}

		for i, f := range sl {
			fullPath := getFullPath(f.parentPath, f.filename)

			So(sl[i].fullPath, ShouldEqual, fullPath)
		}
	})
}

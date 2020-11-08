package mtpx

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestUtils(t *testing.T) {
	//if testing.Short() {
	//	t.Skip("skipping 'TestUtils' testing in short mode")
	//}

	Convey("Test extensions", t, func() {
		filenameList := []string{"abc.txt", "xyz.gz", "123", "123.tar.gz", ".ssh", ".gitignore", "github.com/ganeshrvel/one-archiver/e2e_list_test.go", "one-archiver/e2e_list_test.go", "e2e_list_test.go/.go.psd"}
		extList := []string{"txt", "gz", "", "tar.gz", "ssh", "gitignore", "go", "go", "go.psd"}

		for i, f := range filenameList {
			ext := extension(f, false)

			So(extList[i], ShouldEqual, ext)
		}
	})

	Convey("Test fixSlash", t, func() {
		filenameList := []string{"", ".", "/./", "././", "/../", "/", "//", "/abc", "//bcd", "/cde/", "/def//", "efg/", "fgh", "ghi/124", "hij/124/", "/ijk/124/"}
		dirList := []string{"/", "/", "/", "/", "/", "/", "/", "/abc", "/bcd", "/cde", "/def", "/efg", "/fgh", "/ghi/124", "/hij/124", "/ijk/124"}

		for i, f := range filenameList {
			dir := fixSlash(f)

			So(dir, ShouldEqual, dirList[i])
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
				filename:   "bcd",
				fullPath:   "/bcd",
			},
			s{
				parentPath: "/",
				filename:   "cde/",
				fullPath:   "/cde",
			},
			s{
				parentPath: "/def",
				filename:   "abc/",
				fullPath:   "/def/abc",
			},
			s{
				parentPath: "/efg/",
				filename:   "abc/",
				fullPath:   "/efg/abc",
			},
		}

		for i, f := range sl {
			fullPath := getFullPath(f.parentPath, f.filename)

			So(fullPath, ShouldEqual, sl[i].fullPath)
		}
	})
}

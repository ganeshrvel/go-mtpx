package mtpx

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestUtils(t *testing.T) {
	//if testing.Short() {
	//	t.Skip("skipping 'TestUtils' testing in short mode")
	//}

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

	Convey("Test extension", t, func() {
		type s struct {
			filename, ext string
			isDir         bool
		}

		sl := []s{
			s{
				filename: "",
				ext:      "",
				isDir:    false,
			}, s{
				filename: "abc.xyz.tar.gz",
				ext:      "tar.gz",
				isDir:    false,
			}, s{
				filename: "abc.xyz.tar.tar",
				ext:      "tar.tar",
				isDir:    false,
			}, s{
				filename: "xyz.tar.gz",
				ext:      "tar.gz",
				isDir:    false,
			}, s{
				filename: "tar.gz",
				ext:      "gz",
				isDir:    false,
			}, s{
				filename: "abc.gz",
				ext:      "gz",
				isDir:    false,
			}, s{
				filename: ".gz",
				ext:      "gz",
				isDir:    false,
			}, s{
				filename: ".tar",
				ext:      "tar",
				isDir:    false,
			}, s{
				filename: ".tar.gz",
				ext:      "tar.gz",
				isDir:    false,
			}, s{
				filename: "tar.tar.gz",
				ext:      "tar.gz",
				isDir:    false,
			}, s{
				filename: ".htaccess",
				ext:      "htaccess",
				isDir:    false,
			}, s{
				filename: "abc.txt",
				ext:      "txt",
				isDir:    false,
			}, s{
				filename: "abc",
				ext:      "",
				isDir:    false,
			}, s{
				filename: "github.com/ganeshrvel/one-archiver/e2e_list_test.go",
				ext:      "go",
				isDir:    false,
			}, s{
				filename: "one-archiver/e2e_list_test.go",
				ext:      "go",
				isDir:    false,
			}, s{
				filename: "e2e_list_test.go/.go.psd",
				ext:      "psd",
				isDir:    false,
			}, s{
				filename: "abc",
				ext:      "",
				isDir:    true,
			}, s{
				filename: "abc.tar",
				ext:      "",
				isDir:    true,
			}, s{
				filename: "abc.tar.gz",
				ext:      "",
				isDir:    true,
			},
		}

		for _, f := range sl {
			ext := extension(f.filename, f.isDir)

			So(ext, ShouldEqual, f.ext)
		}
	})
}

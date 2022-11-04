module github.com/ganeshrvel/go-mtpx

go 1.15

require (
	github.com/ganeshrvel/go-mtpfs v1.0.4-0.20221104074511-0d40588840c5
	github.com/smartystreets/goconvey v1.6.4
	golang.org/x/sys v0.0.0-20201231184435-2d18734c6014 // indirect
)

// replace github.com/ganeshrvel/go-mtpfs v1.0.4-0.20201206195153-a90fac923f97 => ../go-mtpfs
///##### Upgrade a package
//go get github.com/<org-name>/<package-name>@<git-commit-hash>
//example: go get github.com/ganeshrvel/go-mtpfs@<git-commit-hash>

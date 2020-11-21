###NOTE
### Build wouldn't work as this package do not have a main function

# build the binary
CGO_CFLAGS='-Wno-deprecated-declarations' go build -gcflags=-trimpath=$(go env GOPATH) -asmflags=-trimpath=$(go env GOPATH) -o build/mtpx .

# copy libusb
cp lib/libusb-1.0.0.dylib build/libusb-1.0.0.dylib

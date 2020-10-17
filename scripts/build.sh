# build the binary
CGO_CFLAGS='-Wno-deprecated-declarations' go build -o build/mtpx .

# copy libusb
cp lib/libusb-1.0.0.dylib build/libusb-1.0.0.dylib

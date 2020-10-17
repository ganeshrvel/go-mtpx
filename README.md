go-mtpx



macOS setup

```shell script
xcode-select --install
```

Install and setup libusb

Method 1:
```shell script
# build the binary
CGO_CFLAGS='-Wno-deprecated-declarations' go build -o build/mtpx .

# copy libusb
cp lib/libusb-1.0.0.dylib build/libusb-1.0.0.dylib

# Run
DYLD_LIBRARY_PATH=./build ./build/mtpx
```

Method 2:
```shell script
brew install libusb

# find the location where libusb is installed
brew info libusb

# run: sudo install_name_tool -id "@executable_path/libusb.dylib" /usr/local/Cellar/libusb/<version>/lib/libusb-1.0.0.dylib

# eg: /usr/local/Cellar/libusb/1.0.23
sudo install_name_tool -id "@executable_path/libusb.dylib" /usr/local/Cellar/libusb/1.0.23/lib/libusb-1.0.0.dylib

# copy the dynamic library to the project
cp /usr/local/Cellar/libusb/1.0.23/lib/libusb-1.0.0.dylib ./libusb.dylib

# confirm whether the libusb.dylib is portable
otool -L libusb.dylib

# the output should look: libusb.dylib: @executable_path/libusb.dylib 

git add libusb.dylib
```

Run mtpx
```shell script
go run ./
```

Build mtpx
```shell script
CGO_CFLAGS='-Wno-deprecated-declarations' go build -o build/mtpx . && cp libusb.dylib build/libusb.dylib

./build/mtpx
```


Samsung test for MTPx

Expected feedback
- Please send me any errors that you may encounter while executing the code. This will help me understand the portability of the code.
- If the execution was successful, then please send me a redacted console output


macOS setup:
```shell script
xcode-select --install
```

- Install GoLang from: https://golang.org/doc/install

Brew setup:
```shell script
# install brew
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)"

# update
brew update

# install libusb
brew install libusb
```

Git clone:
```shell script
git clone https://github.com/ganeshrvel/go-mtpx
cd go-mtpx
git fetch
git checkout test/samsung
```

Build:
```shell script
# build the binary
./scripts/build.sh

# Run
DYLD_LIBRARY_PATH=./build ./build/mtpx
```


If the above steps didn't work then:
```shell script
# build the binary
./scripts/build.sh

# Run
./build/mtpx
```


Cheers!

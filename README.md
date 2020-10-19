macOS setup

```shell script
xcode-select --install
```

Install GoLang from: https://golang.org/doc/install


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
# install brew
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)"

# brew update
brew update

# brew install libusb

./scripts/build.sh

./build
```

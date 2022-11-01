# go-mtpx

## macOS setup

```shell script
xcode-select --install
```

### Before running the tests:

- Open `tests/README.md`
- Follow the instructions to copy the `mtp-test-files` to phone


### Install and setup libusb

```shell
brew install pkg-config
brew install libusb
```


### Test go-mtpx

- Connect android device via usb and Choose File transfer

```shell script
go test
```

##### Upgrade a package

```shell
# example:
go get github.com/ganeshrvel/go-mtpfs@<git-commit-hash>
```
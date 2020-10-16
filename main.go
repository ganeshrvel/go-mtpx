package main

import (
	mtp "github.com/ganeshrvel/go-mtpfs/mtp"
)

// initialize the mtp device
// returns mtp device
func Initialize() (*mtp.Device, error) {
	dev, err := mtp.SelectDevice("")

	if err != nil {
		return nil, MtpDetectFailedError{error: err}
	}

	dev.MTPDebug = isDev
	dev.DataDebug = isDev
	dev.USBDebug = isDev

	dev.Timeout = devTimeout

	if err = dev.Configure(); err != nil {
		return nil, ConfigureError{error: err}
	}

	return dev, nil
}

// close the mtp device
func Dispose(dev *mtp.Device) {
	defer dev.Close()
}

// fetch device info
func FetchDeviceInfo(dev *mtp.Device) (*mtp.DeviceInfo, error) {
	info := mtp.DeviceInfo{}
	err := dev.GetDeviceInfo(&info)

	if err != nil {
		return nil, DeviceInfoError{error: err}
	}

	return &info, nil
}

// fetch storages
func FetchStorages(dev *mtp.Device) (*[]mtp.StorageInfo, error) {
	sids := mtp.Uint32Array{}
	if err := dev.GetStorageIDs(&sids); err != nil {
		return nil, StorageInfoError{error: err}
	}

	if len(sids.Values) < 1 {
		return nil, NoStorageError{}
	}

	var result []mtp.StorageInfo

	for sid := range sids.Values {
		var info mtp.StorageInfo
		if err := dev.GetStorageInfo(uint32(sid), &info); err != nil {
			return nil, StorageInfoError{error: err}
		}

		result = append(result, info)
	}

	return &result, nil
}

func main() {
	dev, _ := Initialize()

	//di, _ := FetchDeviceInfo(dev)

	//si, _ := FetchStorageInfo(dev)
	//
	//pretty.Println("======\n")
	//pretty.Println(si)

	Dispose(dev)
}

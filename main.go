package main

import (
	mtp "github.com/ganeshrvel/go-mtpfs/mtp"
	"github.com/kr/pretty"
	"log"
)

// initialize the mtp device
// returns mtp device
func Initialize(init Init) (*mtp.Device, error) {
	dev, err := mtp.SelectDevice("")

	if err != nil {
		return nil, MtpDetectFailedError{error: err}
	}

	dev.MTPDebug = init.debugMode
	dev.DataDebug = init.debugMode
	dev.USBDebug = init.debugMode

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
func FetchStorages(dev *mtp.Device) ([]StorageData, error) {
	sids := mtp.Uint32Array{}
	if err := dev.GetStorageIDs(&sids); err != nil {
		return nil, StorageInfoError{error: err}
	}

	if len(sids.Values) < 1 {
		return nil, NoStorageError{}
	}

	var result []StorageData

	for sid := range sids.Values {
		var info mtp.StorageInfo
		if err := dev.GetStorageInfo(uint32(sid), &info); err != nil {
			return nil, StorageInfoError{error: err}
		}

		result = append(result, StorageData{
			sid:  sids.Values[0],
			info: info,
		})
	}

	return result, nil
}

func ListDirectory(dev *mtp.Device, storageId uint32) (*[]MtpFileInfo, error) {
	handles := mtp.Uint32Array{}
	if err := dev.GetObjectHandles(storageId, mtp.GOH_ALL_ASSOCS, mtp.GOH_ROOT_PARENT, &handles); err != nil {
		return nil, ListDirectoryError{error: err}
	}

	var fileInfoList []MtpFileInfo

	for _, handle := range handles.Values {
		obj := mtp.ObjectInfo{}
		if err := dev.GetObjectInfo(handle, &obj); err != nil {
			log.Printf("GetObjectInfo for handle %d failed: %v", handle, err)

			continue
		}

		if obj.Filename == "" {
			continue
		}

		var size int64
		if obj.CompressedSize == 0xffffffff {
			var val mtp.Uint64Value
			if err := dev.GetObjectPropValue(handle, mtp.OPC_ObjectSize, &val); err != nil {
				log.Printf("GetObjectPropValue handle %d failed: %v", handle, err)

				continue
			}

			size = int64(val.Value)
		} else {
			size = int64(obj.CompressedSize)
		}

		isDir := isDirectoryObject(&obj)

		fi := MtpFileInfo{
			Info:         &obj,
			Size:         size,
			IsDir:        isDir,
			ModTime:      obj.ModificationDate,
			Name:         obj.Filename,
			FullPath:     "",
			ParentPath:   "",
			Extension:    extension(obj.Filename, isDir),
			ParentObject: obj.ParentObject,
		}

		fileInfoList = append(fileInfoList, fi)
	}

	return &fileInfoList, nil
}

func main() {
	dev, err := Initialize(Init{})

	if err != nil {
		log.Panic(err)
	}

	_, err = FetchDeviceInfo(dev)
	if err != nil {
		log.Panic(err)
	}

	storages, err := FetchStorages(dev)
	if err != nil {
		log.Panic(err)
	}

	files, err := ListDirectory(dev, storages[0].sid)
	if err != nil {
		log.Panic(err)
	}

	pretty.Println("======\n")
	pretty.Println(files)

	Dispose(dev)
}

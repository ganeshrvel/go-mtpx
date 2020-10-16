package main

type MtpDetectFailedError struct {
	error
}

type ConfigureError struct {
	error
}

type DeviceInfoError struct {
	error
}

type StorageInfoError struct {
	error
}

type NoStorageError struct {
	error
}

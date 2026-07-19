package mtpx

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

type ListDirectoryError struct {
	error
}

type FileNotFoundError struct {
	error
}

type FilePermissionError struct {
	error
}

type LocalFileError struct {
	error
}

type InvalidPathError struct {
	error
}

type FileTransferError struct {
	error
}

type FileObjectError struct {
	error
}

type SendObjectError struct {
	error
}

// FileAlreadyExistsError is returned by MoveFile when the destination parent
// directory already contains an entry with the same base name as the source.
type FileAlreadyExistsError struct {
	error
}

// MoveNotSupportedError is returned by MoveFile when the device rejects the
// MTP MoveObject operation with Operation_Not_Supported (0x2005). Callers may
// fall back to copy+delete in this case.
type MoveNotSupportedError struct {
	error
}

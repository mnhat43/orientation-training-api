package cloud

import "io"

// StorageUtility interface
type StorageUtility interface {
	GetFileByFileName(fileName string, directoryCloud string) string
	UploadFileToCloud(file io.Reader, fileName string, directoryCloud string) error
	DeleteFileCloud(fileName string, directoryCloud string) error
}

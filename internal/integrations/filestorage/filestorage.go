package filestorage

import (
	"fmt"
)

var FileDoesNotExist error = fmt.Errorf("the specified file does not exist")

type FileStorageManager interface {
	WriteFile(path string, bytes []byte, shouldEncrypt bool) error
	ReadFile(path string, shouldDecrypt bool) ([]byte, error)
	DeleteFile(path string) error
}

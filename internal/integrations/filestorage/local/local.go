package local

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/hatchet-dev/hatchet/internal/encryption"
	"github.com/hatchet-dev/hatchet/internal/integrations/filestorage"
)

type LocalFileStorageManager struct {
	rootDir       string
	encryptionKey *[32]byte
}

func NewLocalFileStorageManager(rootDir string, encryptionKey *[32]byte) (*LocalFileStorageManager, error) {
	err := os.MkdirAll(rootDir, os.ModePerm)

	if err != nil {
		return nil, fmt.Errorf("error creating file directory: %w", err)
	}

	return &LocalFileStorageManager{rootDir, encryptionKey}, nil
}

func (l *LocalFileStorageManager) WriteFile(path string, fileBytes []byte, shouldEncrypt bool) error {
	body := fileBytes
	var err error

	if shouldEncrypt {
		body, err = encryption.Encrypt(fileBytes, l.encryptionKey)

		if err != nil {
			return err
		}
	}

	fullFilePath := filepath.Join(l.rootDir, path)

	return ioutil.WriteFile(fullFilePath, body, 0666)
}

func (l *LocalFileStorageManager) ReadFile(path string, shouldDecrypt bool) ([]byte, error) {
	fullFilePath := filepath.Join(l.rootDir, path)

	fileBytes, err := os.ReadFile(fullFilePath)

	if err != nil {
		if os.IsNotExist(err) {
			return nil, filestorage.FileDoesNotExist
		}

		return nil, err
	}

	if shouldDecrypt {
		fileBytes, err = encryption.Decrypt(fileBytes, l.encryptionKey)

		if err != nil {
			return nil, err
		}
	}

	return fileBytes, nil
}

func (l *LocalFileStorageManager) DeleteFile(path string) error {
	fullFilePath := filepath.Join(l.rootDir, path)

	err := os.Remove(fullFilePath)

	if err != nil {
		if os.IsNotExist(err) {
			return filestorage.FileDoesNotExist
		}

		return err
	}

	return nil
}

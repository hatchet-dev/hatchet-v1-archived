package s3

import (
	"bytes"
	"context"
	"errors"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/hatchet-dev/hatchet/internal/encryption"
	storage "github.com/hatchet-dev/hatchet/internal/integrations/filestorage"
)

type S3StorageClient struct {
	client        *s3.Client
	bucket        string
	encryptionKey *[32]byte
}

type S3Options struct {
	AWSRegion      string
	AWSAccessKeyID string
	AWSSecretKey   string
	AWSBucketName  string
	EncryptionKey  *[32]byte
}

func NewS3StorageClient(opts *S3Options) (*S3StorageClient, error) {
	client := s3.New(s3.Options{
		Region:      opts.AWSRegion,
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(opts.AWSAccessKeyID, opts.AWSSecretKey, "")),
	})

	return &S3StorageClient{
		bucket:        opts.AWSBucketName,
		encryptionKey: opts.EncryptionKey,
		client:        client,
	}, nil
}

func (s *S3StorageClient) WriteFile(path string, fileBytes []byte, shouldEncrypt bool) error {
	body := fileBytes
	var err error
	if shouldEncrypt {
		body, err = encryption.Encrypt(fileBytes, s.encryptionKey)

		if err != nil {
			return err
		}
	}

	_, err = s.client.PutObject(
		context.Background(),
		&s3.PutObjectInput{
			Body:   manager.ReadSeekCloser(bytes.NewReader(body)),
			Bucket: &s.bucket,
			Key:    aws.String(path),
		})

	return err
}

func (s *S3StorageClient) ReadFile(path string, shouldDecrypt bool) ([]byte, error) {
	output, err := s.client.GetObject(
		context.Background(),
		&s3.GetObjectInput{
			Bucket: &s.bucket,
			Key:    aws.String(path),
		})

	if err != nil {
		var nosuchkey *types.NoSuchKey

		if errors.As(err, &nosuchkey) {
			return nil, storage.FileDoesNotExist
		}

		return nil, err
	}

	if shouldDecrypt {
		var encryptedData bytes.Buffer

		_, err = encryptedData.ReadFrom(output.Body)

		if err != nil {
			return nil, err
		}

		data, err := encryption.Decrypt(encryptedData.Bytes(), s.encryptionKey)

		if err != nil {
			return nil, err
		}

		return data, nil
	} else {
		return io.ReadAll(output.Body)
	}
}

func (s *S3StorageClient) DeleteFile(path string) error {
	_, err := s.client.DeleteObject(
		context.Background(),
		&s3.DeleteObjectInput{
			Bucket: &s.bucket,
			Key:    aws.String(path),
		})

	if err != nil {
		return err
	}

	return nil
}

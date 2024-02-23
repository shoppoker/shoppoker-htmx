package file_storage

import (
	"bytes"
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/w1png/go-htmx-ecommerce-template/config"
)

type ObjectStorageId string

func (o ObjectStorageId) GetUrl(type_, extension string) string {
	return fmt.Sprintf("/file/%s/%s/%s", type_, extension, o)
}

type FileStorage interface {
	UploadFile(file_bytes []byte) (ObjectStorageId, error)

	DoesFileExist(file_bytes []byte) (bool, error)
	DoesObjectIdExist(object_storage_id ObjectStorageId) (bool, error)

	DeleteFile(object_storage_id ObjectStorageId) error
	GetDownloadUrl(object_storage_id ObjectStorageId, filename string) (string, error)
	DownloadFile(object_storage_id ObjectStorageId) ([]byte, error)

	GenerateMD5Hash(file_bytes []byte) (string, error)
}

var FileStorageInstance FileStorage

func InitFileStorage() (err error) {
	FileStorageInstance, err = NewYandexObjecStorage()
	return err
}

type YandexObjectStorage struct {
	Client *s3.Client
}

func NewYandexObjecStorage() (*YandexObjectStorage, error) {
	cfg, err := aws_config.LoadDefaultConfig(context.TODO(), aws_config.WithEndpointResolverWithOptions(
		aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			if service == s3.ServiceID && region == "ru-central1" {
				return aws.Endpoint{
					PartitionID:   "yc",
					URL:           "https://storage.yandexcloud.net",
					SigningRegion: "ru-central1",
				}, nil
			}
			return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
		}),
	))
	if err != nil {
		return nil, err
	}

	return &YandexObjectStorage{
		Client: s3.NewFromConfig(cfg),
	}, nil
}

func (s *YandexObjectStorage) GenerateMD5Hash(file_bytes []byte) (string, error) {
	hasher := md5.New()

	if _, err := io.Copy(hasher, bytes.NewBuffer(file_bytes)); err != nil {
		return "", err
	}

	hasher.Write([]byte(strconv.FormatInt(time.Now().UnixNano(), 10) + strconv.Itoa(os.Getpid()) + strconv.Itoa(os.Getuid())))

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

func (s *YandexObjectStorage) UploadFile(file_bytes []byte) (ObjectStorageId, error) {
	md5_hash, err := s.GenerateMD5Hash(file_bytes)
	if err != nil {
		return "", err
	}

	if _, err := s.Client.PutObject(
		context.Background(),
		&s3.PutObjectInput{
			Bucket: aws.String(config.ConfigInstance.ObjectStorageBucketName),
			Key:    aws.String(md5_hash),
			Body:   bytes.NewBuffer(file_bytes),
		},
	); err != nil {
		return "", err
	}

	return ObjectStorageId(md5_hash), nil
}

func (s *YandexObjectStorage) DeleteFile(object_storage_id ObjectStorageId) error {
	_, err := s.Client.DeleteObject(
		context.Background(),
		&s3.DeleteObjectInput{
			Bucket: aws.String(config.ConfigInstance.ObjectStorageBucketName),
			Key:    aws.String(string(object_storage_id)),
		},
	)

	return err
}

func (s *YandexObjectStorage) GetDownloadUrl(object_storage_id ObjectStorageId, filename string) (string, error) {
	url, err := s3.NewPresignClient(s.Client).PresignGetObject(context.Background(), &s3.GetObjectInput{
		Bucket:                     aws.String(config.ConfigInstance.ObjectStorageBucketName),
		Key:                        aws.String(string(object_storage_id)),
		ResponseContentDisposition: aws.String(fmt.Sprintf("attachment; filename=\"%s\"", filename)),
	}, s3.WithPresignExpires(time.Hour*1))
	if err != nil {
		return "", err
	}

	return url.URL, nil
}

func (s *YandexObjectStorage) DownloadFile(object_storage_id ObjectStorageId) ([]byte, error) {
	downloader := manager.NewDownloader(s.Client)
	buf := manager.NewWriteAtBuffer([]byte{})

	if _, err := downloader.Download(context.TODO(), buf, &s3.GetObjectInput{
		Bucket: aws.String(config.ConfigInstance.ObjectStorageBucketName),
		Key:    aws.String(string(object_storage_id)),
	}); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

var noSuchKey *types.NoSuchKey
var notFound *types.NotFound

func (s *YandexObjectStorage) DoesFileExist(file_bytes []byte) (bool, error) {
	md5_hash, err := s.GenerateMD5Hash(file_bytes)
	if err != nil {
		return false, err
	}

	_, err = s.Client.HeadObject(
		context.Background(),
		&s3.HeadObjectInput{
			Bucket: aws.String(config.ConfigInstance.ObjectStorageBucketName),
			Key:    aws.String(md5_hash),
		},
	)

	if err != nil {
		if errors.As(err, &noSuchKey) || errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (s *YandexObjectStorage) DoesObjectIdExist(object_storage_id ObjectStorageId) (bool, error) {
	_, err := s.Client.HeadObject(
		context.Background(),
		&s3.HeadObjectInput{
			Bucket: aws.String(config.ConfigInstance.ObjectStorageBucketName),
			Key:    aws.String(string(object_storage_id)),
		},
	)

	if err != nil {
		if errors.As(err, &noSuchKey) || errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

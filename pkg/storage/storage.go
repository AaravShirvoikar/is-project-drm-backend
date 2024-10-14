package storage

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/AaravShirvoikar/is-project-drm-backend/internal/models"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type FileStorage struct {
	minioClient *minio.Client
	bucketName  string
}

func NewFileStorage(endpoint, accessKeyID, secretAccessKey, bucketName string, useSSL bool) (*FileStorage, error) {
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	exists, err := minioClient.BucketExists(context.Background(), bucketName)
	if err != nil {
		return nil, err
	}

	if !exists {
		err = minioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
		log.Printf("bucket %s created successfully\n", bucketName)
	} else {
		log.Printf("bucket %s already exists\n", bucketName)
	}

	return &FileStorage{minioClient: minioClient, bucketName: bucketName}, nil
}

func (s *FileStorage) UploadFile(ctx context.Context, reader io.Reader, ext string, size int64) (*models.File, error) {
	fileId, err := generateUniqueFilename(ext)
	if err != nil {
		return nil, err
	}

	_, err = s.minioClient.PutObject(ctx, s.bucketName, fileId, reader, size, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return nil, err
	}

	file := &models.File{
		FileID: fileId,
		Size:   size,
	}

	return file, nil
}

func (s *FileStorage) DownloadFile(ctx context.Context, fileId string) (io.ReadCloser, error) {
	return s.minioClient.GetObject(ctx, s.bucketName, fileId, minio.GetObjectOptions{})
}

func generateUniqueFilename(fileExtension string) (string, error) {
	timestamp := time.Now().Format("20060102-150405")

	randomBytes := make([]byte, 4)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	randomString := hex.EncodeToString(randomBytes)

	return fmt.Sprintf("%s-%s%s", timestamp, randomString, fileExtension), nil
}

package storage

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log"

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

func (s *FileStorage) UploadFile(ctx context.Context, reader io.Reader, size int64) (*models.File, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, reader); err != nil {
		return nil, err
	}
	fileHash := hex.EncodeToString(hash.Sum(nil))

	seeker, ok := reader.(io.Seeker)
	if ok {
		_, err := seeker.Seek(0, io.SeekStart)
		if err != nil {
			return nil, err
		}
	}

	_, err := s.minioClient.PutObject(ctx, s.bucketName, fileHash, reader, size, minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		return nil, err
	}

	file := &models.File{
		Hash: fileHash,
		Size: size,
	}

	return file, nil
}

func (s *FileStorage) DownloadFile(ctx context.Context, fileHash string) (io.ReadCloser, error) {
	return s.minioClient.GetObject(ctx, s.bucketName, fileHash, minio.GetObjectOptions{})
}

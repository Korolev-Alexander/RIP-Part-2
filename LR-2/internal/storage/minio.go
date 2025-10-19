package storage

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIOClient struct {
	client *minio.Client
	bucket string
}

func NewMinIOClient() *MinIOClient {
	// Создаем клиент MinIO
	minioClient, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("myaccesskey123", "mysecretkey123456", ""),
		Secure: false,
	})
	if err != nil {
		log.Printf("⚠️ Failed to create MinIO client: %v", err)
		return &MinIOClient{}
	}

	// Проверяем подключение
	_, err = minioClient.ListBuckets(context.Background())
	if err != nil {
		log.Printf("⚠️ MinIO connection failed: %v", err)
	} else {
		log.Printf("✅ MinIO client initialized successfully")
	}

	return &MinIOClient{
		client: minioClient,
		bucket: "image",
	}
}

func (m *MinIOClient) UploadFile(filename string, fileData []byte) error {
	if m.client == nil {
		return fmt.Errorf("MinIO client not initialized")
	}

	// Создаем bucket если не существует
	exists, err := m.client.BucketExists(context.Background(), m.bucket)
	if err != nil {
		return fmt.Errorf("failed to check bucket: %v", err)
	}

	if !exists {
		err = m.client.MakeBucket(context.Background(), m.bucket, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %v", err)
		}
		log.Printf("✅ Created bucket: %s", m.bucket)
	}

	// Загружаем файл
	_, err = m.client.PutObject(context.Background(), m.bucket, filename,
		bytes.NewReader(fileData), int64(len(fileData)),
		minio.PutObjectOptions{
			ContentType: "image/png",
		})

	if err != nil {
		return fmt.Errorf("failed to upload file: %v", err)
	}

	log.Printf("✅ File uploaded to MinIO: %s (%d bytes)", filename, len(fileData))
	return nil
}

func (m *MinIOClient) DeleteFile(filename string) error {
	if m.client == nil {
		return fmt.Errorf("MinIO client not initialized")
	}

	err := m.client.RemoveObject(context.Background(), m.bucket, filename,
		minio.RemoveObjectOptions{})

	if err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	log.Printf("✅ File deleted from MinIO: %s", filename)
	return nil
}

func (m *MinIOClient) GetImageURL(filename string) string {
	return fmt.Sprintf("http://localhost:9000/%s/%s", m.bucket, filename)
}

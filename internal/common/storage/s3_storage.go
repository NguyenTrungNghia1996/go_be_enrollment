package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/url"
	"path"
	"strings"
	"time"

	"go_be_enrollment/internal/config"
	"go_be_enrollment/pkg/logger"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type StorageService interface {
	UploadFile(file *multipart.FileHeader, objectKey string) error
	DeleteFile(objectKey string) error
	GetPublicURL(objectKey string) string
	BuildObjectKey(appID uint, originalName string) string
}

type s3StorageService struct {
	client     *s3.Client
	bucketName string
	publicURL  string
}

func NewS3StorageService(cfg *config.Config) (StorageService, error) {
	if cfg.R2AccountID == "" || cfg.R2AccessKeyID == "" || cfg.R2SecretAccessKey == "" {
		logger.Log.Warn("Storage R2 credentials are not fully configured")
		return nil, fmt.Errorf("missing R2 credentials")
	}

	endpoint := cfg.R2Endpoint
	if endpoint == "" {
		endpoint = fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.R2AccountID)
	}

	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:               endpoint,
			HostnameImmutable: true,
			SigningRegion:     "auto",
		}, nil
	})

	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithEndpointResolverWithOptions(r2Resolver),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.R2AccessKeyID, cfg.R2SecretAccessKey, "")),
		awsconfig.WithRegion("auto"),
	)

	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg)

	return &s3StorageService{
		client:     client,
		bucketName: cfg.R2BucketName,
		publicURL:  cfg.R2PublicBaseURL,
	}, nil
}

func (s *s3StorageService) UploadFile(fileHeader *multipart.FileHeader, objectKey string) error {
	file, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer file.Close()

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err = s.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.bucketName),
		Key:         aws.String(objectKey),
		Body:        file,
		ContentType: aws.String(contentType),
	})

	if err != nil {
		logger.Log.Error("Failed to upload file to S3", zap.Error(err))
		return err
	}

	return nil
}

func (s *s3StorageService) DeleteFile(objectKey string) error {
	_, err := s.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		logger.Log.Error("Failed to delete file from S3", zap.Error(err))
		return err
	}
	return nil
}

func (s *s3StorageService) GetPublicURL(objectKey string) string {
	if s.publicURL == "" {
		return ""
	}
	// Join publicURL with objectKey correctly
	baseURL := strings.TrimRight(s.publicURL, "/")
	// Important to escape objectKey
	// Some r2 edge cases need raw path but url.PathEscape helps with spaces etc. 
	// However, path separator '/' should not be escaped. 
	// The standard way:
	escapedPath := (&url.URL{Path: objectKey}).String()
	
	return fmt.Sprintf("%s%s", baseURL, escapedPath)
}

func (s *s3StorageService) BuildObjectKey(appID uint, originalName string) string {
	ext := path.Ext(originalName)
	newID := uuid.New().String()
	
	now := time.Now().Format("2006/01")
	return fmt.Sprintf("applications/%d/%s/%s%s", appID, now, newID, ext)
}

package s3client

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type BucketClient interface {
	UploadObject(ctx context.Context, bucket, fileName string, body io.Reader) (string, error)
}

type S3 struct {
	timeout    time.Duration
	client     *s3.S3
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
}

// NewS3 creates a new session with S3 AWS
func NewS3(session *session.Session, timeout time.Duration) BucketClient {
	return &S3{
		timeout:    timeout,
		client:     s3.New(session),
		uploader:   s3manager.NewUploader(session),
		downloader: s3manager.NewDownloader(session),
	}
}

// UploadObject uploads an object to the specified S3 bucket
func (s S3) UploadObject(ctx context.Context, bucket, fileName string, body io.Reader) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	res, err := s.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Body:   body,
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	})
	if err != nil {
		return "", err
	}

	return res.Location, nil
}

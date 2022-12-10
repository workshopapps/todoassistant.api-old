package awss3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// type AWSSession struct {
// 	s3session *s3.S3
// }

func NewAWSSession(id, secret, token string) (*s3.S3, error) {
	// var s3session *s3.S3
	s3session := s3.New(session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(id, secret, token),
	})))

	return s3session, nil
	// return &AWSSession{s3session: s3session}, nil
}

// func (a *AWSSession) UploadImage(file *multipart.FileHeader, filename string) error {
// 	image, err := file.Open()
// 	if err != nil {
// 		return err
// 	}
// 	defer image.Close()

// 	_, err = a.s3session.PutObject(&s3.PutObjectInput{
// 		Body:        image,
// 		Bucket:      aws.String("ticked-v1-backend-bucket"),
// 		Key:         aws.String(filename),
// 		ContentType: aws.String("image/jpeg"),
// 		// ACL:    aws.String(s3.BucketCannedACLPublicRead),
// 	})

// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

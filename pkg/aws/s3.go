package aws

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3 struct {
	Region string
	Bucket string
}

func NewS3(regin string, bucket string) S3 {
	return S3{
		Region: regin,
		Bucket: bucket,
	}
}

func (s *S3) Upload(imagePath string, objectKey string) error {
	// sess := session.Must(session.NewSessionWithOptions(session.Options{
	// 	SharedConfigState: session.SharedConfigEnable,
	// }))
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(s.Region),
	}))

	svc := s3.New(sess)

	file, err := os.Open(imagePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(objectKey),
		Body:   file,
	})
	if err != nil {
		return err
	}

	log.Printf("Successfully uploaded image to S3: %s/%s\n", imagePath, objectKey)
	return nil
}

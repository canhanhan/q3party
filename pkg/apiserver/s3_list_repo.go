package apiserver

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3ListRepository struct {
	downloader *s3manager.Downloader
	uploader   *s3manager.Uploader
	bucket     string
}

func (r *S3ListRepository) Get(id string) (*List, error) {
	buf := aws.NewWriteAtBuffer([]byte{})
	_, err := r.downloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(fmt.Sprintf("%s.json", id)),
	})
	if err != nil {
		e, ok := err.(awserr.RequestFailure)
		if ok {
			if e.StatusCode() == 404 {
				return r.Create(id)
			}
		}

		return nil, err
	}

	var list List
	err = json.Unmarshal(buf.Bytes(), &list)
	if err != nil {
		return nil, err
	}

	return &list, nil
}

func (r *S3ListRepository) Create(id string) (*List, error) {
	list := &List{
		ID:      id,
		Servers: []string{},
	}

	err := r.Save(list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (r *S3ListRepository) Save(list *List) error {
	data, err := json.Marshal(list)
	if err != nil {
		return err
	}

	_, err = r.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(r.bucket),
		Key:    aws.String(fmt.Sprintf("%s.json", list.ID)),
		Body:   bytes.NewReader(data),
	})
	return err
}

func NewS3ListRepository(region string, id string, secret string, bucket string) (*S3ListRepository, error) {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(id, secret, ""),
	}))
	downloader := s3manager.NewDownloader(sess)
	uploader := s3manager.NewUploader(sess)

	return &S3ListRepository{
		downloader: downloader,
		uploader:   uploader,
		bucket:     bucket,
	}, nil
}

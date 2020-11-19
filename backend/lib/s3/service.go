package s3

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

type s3Service struct {
	AWSs3 s3iface.S3API
}

type AWSServiceS3Params struct {
	Region    string
	AccessKey string
	SecretKey string
	PublicUrl string
	PathStyle bool
}

func NewService(params AWSServiceS3Params) (s3Service, error) {
	if params.Region == "" {
		return s3Service{}, fmt.Errorf("missing region")
	}

	return s3Service{
		AWSs3: s3.New(session.Must(session.NewSession(&aws.Config{
			Region:           aws.String(params.Region),
			Credentials:      credentials.NewStaticCredentials(params.AccessKey, params.SecretKey, ""),
			Endpoint:         &params.PublicUrl,
			S3ForcePathStyle: &params.PathStyle,
		}))),
	}, nil
}

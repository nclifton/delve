package rpc

import (
	"encoding/gob"

	"github.com/burstsms/mtmo-tp/backend/lib/rabbit"
	"github.com/burstsms/mtmo-tp/backend/lib/redis"
	"github.com/burstsms/mtmo-tp/backend/lib/rpc"
	mms "github.com/burstsms/mtmo-tp/backend/mms/rpc/client"
)

const Name = "MM7"

type Env struct {
	ContainerName      string `envconfig:"CONTAINER_NAME"`
	ContainerPort      int    `envconfig:"CONTAINER_PORT"`
	AWSRegion          string `envconfig:"AWS_REGION"`
	AWSAccessKey       string `envconfig:"AWS_ACCESS_KEY"`
	AWSSecretKey       string `envconfig:"AWS_SECRET_KEY"`
	AWSS3PublicUrl     string `envconfig:"AWS_S3_PUBLIC_URL"`
	AWSS3PathStyle     bool   `envconfig:"AWS_S3_PATH_STYLE"`
	MMSMediaBucket     string `envconfig:"MMS_MEDIA_BUCKET"`
	RabbitURL          string `envconfig:"RABBIT_URL"`
	RabbitExchange     string `envconfig:"RABBIT_EXCHANGE"`
	RabbitExchangeType string `envconfig:"RABBIT_EXCHANGE_TYPE"`
	RedisURL           string `envconfig:"REDIS_URL"`

	MMSRPCAddress string `envconfig:"MMS_RPC_ADDRESS"`

	NRName    string `envconfig:"NR_NAME"`
	NRLicense string `envconfig:"NR_LICENSE"`
	NRTracing bool   `envconfig:"NR_TRACING"`
}

type s3Svc interface {
	PutS3Content(content []byte, bucket, key string) error
}

type mmsSvc interface {
	UpdateStatus(p mms.UpdateStatusParams) (err error)
}

type ConfigSvc struct {
	S3  s3Svc
	MMS mmsSvc
}

type ConfigVar struct {
	AWSRegion      string
	MMSMediaBucket string
}

type MM7 struct {
	db        *db
	name      string
	svc       ConfigSvc
	configVar ConfigVar
}

type Service struct {
	receiver *MM7
}

func (s *Service) Name() string {
	return s.receiver.name
}

func (s *Service) Receiver() interface{} {
	return s.receiver
}

func NewService(r rabbit.Conn, opts RabbitPublishOptions, redis *redis.Connection, limiter *redis.Limiter, svc ConfigSvc, configVar ConfigVar) rpc.Service {
	gob.Register(map[string]interface{}{})
	return &Service{
		receiver: &MM7{db: &db{rabbit: r, opts: opts, redis: redis, limiter: limiter}, svc: svc, name: Name, configVar: configVar},
	}
}

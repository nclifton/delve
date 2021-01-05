package client

import (
	"fmt"
	"log"

	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"

	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
)

func NewClient(host string, port int, tracer opentracing.Tracer) webhookpb.ServiceClient {
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", host, port),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(tracer)))
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return webhookpb.NewServiceClient(conn)
}

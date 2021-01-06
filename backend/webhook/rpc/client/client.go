package client

import (
	"fmt"
	"log"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"

	"github.com/burstsms/mtmo-tp/backend/webhook/rpc/webhookpb"
)

func NewClient(host string, port int, tracer opentracing.Tracer) webhookpb.ServiceClient {
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%d", host, port),
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(
			otgrpc.OpenTracingClientInterceptor(tracer, otgrpc.LogPayloads()),
		),
	)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	return webhookpb.NewServiceClient(conn)
}

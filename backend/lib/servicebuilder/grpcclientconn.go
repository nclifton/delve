package servicebuilder

import (
	"fmt"
	"log"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

func NewClientConn(host string, port int, tracer opentracing.Tracer) *grpc.ClientConn {
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

	return conn
}

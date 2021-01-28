package rpcbuilder

import (
	"log"

	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
)

func NewClientConn(addr string, tracer opentracing.Tracer) *grpc.ClientConn {
	conn, err := grpc.Dial(
		addr,
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

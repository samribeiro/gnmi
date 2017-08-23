package main

import (
	"context"
	"flag"
	"net"

	"ribeiro/gnmi/helper"

	log "github.com/golang/glog"
	"github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
)

var (
	bind = flag.String("port", ":32123", "Bind to address:port or just :port.")
)

type subscriber struct {
	grpc.ServerStream
}

func (subs *subscriber) Send(m *gnmi.SubscribeResponse) error {
	return nil
}

func (subs *subscriber) Recv() (*gnmi.SubscribeRequest, error) {
	return nil, nil
}

type server struct{}

func (s *server) Capabilities(ctx context.Context, in *gnmi.CapabilityRequest) (*gnmi.CapabilityResponse, error) {
	log.Infoln("served Capabilities request")
	return nil, grpc.Errorf(codes.Unimplemented, "Capabilities() is not implemented.")
}

func (s *server) Get(ctx context.Context, in *gnmi.GetRequest) (*gnmi.GetResponse, error) {
	log.Infoln("served a Get request")
	return helper.ReflectGetRequest(in), nil
}

func (s *server) Set(ctx context.Context, in *gnmi.SetRequest) (*gnmi.SetResponse, error) {
	log.Infoln("served a Set request")
	return nil, grpc.Errorf(codes.Unimplemented, "Set() is not implemented.")
}

func (s *server) Subscribe(subs gnmi.GNMI_SubscribeServer) error {
	log.Infoln("served a Subscribe request")
	return grpc.Errorf(codes.Unimplemented, "Subscribe() is not implemented.")
}

func main() {
	flag.Parse()

	creds := helper.ServerCertificates()

	s := grpc.NewServer(grpc.Creds(creds))

	gnmi.RegisterGNMIServer(s, &server{})
	reflection.Register(s)

	log.Infoln("starting to listen on", *bind)
	listen, err := net.Listen("tcp", *bind)
	if err != nil {
		log.Exitf("failed to listen: %v", err)
	}

	log.Infoln("starting to serve")
	if err := s.Serve(listen); err != nil {
		log.Exitf("failed to serve: %v", err)
	}
}

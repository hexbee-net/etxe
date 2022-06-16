package plugin

import (
	"crypto/tls"
	"fmt"
	"github.com/hexbee-net/etxe/internal/plugin"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

const GRPCServiceName = "plugin"

type GRPCServer struct {
	// Server is the actual server that will accept connections.
	// This will be used for plugin registration as well.
	Server func([]grpc.ServerOption) *grpc.Server

	// TLS should be the TLS configuration if available.
	// If this is nil, the connection will not have transport security.
	TLS *tls.Config

	server *grpc.Server
	broker *GRPCBroker
}

func (s *GRPCServer) Init() error {
	var opts []grpc.ServerOption
	if s.TLS != nil {
		opts = append(opts, grpc.Creds(credentials.NewTLS(s.TLS)))
	}

	s.server = s.Server(opts)

	// Register the health service
	healthCheck := health.NewServer()
	healthCheck.SetServingStatus(GRPCServiceName, grpc_health_v1.HealthCheckResponse_SERVING)
	grpc_health_v1.RegisterHealthServer(s.server, healthCheck)

	// Register the reflection service
	reflection.Register(s.server)

	// Register the broker service
	brokerServer := newGRPCBrokerServer()
	plugin.RegisterGRPCBrokerServer(s.server, brokerServer)
	s.broker = newGRPCBroker(brokerServer, s.TLS)
	go s.broker.Run()

	// Register the controller
	controllerServer := &grpcControllerServer{server: s}
	plugin.RegisterGRPCControllerServer(s.server, controllerServer)

	// Register the stdio service
	s.stdioServer = newGRPCStdioServer(s.logger, s.Stdout, s.Stderr)
	plugin.RegisterGRPCStdioServer(s.server, s.stdioServer)

	// Register all our plugins onto the gRPC server.
	for k, raw := range s.Plugins {
		p, ok := raw.(GRPCPlugin)
		if !ok {
			return fmt.Errorf("%q is not a GRPC-compatible plugin", k)
		}

		if err := p.GRPCServer(s.broker, s.server); err != nil {
			return fmt.Errorf("error registering %q: %s", k, err)
		}
	}

	return nil
}

func (s *GRPCServer) Config() string {
	// TODO implement me
	panic("implement me")
}

func (s *GRPCServer) Serve(listener net.Listener) {
	// TODO implement me
	panic("implement me")
}

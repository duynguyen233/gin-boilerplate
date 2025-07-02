package grpcserver

import (
	"context"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
)

const (
	_defaultShutdownTimeout = 3 * time.Second
)

// Server -.
type Server struct {
	server          *grpc.Server
	notify          chan error
	shutdownTimeout time.Duration
	port            int
}

// New -.
func New(port int, opts ...grpc.ServerOption) *Server {
	grpcServer := grpc.NewServer(opts...)
	s := &Server{
		server:          grpcServer,
		notify:          make(chan error, 1),
		shutdownTimeout: _defaultShutdownTimeout,
		port:            port,
	}

	return s
}

// RegisterService registers a service with the gRPC server.
func (s *Server) RegisterService(registerFunc func(*grpc.Server)) {
	if registerFunc != nil {
		registerFunc(s.server)
	}
}

func (s *Server) Start() {
	go func() {
		lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
		if err != nil {
			s.notify <- fmt.Errorf("failed to listen on port %d: %w", s.port, err)
			close(s.notify)
			return
		}
		err = s.server.Serve(lis)
		s.notify <- err
		close(s.notify)
	}()
}

// Notify -.
func (s *Server) Notify() <-chan error {
	return s.notify
}

// Shutdown -.
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	done := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(done)
	}()

	select {
	case <-ctx.Done():
		s.server.Stop()
		return ctx.Err()
	case <-done:
		return nil
	}
}

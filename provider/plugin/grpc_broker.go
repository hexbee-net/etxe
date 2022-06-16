package plugin

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"math"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/oklog/run"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/hexbee-net/etxe/internal/plugin"
)

// streamer interface is used in the broker to send/receive connection information.
type streamer interface {
	Send(*plugin.ConnInfo) error
	Receive() (*plugin.ConnInfo, error)
	Close()
}

// sendErr is used to pass errors back during a send.
type sendErr struct {
	i  *plugin.ConnInfo
	ch chan error
}

// /////////////////////////////////////////////////////////////////////////////

// grpcBrokerServer is used by the plugin to start a stream and to send
// connection information to/from the plugin. Implements GRPCBrokerServer and
// streamer interfaces.
type grpcBrokerServer struct {
	plugin.UnimplementedGRPCBrokerServer

	// send is used to send connection info to the gRPC stream.
	send chan *sendErr

	// receive is used to receive connection info from the gRPC stream.
	receive chan *plugin.ConnInfo

	// quit closes down the stream.
	quit chan struct{}

	// o is used to ensure we close the quit channel only once.
	o sync.Once
}

func newGRPCBrokerServer() *grpcBrokerServer {
	return &grpcBrokerServer{
		send:    make(chan *sendErr),
		receive: make(chan *plugin.ConnInfo),
		quit:    make(chan struct{}),
	}
}

// StartStream implements the GRPCBrokerServer interface and will block until
// the quit channel is closed or the context reports Done.
// The stream will pass connection information to/from the client.
func (s *grpcBrokerServer) StartStream(stream plugin.GRPCBroker_StartStreamServer) error {
	doneCh := stream.Context().Done()
	defer s.Close()

	// Process send stream
	go func() {
		for {
			select {
			case <-doneCh:
				return
			case <-s.quit:
				return
			case se := <-s.send:
				err := stream.Send(se.i)
				se.ch <- err
			}
		}
	}()

	// Process receive stream
	for {
		i, err := stream.Recv()
		if err != nil {
			return err
		}
		select {
		case <-doneCh:
			return nil
		case <-s.quit:
			return nil
		case s.receive <- i:
		}
	}
}

// Send is used by the GRPCBroker to pass connection information into the stream
// to the client.
func (s *grpcBrokerServer) Send(i *plugin.ConnInfo) error {
	ch := make(chan error)
	defer close(ch)

	select {
	case <-s.quit:
		return errors.New("broker closed")
	case s.send <- &sendErr{
		i:  i,
		ch: ch,
	}:
	}

	return <-ch
}

// Receive is used by the GRPCBroker to pass connection information that has been
// sent from the client from the stream to the broker.
func (s *grpcBrokerServer) Receive() (*plugin.ConnInfo, error) {
	select {
	case <-s.quit:
		return nil, errors.New("broker closed")
	case i := <-s.receive:
		return i, nil
	}
}

// Close closes the quit channel, shutting down the stream.
func (s *grpcBrokerServer) Close() {
	s.o.Do(func() {
		close(s.quit)
	})
}

// /////////////////////////////////////////////////////////////////////////////

// grpcBrokerClient is used by the client to start a stream and to send
// connection information to/from the client. Implements GRPCBrokerClient and
// streamer interfaces.
type grpcBrokerClient struct {
	// client is the underlying GRPC client used to make calls to the server.
	client plugin.GRPCBrokerClient

	// send is used to send connection info to the gRPC stream.
	send chan *sendErr

	// receive is used to receive connection info from the gRPC stream.
	receive chan *plugin.ConnInfo

	// quit closes down the stream.
	quit chan struct{}

	// o is used to ensure we close the quit channel only once.
	o sync.Once
}

func newGRPCBrokerClient(conn *grpc.ClientConn) *grpcBrokerClient {
	return &grpcBrokerClient{
		client:  plugin.NewGRPCBrokerClient(conn),
		send:    make(chan *sendErr),
		receive: make(chan *plugin.ConnInfo),
		quit:    make(chan struct{}),
	}
}

// StartStream implements the GRPCBrokerClient interface and will block until
// the quit channel is closed or the context reports Done.
// The stream will pass connection information to/from the plugin.
func (s *grpcBrokerClient) StartStream() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	defer s.Close()

	stream, err := s.client.StartStream(ctx)
	if err != nil {
		return err
	}
	doneCh := stream.Context().Done()

	go func() {
		for {
			select {
			case <-doneCh:
				return
			case <-s.quit:
				return
			case se := <-s.send:
				err := stream.Send(se.i)
				se.ch <- err
			}
		}
	}()

	for {
		i, err := stream.Recv()
		if err != nil {
			return err
		}
		select {
		case <-doneCh:
			return nil
		case <-s.quit:
			return nil
		case s.receive <- i:
		}
	}
}

// Send is used by the GRPCBroker to pass connection information into the stream
// to the plugin.
func (s *grpcBrokerClient) Send(i *plugin.ConnInfo) error {
	ch := make(chan error)
	defer close(ch)

	select {
	case <-s.quit:
		return errors.New("broker closed")
	case s.send <- &sendErr{
		i:  i,
		ch: ch,
	}:
	}

	return <-ch
}

// Receive is used by the GRPCBroker to pass connection information that has been
// sent from the plugin to the broker.
func (s *grpcBrokerClient) Receive() (*plugin.ConnInfo, error) {
	select {
	case <-s.quit:
		return nil, errors.New("broker closed")
	case i := <-s.receive:
		return i, nil
	}
}

// Close closes the quit channel, shutting down the stream.
func (s *grpcBrokerClient) Close() {
	s.o.Do(func() {
		close(s.quit)
	})
}

// /////////////////////////////////////////////////////////////////////////////

// GRPCBroker is responsible for brokering connections by unique ID.
//
// It is used by plugins to create multiple gRPC connections and data
// streams between the plugin process and the host process.
//
// This allows a plugin to request a channel with a specific ID to connect to
// or accept a connection from, and the broker handles the details of
// holding these channels open while they're being negotiated.
//
// The Plugin interface has access to these for both Server and Client.
// The broker can be used by either (optionally) to reserve and connect to
// new streams. This is useful for complex args and return values,
// or anything else you might need a data stream for.
type GRPCBroker struct {
	nextId   uint32
	streamer streamer
	streams  map[uint32]*gRPCBrokerPending
	tls      *tls.Config
	doneCh   chan struct{}
	o        sync.Once

	sync.Mutex
}

type gRPCBrokerPending struct {
	ch     chan *plugin.ConnInfo
	doneCh chan struct{}
}

func newGRPCBroker(s streamer, tls *tls.Config) *GRPCBroker {
	return &GRPCBroker{
		streamer: s,
		streams:  make(map[uint32]*gRPCBrokerPending),
		tls:      tls,
		doneCh:   make(chan struct{}),
	}
}

// Accept accepts a connection by ID.
//
// This should not be called multiple times with the same ID.
func (b *GRPCBroker) Accept(id uint32) (net.Listener, error) {
	listener, err := serverListener()
	if err != nil {
		return nil, err
	}

	err = b.streamer.Send(&plugin.ConnInfo{
		ServiceId: id,
		Network:   listener.Addr().Network(),
		Address:   listener.Addr().String(),
	})
	if err != nil {
		return nil, err
	}

	return listener, nil
}

// AcceptAndServe is used to accept a specific stream ID and immediately
// serve a gRPC server on that stream ID.
//
// This is used to easily serve complex arguments.
// Each AcceptAndServe call opens a new listener socket and sends the connection
// info down the stream to the dialer. Since a new connection is opened every
// call, these calls should be used sparingly.
// Multiple gRPC server implementations can be registered to a single
// AcceptAndServe call.
func (b *GRPCBroker) AcceptAndServe(id uint32, s func([]grpc.ServerOption) *grpc.Server) error {
	listener, err := b.Accept(id)
	if err != nil {
		return err
	}

	defer func(listener net.Listener) {
		_ = listener.Close()
	}(listener)

	var opts []grpc.ServerOption
	if b.tls != nil {
		opts = []grpc.ServerOption{grpc.Creds(credentials.NewTLS(b.tls))}
	}

	server := s(opts)

	// Here we use a run group to close this goroutine if the server is shutdown
	// or the broker is shutdown.
	var g run.Group

	// Serve on the listener, if shutting down call GracefulStop.
	g.Add(func() error {
		return server.Serve(listener)
	}, func(err error) {
		server.GracefulStop()
	})

	// block on the closeCh or the doneCh.
	// If we are shutting down, close the closeCh.
	closeCh := make(chan struct{})
	g.Add(func() error {
		select {
		case <-b.doneCh:
		case <-closeCh:
		}
		return nil
	}, func(err error) {
		close(closeCh)
	})

	// Block until we are done
	return g.Run()
}

// Close closes the stream and all servers.
func (b *GRPCBroker) Close() error {
	b.streamer.Close()
	b.o.Do(func() {
		close(b.doneCh)
	})
	return nil
}

// Dial opens a connection by ID.
func (b *GRPCBroker) Dial(id uint32) (conn *grpc.ClientConn, err error) {
	var c *plugin.ConnInfo

	// Open the stream
	p := b.getStream(id)
	select {
	case c = <-p.ch:
		close(p.doneCh)
	case <-time.After(5 * time.Second):
		return nil, ErrConnectionInfoTimeout
	}

	var addr net.Addr
	switch c.Network {
	case "tcp":
		addr, err = net.ResolveTCPAddr("tcp", c.Address)
	case "unix":
		addr, err = net.ResolveUnixAddr("unix", c.Address)
	default:
		err = fmt.Errorf("%s: %w", c.Address, ErrUnknownAddressType)
	}

	if err != nil {
		return nil, err
	}

	return dialGRPCConn(b.tls, netAddrDialer(addr))
}

// NextId returns a unique ID to use next.
func (b *GRPCBroker) NextId() uint32 {
	return atomic.AddUint32(&b.nextId, 1)
}

// Run starts the brokering and should be executed in a goroutine, since it
// blocks forever, or until the session closes.
//
// Uses of GRPCBroker never need to call this. It is called internally by
// the plugin host/client.
func (b *GRPCBroker) Run() {
	for {
		stream, err := b.streamer.Receive()
		if err != nil {
			// Once we receive an error, just exit
			break
		}

		// Initialize the waiter
		p := b.getStream(stream.ServiceId)
		select {
		case p.ch <- stream:
		default:
		}

		go b.timeoutWait(stream.ServiceId, p)
	}
}

func (b *GRPCBroker) getStream(id uint32) *gRPCBrokerPending {
	b.Lock()
	defer b.Unlock()

	p, ok := b.streams[id]
	if ok {
		return p
	}

	b.streams[id] = &gRPCBrokerPending{
		ch:     make(chan *plugin.ConnInfo, 1),
		doneCh: make(chan struct{}),
	}
	return b.streams[id]
}

func (b *GRPCBroker) timeoutWait(id uint32, p *gRPCBrokerPending) {
	// Wait for the stream to either be picked up and connected,
	// or for a timeout.
	select {
	case <-p.doneCh:
	case <-time.After(5 * time.Second):
	}

	b.Lock()
	defer b.Unlock()

	// Delete the stream so no one else can grab it
	delete(b.streams, id)
}

func netAddrDialer(addr net.Addr) func(ctx context.Context, address string) (net.Conn, error) {
	return func(ctx context.Context, address string) (net.Conn, error) {
		// Connect to the client
		conn, err := net.Dial(addr.Network(), addr.String())
		if err != nil {
			return nil, fmt.Errorf("failed to dial (%s, %s): %w", addr.Network(), addr.String(), err)
		}

		if tcpConn, ok := conn.(*net.TCPConn); ok {
			// Make sure to set keep alive so that the connection doesn't die
			_ = tcpConn.SetKeepAlive(true)
		}

		return conn, nil
	}
}

func dialGRPCConn(tls *tls.Config, dialer func(context.Context, string) (net.Conn, error), dialOpts ...grpc.DialOption) (*grpc.ClientConn, error) {
	var transportCredentials credentials.TransportCredentials
	if tls == nil {
		transportCredentials = insecure.NewCredentials()
	} else {
		transportCredentials = credentials.NewTLS(tls)
	}

	opts := []grpc.DialOption{
		grpc.WithContextDialer(dialer),
		grpc.FailOnNonTempDialError(true),
		grpc.WithTransportCredentials(transportCredentials),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32)),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(math.MaxInt32)),
	}

	// Add custom options if we have any
	opts = append(opts, dialOpts...)

	// Connect.
	// The first parameter is unused because we use a custom dialer that has
	// the state to see the address.
	conn, err := grpc.Dial("unused", opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to dial connection: %w", err)
	}

	return conn, nil
}

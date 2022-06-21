package plugin

import (
	"bufio"
	"bytes"
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/hexbee-net/etxe/internal/plugin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
)

// grpcStdioBuffer is the buffer size we try to fill when sending a chunk of
// stdio data.
// This is currently 1 KB for no reason other than that seems like enough (stdio
// data isn't that common) and is fairly low.
const grpcStdioBuffer = 1 * 1024

// grpcStdioServer implements the Stdio service and streams stdout/stderr.
type grpcStdioServer struct {
	plugin.UnimplementedGRPCStdioServer

	stdoutCh <-chan []byte
	stderrCh <-chan []byte
}

// newGRPCStdioServer creates a new grpcStdioServer and starts the stream
// copying for the given out and err readers.
//
// This must only be called once per srcOut, srcErr.
func newGRPCStdioServer(log *zap.Logger, srcOut, srcErr io.Reader) *grpcStdioServer {
	stdoutCh := make(chan []byte)
	stderrCh := make(chan []byte)

	// Begin copying the streams
	go copyChan(log, stdoutCh, srcOut)
	go copyChan(log, stderrCh, srcErr)

	// Construct our server
	return &grpcStdioServer{
		stdoutCh: stdoutCh,
		stderrCh: stderrCh,
	}
}

// StreamStdio streams our stdout/err as the response.
func (s *grpcStdioServer) StreamStdio(_ *empty.Empty, srv plugin.GRPCStdio_StreamStdioServer) error {
	// Share the same data value between runs.
	// Sending this over the wire marshals it, so we can reuse this.
	var data plugin.StdioData

	for {
		// Read our data
		select {
		case data.Data = <-s.stdoutCh:
			data.Channel = plugin.StdioData_STDOUT

		case data.Data = <-s.stderrCh:
			data.Channel = plugin.StdioData_STDERR

		case <-srv.Context().Done():
			return nil
		}

		// Not sure if this is possible, but if we somehow got here,
		// and we didn't populate any data at all, then just continue.
		if len(data.Data) == 0 {
			continue
		}

		// Send our data to the client.
		if err := srv.Send(&data); err != nil {
			return err
		}
	}
}

// grpcStdioClient wraps the stdio service as a client to copy
// the stdio data to output writers.
type grpcStdioClient struct {
	log         *zap.Logger
	stdioClient plugin.GRPCStdio_StreamStdioClient
}

// newGRPCStdioClient creates a grpcStdioClient. This will perform the
// initial connection to the stdio service. If the stdio service is unavailable
// then this will be a no-op. This allows this to work without error for
// plugins that don't support this.
func newGRPCStdioClient(ctx context.Context, log *zap.Logger, conn *grpc.ClientConn) (*grpcStdioClient, error) {
	client := plugin.NewGRPCStdioClient(conn)

	// Connect immediately to the endpoint
	stdioClient, err := client.StreamStdio(ctx, &empty.Empty{})

	// If we get an Unavailable or Unimplemented error, this means that the plugin isn't
	// updated and linking to the latest version of go-plugin that supports
	// this. We fall back to the previous behavior of just not syncing anything.
	if status.Code(err) == codes.Unavailable || status.Code(err) == codes.Unimplemented {
		log.Warn("stdio service not available, stdout/stderr syncing unavailable")
		stdioClient = nil
		err = nil
	}
	if err != nil {
		return nil, err
	}

	return &grpcStdioClient{
		log:         log,
		stdioClient: stdioClient,
	}, nil
}

// Run starts the loop that receives stdio data and writes it to the given
// writers. This blocks and should be run in a goroutine.
func (c *grpcStdioClient) Run(stdout, stderr io.Writer) {
	// This will be nil if stdio is not supported by the plugin
	if c.stdioClient == nil {
		c.log.Warn("stdio service unavailable, run will do nothing")
		return
	}

	for {
		data, err := c.stdioClient.Recv()
		if err != nil {
			if err == io.EOF ||
				status.Code(err) == codes.Unavailable ||
				status.Code(err) == codes.Canceled ||
				status.Code(err) == codes.Unimplemented ||
				err == context.Canceled {

				return
			}

			c.log.Error("error receiving data", zap.Error(err))
			return
		}

		// Determine our output writer based on channel
		var w io.Writer
		switch data.Channel {
		case plugin.StdioData_STDOUT:
			w = stdout

		case plugin.StdioData_STDERR:
			w = stderr

		default:
			c.log.Warn("unknown channel, dropping", zap.Any("channel", data.Channel))
			continue
		}

		if _, err := io.Copy(w, bytes.NewReader(data.Data)); err != nil {
			c.log.Error("failed to copy all bytes", zap.Error(err))
		}
	}
}

// copyChan copies an io.Reader into a channel.
func copyChan(log *zap.Logger, dst chan<- []byte, src io.Reader) {
	sourceBuf := bufio.NewReader(src)

	for {
		// Make our data buffer.
		// We allocate a new one per loop iteration so that we can send it over the channel.
		var data [grpcStdioBuffer]byte

		// Read the data, this will block until data is available
		n, err := sourceBuf.Read(data[:])

		// We have to check if we have data BEFORE err != nil.
		// The bufio docs guarantee n == 0 on EOF, but it's better to be safe.
		if n > 0 {
			// We have data.
			// Send it on the channel. This will block if there is no reader
			// on the other side.
			// We expect that go-plugin will connect immediately to the stdio
			// server to drain this, so we want this block to happen
			// for backpressure.
			dst <- data[:n]
		}

		// If we hit EOF we're done copying
		if err == io.EOF {
			log.Debug("stdio EOF, exiting copy loop")
			return
		}

		// Any other error we just exit the loop.
		// We don't expect there to be errors since our use case for this is
		// reading/writing from an in-process pipe (os.Pipe).
		if err != nil {
			log.Warn("error copying stdio data, stopping copy", zap.Error(err))
			return
		}
	}
}

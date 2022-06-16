package plugin

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"runtime"
	"strconv"

	"go.uber.org/zap"
)

const (
	envClientCert = "PLUGIN_CLIENT_CERT"
)

type ServerConfig struct {
	// TLSProvider is a function that returns a configured tls.Config.
	TLSProvider func() (*tls.Config, error)

	// Logger is used to pass a logger into the server.
	// If none is provided the server will create a default logger.
	Logger *zap.Logger
}

// Serve serves the plugins given by ServerConfig.
//
// This is the method that plugins should call in their main() functions.
func Serve(config *ServerConfig) (err error) {
	logger := config.Logger
	if logger == nil {
		logger, err = zap.NewDevelopment()
		if err != nil {
			return fmt.Errorf("failed to create default logger: %w", err)
		}
	}

	// Register a listener so we can accept a connection
	listener, err := serverListener()
	if err != nil {
		return fmt.Errorf("failed to open listener: %w", err)
	}

	// Close the listener on return.
	defer func() {
		_ = listener.Close()
	}()

	var tlsConfig *tls.Config
	if config.TLSProvider != nil {
		tlsConfig, err = config.TLSProvider()
		if err != nil {
			return fmt.Errorf("failed to process TLS configuration: %w", err)
		}
	}

	var serverCert string
	clientCert := os.Getenv(envClientCert)

	// If the client is configured using AutoMTLS, the certificate will be here,
	// and we need to generate our own in response.
	if tlsConfig == nil && clientCert != "" {
		tlsConfig, serverCert, err = configureServerMTLS(logger, clientCert)
		if err != nil {
			return err
		}
	}

	if tlsConfig != nil {
		listener = tls.NewListener(listener, tlsConfig)
	}

	server := &GRPCServer{}

	if err := server.Init(); err != nil {
		return fmt.Errorf("server init failed: %w", err)
	}

	logger.Debug("server initialized",
		zap.String("network", listener.Addr().Network()),
		zap.String("address", listener.Addr().String()),
	)

	// Accept connections and wait for completion
	go server.Serve(listener)

	ctx := context.Background()

	select {
	case <-ctx.Done():
		_ = listener.Close()
	}

	// TODO implement me
	panic("implement me")
}

func configureServerMTLS(logger *zap.Logger, clientCert string) (*tls.Config, string, error) {
	logger.Info("configuring server automatic mTLS")
	clientCertPool := x509.NewCertPool()
	if !clientCertPool.AppendCertsFromPEM([]byte(clientCert)) {
		logger.Error("client cert provided but failed to parse", zap.String("cert", clientCert))
	}

	certPEM, keyPEM, err := generateCert()
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate server certificate: %w", err)
	}

	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, "", fmt.Errorf("failed to parse server certificate: %w", err)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    clientCertPool,
		MinVersion:   tls.VersionTLS12,
		RootCAs:      clientCertPool,
		ServerName:   "localhost",
	}

	// We send back the raw leaf cert data for the client rather than the
	// PEM, since the protocol can't handle newlines.
	serverCert := base64.RawStdEncoding.EncodeToString(cert.Certificate[0])

	return tlsConfig, serverCert, nil
}

func serverListener() (net.Listener, error) {
	switch runtime.GOOS {
	case "windows":
		return serverListenerTCP()
	default:
		return serverListenerUnix()
	}
}

func serverListenerTCP() (net.Listener, error) {
	envMinPort := os.Getenv("PLUGIN_MIN_PORT")
	envMaxPort := os.Getenv("PLUGIN_MAX_PORT")

	var minPort, maxPort int64
	var err error

	switch {
	case len(envMinPort) == 0:
		minPort = 0
	default:
		minPort, err = strconv.ParseInt(envMinPort, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("failed to get value from PLUGIN_MIN_PORT: %w", err)
		}
	}

	switch {
	case len(envMaxPort) == 0:
		maxPort = 0
	default:
		maxPort, err = strconv.ParseInt(envMaxPort, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("failed to get value from PLUGIN_MAX_PORT: %w", err)
		}
	}

	if minPort > maxPort {
		return nil, ErrMinPortGreaterThanMaxPort
	}

	for port := minPort; port <= maxPort; port++ {
		address := fmt.Sprintf("127.0.0.1:%d", port)
		listener, err := net.Listen("tcp", address)
		if err == nil {
			return listener, nil
		}
	}

	return nil, ErrTCPListenerBindFailed
}

func serverListenerUnix() (net.Listener, error) {
	tf, err := ioutil.TempFile("", "plugin")
	if err != nil {
		return nil, err
	}
	path := tf.Name()

	// Close the file and remove it because it has to not exist for
	// the domain socket.
	if err := tf.Close(); err != nil {
		return nil, err
	}
	if err := os.Remove(path); err != nil {
		return nil, err
	}

	l, err := net.Listen("unix", path)
	if err != nil {
		return nil, err
	}

	// Wrap the listener in rmListener so that the Unix domain socket file
	// is removed on close.
	return &rmListener{
		Listener: l,
		Path:     path,
	}, nil
}

// rmListener is an implementation of net.Listener that forwards most
// calls to the listener but also removes a file as part of the close.
// We use this to clean up the unix domain socket on close.
type rmListener struct {
	net.Listener
	Path string
}

func (l *rmListener) Close() error {
	// Close the listener itself
	if err := l.Listener.Close(); err != nil {
		return err
	}

	// Remove the file
	return os.Remove(l.Path)
}

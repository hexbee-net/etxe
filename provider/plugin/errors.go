package plugin

import "errors"

var (
	ErrConnectionInfoTimeout = errors.New("timeout waiting for connection info")
	ErrUnknownAddressType    = errors.New("unknown address type")

	ErrTCPListenerBindFailed     = errors.New("couldn't bind plugin TCP listener")
	ErrMinPortGreaterThanMaxPort = errors.New("plugin min port is greater than max port")
)

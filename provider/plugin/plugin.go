package plugin

import "net/rpc"

// Plugin is the interface that is implemented to serve/connect to an
// inteface implementation.
type Plugin interface {
	// Server should return the RPC server compatible struct to serve
	// the methods that the Client calls over net/rpc.
	Server(*MuxBroker) (interface{}, error)

	// Client returns an interface implementation for the plugin you're
	// serving that communicates to the server end of the plugin.
	Client(*MuxBroker, *rpc.Client) (interface{}, error)
}

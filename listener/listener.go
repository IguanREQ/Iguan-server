package listener

import (
	"net/rpc"
)

var rpcServer *rpc.Server

func init() {
	rpcServer = rpc.NewServer()
}

// Entry point for other packages to register their methods
func RegisterMethods(rcvr interface{}) error {
	return rpcServer.Register(rcvr)
}

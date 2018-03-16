package listener

import (
	"iguan/logs"
	"net"
	"time"
	"iguan/auth"
	"github.com/amkulikov/extrpc"
	"log"
)

const (
	connKeepalivePeriod = 30 * time.Second
	connTtl             = 10 * time.Second
	connAcceptFailPause = 1 * time.Second
)

func handleTCPConnection(conn net.Conn) {
	defer conn.Close()

	conn.(*net.TCPConn).SetKeepAlive(true)
	conn.(*net.TCPConn).SetKeepAlivePeriod(connKeepalivePeriod)
	conn.SetReadDeadline(time.Now().Add(connTtl))

	caller := auth.NewCaller()
	if err := caller.ParseConn(conn); err != nil {
		log.Println(err)
		return
	}
	rpcServer.ServeCodec(extrpc.ExtendCodec(extrpc.NewGobServerCodec(conn), caller))
}

// Run json-rpc over tcp
func RunTCP(addr string) {
	logs.Info("LI :: Start RPC TCP Listener...")

	serv, err := net.Listen("tcp", addr)
	if err != nil {
		logs.Error("LI :: Can't start RPC TCP Listener: %v", err)
		return
	}

	logs.Info("LI :: RPC TCP Listener started!")
	for {
		conn, err := serv.Accept()
		if err != nil {
			logs.Error("SP :: Can't accept new connection: %v", err)
			time.Sleep(connAcceptFailPause)
			continue
		}

		go handleTCPConnection(conn)
	}
}

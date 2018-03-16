package listener

import (
	"io"
	"net/http"
	"net/rpc/jsonrpc"

	"github.com/amkulikov/extrpc"
	"iguan/auth"
	"iguan/logs"
)

// http://www.jsonrpc.org/historical/json-rpc-over-http.html
func httpHandler(rw http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		rw.Header().Set("Content-Type", "text/plain; charset=utf-8")
		rw.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(rw, "405 must POST\n")
		logs.Error("Disallowed http method: %v", req.Method)
		return
	}

	wr := extrpc.NewServerWrapper(rw, req)
	caller := auth.NewCaller()
	if err := caller.ParseHTTP(req); err != nil {
		logs.Error("ParseHTTP: %v", err)
		return
	}
	rpcServer.ServeCodec(extrpc.ExtendCodec(jsonrpc.NewServerCodec(wr), caller))
	<-wr.Wait()
	return
}

// Run json-rpc over tcp
func RunHTTP(server http.Server) {
	logs.Info("LI :: Start RPC HTTP Listener...")

	http.HandleFunc("/iguan", httpHandler)

	logs.Fatal("HTTP server fails: %v", server.ListenAndServe())
}

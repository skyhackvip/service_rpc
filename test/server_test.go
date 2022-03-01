package main

import (
	"github.com/skyhackvip/service_rpc/data"
	"github.com/skyhackvip/service_rpc/provider"
	"testing"
)

func TestServer(t *testing.T) {
	svr := provider.NewRPCServer("localhost:5558")
	go svr.Run()

	svr.Register("test", func() string {
		return "hello"
	})
	req := data.RPCData{Name: "test", Args: nil}
	res := svr.Call(req)
	t.Log("result:", res.Name, res.Args, res.Err)
}

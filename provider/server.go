package provider

import (
	"log"
	"reflect"
)

type RPCServer struct {
	listener Listener //*Listener is error
}

func NewRPCServer(ip string, port int) *RPCServer {
	return &RPCServer{
		listener: NewRPCListener(ip, port),
	}
}

//register service
func (svr *RPCServer) Register(name string, function interface{}) {
	handler := &RPCServerHandler{f: reflect.ValueOf(function)}
	svr.listener.SetHandler(name, handler) //check exitsted first
	log.Printf("%s registered success!\n", name)
}

//service start
func (svr *RPCServer) Run() {
	go svr.listener.Run()
}

//service stop
func (svr *RPCServer) Close() {
	if svr.listener != nil {
		log.Println("server close")
		svr.listener.Close()
	}
}

package provider

import (
	"log"
	"reflect"
)

type Server interface {
	Register(string, interface{}) //error
	Run()
	Close()
}

type RPCServer struct {
	listener Listener //*Listener is error
}

func NewRPCServer(ip string, port int) *RPCServer {
	return &RPCServer{
		listener: NewRPCListener(ip, port),
	}
}

//register service
func (svr *RPCServer) Register(class interface{}) {
	name := reflect.Indirect(reflect.ValueOf(class)).Type().Name()
	svr.RegisterName(name, class)
}

func (svr *RPCServer) RegisterName(name string, class interface{}) {
	handler := &RPCServerHandler{class: reflect.ValueOf(class)}
	svr.listener.SetHandler(name, handler)
	log.Printf("%s registered success!\n", name)
}

//service start
func (svr *RPCServer) Run() {
	go svr.listener.Run()
}

//service stop
func (svr *RPCServer) Close() {
	if svr.listener != nil {
		svr.listener.Close()
	}
}

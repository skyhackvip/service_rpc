package provider

import (
	"context"
	"github.com/skyhackvip/service_rpc/naming"
	"log"
	"reflect"
	"time"
)

type Server interface {
	Register(string, interface{}) //error
	Run()
	Close()
}

type Option struct {
	Ip           string
	Port         int
	Hostname     string
	AppId        string
	Env          string
	NetProtocol  string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var DefaultOption = Option{
	NetProtocol:  "tcp",
	ReadTimeout:  5 * time.Millisecond,
	WriteTimeout: 5 * time.Millisecond,
}

type RPCServer struct {
	listener Listener //*Listener is error
	registry naming.Registry
	option   Option
}

func NewRPCServer(option Option, registry naming.Registry) *RPCServer {
	return &RPCServer{
		listener: NewRPCListener(option),
		registry: registry,
		option:   option,
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
	//register in registry
	svr.registerToCenter()
}

//service stop
func (svr *RPCServer) Close() {
	if svr.listener != nil {
		svr.listener.Close()
	}
}

func (svr *RPCServer) registerToCenter() {
	instance := &naming.Instance{
		Env:      svr.option.Env,
		AppId:    svr.option.AppId,
		Hostname: svr.option.Hostname,
		Addrs:    svr.listener.GetAddrs(),
	}
	svr.registry.Register(context.Background(), instance)
}

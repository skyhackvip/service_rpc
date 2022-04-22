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
	Shutdown()
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
	ReadTimeout:  15 * time.Millisecond,
	WriteTimeout: 15 * time.Millisecond,
}

type RPCServer struct {
	listener   Listener //*Listener is error
	registry   naming.Registry
	cancelFunc context.CancelFunc
	option     Option
	Plugins    PluginContainer
}

func NewRPCServer(option Option, registry naming.Registry) *RPCServer {
	return &RPCServer{
		listener: NewRPCListener(option),
		registry: registry,
		option:   option,
		Plugins:  &pluginContainer{},
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
	svr.Plugins.RegisterHook(name, class)
	log.Printf("%s registered success!\n", name)
}

//service start
func (svr *RPCServer) Run() {
	svr.listener.SetPlugins(svr.Plugins)
	go svr.listener.Run()
	//register in discovery
	svr.registerToNaming()
}

//service close
func (svr *RPCServer) Close() {
	log.Println("close and cancel")
	if svr.listener != nil {
		svr.listener.Close()
	}
	svr.cancelFunc()
}

//service shutdown gracefully
func (svr *RPCServer) Shutdown() {
	log.Println("shutdown and cancel")
	if svr.listener != nil {
		svr.listener.Shutdown()
	}
	svr.cancelFunc()
}

func (svr *RPCServer) registerToNaming() error {
	instance := &naming.Instance{
		Env:      svr.option.Env,
		AppId:    svr.option.AppId,
		Hostname: svr.option.Hostname,
		Addrs:    svr.listener.GetAddrs(),
	}
	cancel, err := svr.registry.Register(context.Background(), instance)
	if err != nil {
		log.Println("register to naming error", err)
		return err
	}
	svr.cancelFunc = cancel
	return nil
}

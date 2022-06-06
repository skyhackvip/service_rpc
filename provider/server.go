package provider

import (
	"context"
	"errors"
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
	if option.NetProtocol == "" {
		option.NetProtocol = DefaultOption.NetProtocol
	}
	if option.ReadTimeout == 0 {
		option.ReadTimeout = DefaultOption.ReadTimeout
	}
	if option.WriteTimeout == 0 {
		option.WriteTimeout = DefaultOption.WriteTimeout
	}
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
	//register in discovery
	err := svr.registerToNaming()
	if err != nil {
		log.Fatal(err)
	}
	svr.listener.SetPlugins(svr.Plugins)
	go svr.listener.Run()
}

//service close
func (svr *RPCServer) Close() {
	log.Println("close and cancel: ", svr.option.AppId, svr.option.Hostname)
	if svr.listener != nil {
		svr.listener.Close()
	}
	svr.cancelFunc()
}

//service shutdown gracefully
func (svr *RPCServer) Shutdown() {
	log.Println("shutdown and cancel:", svr.option.AppId, svr.option.Hostname)
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
	retries := 2
	for retries > 0 {
		retries--
		cancel, err := svr.registry.Register(context.Background(), instance)
		if err == nil {
			log.Println("register to naming server success: ", svr.option.AppId, svr.option.Hostname)
			svr.cancelFunc = cancel
			return nil
		}
	}
	return errors.New("register to naming server fail")
}

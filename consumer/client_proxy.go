package consumer

import (
	"context"
	"errors"
	"github.com/skyhackvip/service_rpc/naming"
	"log"
	"strings"
)

type ClientProxy interface {
	//select
	//auth
	Call(context.Context, string, interface{}, ...interface{}) (interface{}, error)
}

type RPCClientProxy struct {
	option   Option
	registry naming.Registry
}

//, discovery ServiceDiscovery, lbMode LoadBalanceMode,failMode FailMode
func NewClientProxy(option Option, registry naming.Registry) ClientProxy {
	return &RPCClientProxy{
		//	failMode:   failMode,
		//	selectMode: selectMode,
		//	discovery:  discovery,
		option:   option,
		registry: registry,
	}
}

func (cp *RPCClientProxy) Call(ctx context.Context, servicePath string, stub interface{}, params ...interface{}) (interface{}, error) {
	service, err := NewService(servicePath)
	if err != nil {
		return nil, err
	}
	client := NewClient(cp.option)
	servers, err := cp.discoveryService(ctx, service.AppId)
	if err != nil {
		return nil, err
	}
	service.Addrs = servers
	addr := service.SelectAddr()
	addr = strings.Replace(addr, "tcp://", "", -1)
	log.Println("get server addr:" + addr)
	err = client.Connect(addr) //长连接管理
	if err != nil {
		return nil, err
	}
	retries := cp.option.Retries
	for retries > 0 {
		retries--
		return client.Invoke(ctx, service, stub, params...)
	}
	return nil, errors.New("error")
}

func (cp *RPCClientProxy) discoveryService(ctx context.Context, appId string) ([]string, error) {
	instances, ok := cp.registry.Fetch(ctx, appId)
	if !ok {
		return nil, errors.New("service not found")
	}
	var servers []string
	for _, instance := range instances {
		servers = append(servers, instance.Addrs...)
	}
	return servers, nil
}

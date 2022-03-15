package consumer

import (
	"context"
	"errors"
)

type ClientProxy interface {
	//select
	//auth
	Call(context.Context, string, interface{}, ...interface{}) (interface{}, error)
}

type RPCClientProxy struct {
	option Option
}

//, discovery ServiceDiscovery, lbMode LoadBalanceMode,failMode FailMode
func NewClientProxy(option Option) ClientProxy {
	return &RPCClientProxy{
		//	failMode:   failMode,
		//	selectMode: selectMode,
		//	discovery:  discovery,
		option: option,
	}
}

func (cp *RPCClientProxy) Call(ctx context.Context, servicePath string, stub interface{}, params ...interface{}) (interface{}, error) {
	service, err := NewService(servicePath)
	if err != nil {
		return nil, err
	}
	client := NewClient(cp.option)
	addr := service.SelectAddr()
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

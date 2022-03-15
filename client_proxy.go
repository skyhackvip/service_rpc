package consumer

import (
	"context"
	"errors"
)

type ClientProxy interface {
	//select
	//auth
	Call()
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

func (cp *RPCClientProxy) Call(ctx context.Context, methodName string, stub interface{}, params ...interface{}) (interface{}, error) {
	service, err := NewService(methodName)
	if err != nil {
		return nil, err
	}
	client := NewClient(option)
	addr := service.SelectAddr()
	err = client.Connect(addr) //长连接管理
	defer client.Close()
	if err != nil {
		return nil, err
	}
	retries := cp.option.Retries
	for retries > 0 {
		retries--
		return client.Invoke(ctx, service.Name, stub, params...)
	}
	return nil, errors.New("error")
}

package consumer

import (
	"context"
	"errors"
	"github.com/skyhackvip/service_rpc/global"
	"github.com/skyhackvip/service_rpc/naming"
	"log"
	"strings"
	"sync"
)

type ClientProxy interface {
	Call(context.Context, string, interface{}, ...interface{}) (interface{}, error)
}

type RPCClientProxy struct {
	failMode FailMode
	option   Option
	registry naming.Registry
	client   Client

	mutex       sync.RWMutex
	servers     []string
	loadBalance LoadBalance
}

func NewClientProxy(appId string, option Option, registry naming.Registry) ClientProxy {
	cp := &RPCClientProxy{
		option:   option,
		failMode: option.FailMode,
		registry: registry,
	}
	servers, err := cp.discoveryService(context.Background(), appId)
	if err != nil {
		log.Fatal(err)
	}
	cp.servers = servers
	cp.loadBalance = LoadBalanceFactory(option.LoadBalanceMode, cp.servers)
	cp.client = NewClient(cp.option)
	//watch server:if server addrs change, update loadBalance
	return cp
}

func (cp *RPCClientProxy) Call(ctx context.Context, servicePath string, stub interface{}, params ...interface{}) (interface{}, error) {
	service, err := NewService(servicePath)
	if err != nil {
		return nil, err
	}

	err = cp.getConn()
	if err != nil && cp.failMode == Failfast {
		log.Println("failfast:", err)
		return nil, err
	}

	//失败策略
	switch cp.failMode {
	case Failretry:
		retries := cp.option.Retries
		for retries > 0 {
			retries--
			if cp.client != nil {
				rs, err := cp.client.Invoke(ctx, service, stub, params...)
				if err == nil {
					return rs, nil
				}
			}
		}
	case Failover:
		retries := cp.option.Retries
		for retries > 0 {
			retries--
			if cp.client != nil {
				rs, err := cp.client.Invoke(ctx, service, stub, params...)
				//err == global.paramErr
				if err == nil || err == global.ParamErr {
					return rs, nil
				}
			}
			err = cp.getConn()
			log.Println("--failover new server--", cp.client.GetAddr())
		}
	case Failfast:
		if cp.client != nil {
			rs, err := cp.client.Invoke(ctx, service, stub, params...)
			if err == nil {
				return rs, nil
			}
			return nil, err
		}

	}
	return nil, errors.New("call error")
}

func (cp *RPCClientProxy) getConn() error {
	addr := strings.Replace(cp.loadBalance.Get(), cp.option.NetProtocol+"://", "", -1)
	err := cp.client.Connect(addr) //长连接管理
	if err != nil {
		log.Println("connect server fail:", err)
		return err
	}
	log.Println("connect server:" + addr)
	return nil
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
	log.Println(appId, " found service addrs: ", servers)
	return servers, nil
}

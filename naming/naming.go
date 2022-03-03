package naming

import "context"

type Instance struct {
	Env      string   `json:"env"`
	AppId    string   `json:"appid"`
	Hostname string   `json:"hostname"`
	Addrs    []string `json:"addrs"`
	Version  string   `json:"version"`
	Status   uint32   `json:"status"`
}

//注册
type Registry interface {
	Register(context.Context, *Instance) (context.CancelFunc, error)
	Close() error
}

//发现
type Resolver interface {
	Fetch(context.Context) (map[string]*Instance, bool)
	Close() error
}

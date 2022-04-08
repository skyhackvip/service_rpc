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

type Registry interface {
	Register(context.Context, *Instance) (context.CancelFunc, error)
	Fetch(context.Context, string) ([]*Instance, bool)
	Close() error
}

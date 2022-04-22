package provider

import (
	"errors"
	"fmt"
	"github.com/skyhackvip/service_rpc/protocol"
	"net"
)

type PluginContainer interface {
	Add(plugin Plugin)
	Remove(plugin Plugin)
	All() []Plugin

	RegisterHook(string, interface{}) error
	UnregisterHook(string) error
	ConnAcceptHook(net.Conn) (net.Conn, bool)
	//ConnCloseHook() error
	BeforeReadHook() error
	AfterReadHook(*protocol.RPCMsg, error) error
	BeforeCallHook(string, string, []interface{}) error
	AfterCallHook(string, string, []interface{}, []interface{}, error) error
	BeforeWriteHook([]byte) error
	AfterWriteHook([]byte, error) error
}

type Plugin interface{}

type RegisterPlugin interface {
	Register(name string, class interface{}) error
	Unregister(name string) error
}

type ConnAcceptPlugin interface {
	HandleConnAccept(net.Conn) (net.Conn, bool)
}

type BeforeReadPlugin interface {
	BeforeRead() error
}

type AfterReadPlugin interface {
	AfterRead(*protocol.RPCMsg, error) error
}

type BeforeCallPlugin interface {
	BeforeCall(string, string, []interface{}) error
}

type AfterCallPlugin interface {
	AfterCall(string, string, []interface{}, []interface{}, error) error
}

type BeforeWritePlugin interface {
	BeforeWrite([]byte) error
}

type AfterWritePlugin interface {
	AfterWrite([]byte, error) error
}

type pluginContainer struct {
	plugins []Plugin
}

func (p *pluginContainer) Add(plugin Plugin) {
	if plugin == nil {
		return
	}
	p.plugins = append(p.plugins, plugin)
}

func (p *pluginContainer) Remove(plugin Plugin) {
	if p.plugins == nil {
		return
	}
	res := make([]Plugin, 0, len(p.plugins))
	for _, v := range p.plugins {
		if v != plugin {
			res = append(res, v)
		}
	}
	p.plugins = res
}

func (p *pluginContainer) All() []Plugin {
	return p.plugins
}

func (p *pluginContainer) RegisterHook(name string, class interface{}) error {
	var errs string
	for _, v := range p.plugins {
		if registerPlugin, ok := v.(RegisterPlugin); ok { //*RegisterPlugin
			err := registerPlugin.Register(name, class)
			if err != nil {
				errs = fmt.Sprintf("%v\r%v", errs, err)
			}
		}
	}
	if len(errs) > 0 && errs != "" {
		return errors.New(errs)
	}
	return nil
}

func (p *pluginContainer) UnregisterHook(name string) error {
	var errs string
	for _, v := range p.plugins {
		if registerPlugin, ok := v.(RegisterPlugin); ok {
			err := registerPlugin.Unregister(name)
			if err != nil {
				errs = fmt.Sprintf("%v\r%v", errs, err)
			}
		}
	}
	if len(errs) > 0 && errs != "" {
		return errors.New(errs)
	}
	return nil
}

func (p *pluginContainer) ConnAcceptHook(conn net.Conn) (net.Conn, bool) {
	var flag bool
	for _, v := range p.plugins {
		if connAcceptPlugin, ok := v.(ConnAcceptPlugin); ok {
			conn, flag = connAcceptPlugin.HandleConnAccept(conn)
			if !flag {
				conn.Close()
				return conn, false
			}
		}
	}
	return conn, true
}

func (p *pluginContainer) BeforeReadHook() error {
	for _, v := range p.plugins {
		if plugin, ok := v.(BeforeReadPlugin); ok {
			err := plugin.BeforeRead()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *pluginContainer) AfterReadHook(msg *protocol.RPCMsg, err error) error {
	for _, v := range p.plugins {
		if plugin, ok := v.(AfterReadPlugin); ok {
			err := plugin.AfterRead(msg, err)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *pluginContainer) BeforeCallHook(class string, method string, args []interface{}) error {
	for _, v := range p.plugins {
		if plugin, ok := v.(BeforeCallPlugin); ok {
			err := plugin.BeforeCall(class, method, args)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *pluginContainer) AfterCallHook(class string, method string, args []interface{}, result []interface{}, err error) error {
	for _, v := range p.plugins {
		if plugin, ok := v.(AfterCallPlugin); ok {
			err := plugin.AfterCall(class, method, args, result, err)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *pluginContainer) BeforeWriteHook(res []byte) error {
	for _, v := range p.plugins {
		if plugin, ok := v.(BeforeWritePlugin); ok {
			err := plugin.BeforeWrite(res)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func (p *pluginContainer) AfterWriteHook(res []byte, err error) error {
	for _, v := range p.plugins {
		if plugin, ok := v.(AfterWritePlugin); ok {
			err := plugin.AfterWrite(res, err)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

package plugin

import (
	"fmt"
	"github.com/skyhackvip/service_rpc/protocol"
)

type BeforeReadPlugin struct{}

func (p BeforeReadPlugin) BeforeRead() error {
	fmt.Println("==== before read plugin ====")
	return nil
}

type AfterReadPlugin struct{}

func (p AfterReadPlugin) AfterRead(msg *protocol.RPCMsg, err error) error {
	fmt.Println("==== after read plugin ====", msg, err)
	return nil
}

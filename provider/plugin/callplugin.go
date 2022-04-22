package plugin

import (
	"fmt"
)

type BeforeCallPlugin struct {
}

func (p BeforeCallPlugin) BeforeCall(class string, method string, args []interface{}) error {
	fmt.Println("==== before call plugin ====", class, method, args)
	return nil
}

type MonitorPlugin struct {
}

func (p MonitorPlugin) AfterCall(class string, method string, args []interface{}, result []interface{}, err error) error {
	fmt.Println("==== after call plugin ====", class, method, args, result, err)
	return nil
}

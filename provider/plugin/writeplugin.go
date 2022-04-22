package plugin

import (
	"fmt"
)

type BeforeWritePlugin struct {
}

func (p BeforeWritePlugin) BeforeWrite(res []byte) error {
	fmt.Println("==== before write plugin ====", res)
	return nil
}

type AfterWritePlugin struct {
}

func (p AfterWritePlugin) AfterWrite(res []byte, err error) error {
	fmt.Println("==== after write plugin ====", res, err)
	return nil
}

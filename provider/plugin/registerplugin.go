package plugin

import "fmt"

type RegisterPlugin struct {
}

//*RegisterPlugin x
func (p RegisterPlugin) Register(name string, class interface{}) error {
	fmt.Println("==== register plugin ====", name, class)
	return nil
}

func (p RegisterPlugin) Unregister(name string) error {
	fmt.Println("==== unregister plugin ====", name)
	return nil
}

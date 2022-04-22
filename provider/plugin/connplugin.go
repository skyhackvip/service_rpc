package plugin

import (
	"fmt"
	"net"
)

type ConnPlugin struct{}

func (plugin ConnPlugin) HandleConnAccept(conn net.Conn) (net.Conn, bool) {
	fmt.Println("==== conn accept plugin ====", conn)
	return conn, true
}

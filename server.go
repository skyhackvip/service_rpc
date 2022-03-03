package main

import (
	"encoding/gob"
	"errors"
	"flag"
	"fmt"
	_ "github.com/skyhackvip/service_rpc/naming"
	"github.com/skyhackvip/service_rpc/provider"
	"github.com/skyhackvip/service_rpc/user"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	hostname string
	appid    string
	port     string
	ip       string
	env      string
)

func init() {
	fmt.Println(os.Hostname)
	if hostname, err := os.Hostname(); err != nil || hostname == "" {
		fmt.Println(err)
		hostname = os.Getenv("HOSTNAME") //system enviorment
	}
	if ip = InternalIP(); ip == "" {
		ip = os.Getenv("IP")
	}
	flag.StringVar(&appid, "appid", os.Getenv("APPID"), "appid required")
	flag.StringVar(&port, "port", os.Getenv("PORT"), "port required")
	flag.StringVar(&env, "env", os.Getenv("ENV"), "env required")
}

func main() {
	addr := ip + ":" + port
	fmt.Println(addr)

	//listen port
	srv := provider.NewRPCServer(addr)

	//register local service
	srv.Register("Test", Test)
	srv.Register("TestInt", TestInt)
	srv.Register("QueryUser", QueryUser)
	srv.Register("QueryUser1", QueryUser1)
	//gob
	gob.Register(user.User{})

	go srv.Run()

	//register to center
	/*nodes := []string{"localhost:8881"}
	conf := &naming.Config{Nodes: nodes, Env: env}
	dis := naming.New(conf)
	instance := &naming.Instance{
		AppId:    appid,
		Addrs:    []string{addr},
		Hostname: hostname,
	}
	dis.Register(instance)
	*/

	//graceful restart
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	<-quit
	srv.Close()
	//dis.Cancel()
}

func Test() string {
	return "hello world"
}

func TestInt(n string) (string, error) {
	if n == "1" {
		return "88888", nil
	} else {
		return "-1", errors.New("int err")
	}
}

var userList = map[int]user.User{
	1: user.User{"hero"},
	2: user.User{"kavin"},
}

func QueryUser(id int) (user.User, error) {
	fmt.Println("begin to query user", id)
	if u, ok := userList[id]; ok {
		return u, nil
	}
	return user.User{}, fmt.Errorf("id %d not found", id)
}

var userList1 = map[string]user.User{
	"a": user.User{"xxxxxxxwwwww"},
}

func QueryUser1(id string) (user.User, error) {
	fmt.Println("begin to query user1", id)
	if u, ok := userList1[id]; ok {
		return u, nil
	}
	return user.User{}, fmt.Errorf("id %s not found", id)
}

func InternalIP() string {
	inters, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, inter := range inters {
		if !strings.HasPrefix(inter.Name, "lo") {
			addrs, err := inter.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return ipnet.IP.String()
					}
				}
			}
		}
	}
	return ""
}

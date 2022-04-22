package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"github.com/skyhackvip/service_rpc/naming"
	"github.com/skyhackvip/service_rpc/provider"
	"github.com/skyhackvip/service_rpc/provider/plugin"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	//	"time"
)

var (
	hostname string
	appid    string
	port     int
	ip       string
	env      string
)

func init() {
	if ip = os.Getenv("IP"); ip == "" {
		ip = "localhost"
	}
	flag.StringVar(&appid, "appid", os.Getenv("APPID"), "appid required")
	flag.StringVar(&hostname, "hostname", os.Getenv("HOSTNAME"), "hostname required")
	flag.StringVar(&env, "env", os.Getenv("ENV"), "env required")
	port, _ = strconv.Atoi(os.Getenv("PORT"))
	flag.IntVar(&port, "port", port, "port required")
}

func main() {
	flag.Parse()
	if ip == "" || port == 0 || env == "" || appid == "" || hostname == "" {
		panic("init ip,port,env,appid,hostname error")
	}

	nodes := []string{"localhost:8881"}
	conf := &naming.Config{Nodes: nodes, Env: env}
	discovery := naming.New(conf)

	option := provider.Option{
		Ip:           ip,
		Port:         port,
		Hostname:     hostname,
		Env:          env,
		AppId:        appid,
		NetProtocol:  provider.DefaultOption.NetProtocol,
		ReadTimeout:  provider.DefaultOption.ReadTimeout,
		WriteTimeout: provider.DefaultOption.WriteTimeout,
	}

	srv := provider.NewRPCServer(option, discovery)
	srv.Plugins.Add(plugin.RegisterPlugin{})
	srv.Plugins.Add(plugin.ConnPlugin{})
	srv.Plugins.Add(plugin.BeforeReadPlugin{})
	srv.Plugins.Add(plugin.AfterReadPlugin{})
	srv.Plugins.Add(plugin.BeforeCallPlugin{})
	srv.Plugins.Add(plugin.MonitorPlugin{})
	srv.Plugins.Add(plugin.BeforeWritePlugin{})
	srv.Plugins.Add(plugin.AfterWritePlugin{})

	//register local service
	srv.RegisterName("User", &UserHandler{})
	srv.RegisterName("Test", &TestHandler{})

	//register gob
	gob.Register(User{})

	go srv.Run()

	//graceful restart
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	<-quit
	srv.Shutdown()
}

//test
type TestHandler struct{}

func (t *TestHandler) Hello() string {
	return "hello world"
}

func (t *TestHandler) Add(a, b int) int {
	return a + b
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var userList = map[int]User{
	1: User{1, "hero", 11},
	2: User{2, "kavin", 12},
}

type UserHandler struct{}

func (u *UserHandler) Login(name, pass string) bool {
	if name == "kavin" && pass == "123456" {
		return true
	}
	return false
}

func (u *UserHandler) GetUserById(id int) (User, error) {
	//time.Sleep(10 * time.Second)
	log.Println("start to query user", id)
	if u, ok := userList[id]; ok {
		return u, nil
	}
	return User{}, fmt.Errorf("id %d not found", id)
}

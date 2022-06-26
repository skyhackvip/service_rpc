package main

import (
	"encoding/gob"
	"errors"
	"flag"
	"github.com/skyhackvip/service_rpc/naming"
	"github.com/skyhackvip/service_rpc/provider"
	//"github.com/skyhackvip/service_rpc/provider/plugin"
	"context"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	c := flag.String("c", "", "config file path")
	flag.Parse()
	config, err := loadConfig(*c)
	if err != nil {
		log.Fatal("load config fail!", err)
	}

	conf := &naming.Config{Nodes: config.RegistryAddrs, Env: config.Env}
	discovery := naming.New(conf)

	option := provider.Option{
		Ip:       config.Ip,
		Port:     config.Port,
		Hostname: config.Hostname,
		Env:      config.Env,
		AppId:    config.Appid,
	}

	srv := provider.NewRPCServer(option, discovery)
	/*	srv.Plugins.Add(plugin.RegisterPlugin{})
		srv.Plugins.Add(plugin.ConnPlugin{})
		srv.Plugins.Add(plugin.BeforeReadPlugin{})
		srv.Plugins.Add(plugin.AfterReadPlugin{})
		srv.Plugins.Add(plugin.BeforeCallPlugin{})
		srv.Plugins.Add(plugin.MonitorPlugin{})
		srv.Plugins.Add(plugin.BeforeWritePlugin{})
		srv.Plugins.Add(plugin.AfterWritePlugin{})
	*/

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
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Println("start shutdown server")
	srv.Shutdown()
	select {
	case <-ctx.Done():
		log.Println("server close timeout")
	}
	log.Println("server exiting")
}

type Config struct {
	Hostname      string   `yaml:"hostname"`
	Appid         string   `yaml:"appid"`
	Port          int      `yaml:"port"`
	Ip            string   `yaml:"ip"`
	Env           string   `yaml:"env"`
	RegistryAddrs []string `yaml:"registry_addrs"`
}

func loadConfig(path string) (*Config, error) {
	configFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := new(Config)
	err = yaml.Unmarshal(configFile, config)
	if err != nil {
		return nil, err
	}
	return config, nil
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
	log.Println("start to query user", id)
	//time.Sleep(3 * time.Second)
	if u, ok := userList[id]; ok {
		return u, nil
	}
	return User{}, errors.New("id not found error! id=" + string(id))
}

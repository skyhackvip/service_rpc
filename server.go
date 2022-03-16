package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"github.com/skyhackvip/service_rpc/provider"
	"github.com/skyhackvip/service_rpc/util"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var (
	hostname string
	appid    string
	port     int
	ip       string
	env      string
)

func init() {
	if hostname, err := os.Hostname(); err != nil || hostname == "" {
		hostname = os.Getenv("HOSTNAME") //system enviorment
	}
	if ip = util.InternalIP(); ip == "" {
		ip = os.Getenv("IP")
	}
	flag.StringVar(&appid, "appid", os.Getenv("APPID"), "appid required")
	flag.StringVar(&env, "env", os.Getenv("ENV"), "env required")
	port, _ = strconv.Atoi(os.Getenv("PORT"))
	flag.IntVar(&port, "port", port, "port required")
}

func main() {
	flag.Parse()
	if ip == "" || port == 0 {
		panic("init ip and port error")
	}

	srv := provider.NewRPCServer(ip, port)

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
	srv.Close()
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
	if u, ok := userList[id]; ok {
		return u, nil
	}
	return User{}, fmt.Errorf("id %d not found", id)
}

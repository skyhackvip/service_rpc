package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/skyhackvip/service_rpc/provider"
	"github.com/skyhackvip/service_rpc/user"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	//fmt.Println(Test())
	//fmt.Println(QueryUser(1))
	//fmt.Println(QueryUser(3))

	//listen port
	srv := provider.NewRPCServer("localhost:3332")

	//register service
	srv.Register("Test", Test)
	srv.Register("TestInt", TestInt)
	srv.Register("QueryUser", QueryUser)
	srv.Register("QueryUser1", QueryUser1)
	//gob
	gob.Register(user.User{})

	go srv.Run()

	//graceful restart
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)
	<-quit
	srv.Close()
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

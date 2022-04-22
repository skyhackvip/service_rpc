package main

import (
	"context"
	"encoding/gob"
	"github.com/skyhackvip/service_rpc/consumer"
	"github.com/skyhackvip/service_rpc/naming"
	"log"
)

func main() {
	nodes := []string{"localhost:8881"}
	conf := &naming.Config{Nodes: nodes, Env: "dev"}
	discovery := naming.New(conf)

	gob.Register(User{})
	cli := consumer.NewClientProxy(consumer.DefaultOption, discovery)
	ctx := context.Background()

	var GetUserById func(id int) (User, error)
	cli.Call(ctx, "UserService.User.GetUserById", &GetUserById)
	u, err := GetUserById(2)
	log.Println("result:", u, err)
	/*
		var Hello func() string
		cli.Call(ctx, "UserService.Test.Hello", &Hello)
		r := Hello()
		log.Println("result:", r, err)

		var Add func(a, b int) int
		cli.Call(ctx, "UserService.Test.Add", &Add)
		w := Add(1, 2)
		log.Println("result:", w)

		var Login func(string, string) bool
		//for {
		cli.Call(ctx, "UserService.User.Login", &Login)
		v := Login("kavin", "123456")
		log.Println("result:", v)
	*/

	//	}
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

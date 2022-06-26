package main

import (
	"context"
	"encoding/gob"
	"github.com/skyhackvip/service_rpc/consumer"
	"github.com/skyhackvip/service_rpc/naming"
	"log"
	"reflect"
)

func main() {
	nodes := []string{"localhost:8881"}
	conf := &naming.Config{Nodes: nodes, Env: "dev"}
	discovery := naming.New(conf)

	gob.Register(User{})
	cli := consumer.NewClientProxy("UserService", consumer.DefaultOption, discovery)

	var GetUserById func(id int) (User, error)

	//wrap call
	ret, err := cli.Call(context.Background(), "User.GetUserById", &GetUserById, 1)
	if err != nil {
		log.Println("call error:", err)
	} else {
		val := ret.([]reflect.Value)
		user := val[0].Interface().(User)
		log.Println("rpc return result:", user)
	}

	//makefunc and call
	//u, err := GetUserById(2)
	for {

	}

	/*var Hello func() string
	cli.Call(ctx, "Test.Hello", &Hello)
	r := Hello()
	log.Println("result:", r, err)

	var Add func(a, b int) int
	cli.Call(ctx, "Test.Add", &Add)
	w := Add(1, 2)
	log.Println("result:", w)

	var Login func(string, string) bool
	cli.Call(ctx, "User.Login", &Login)
	v := Login("kavin", "123456")
	log.Println("result:", v)
	*/

}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

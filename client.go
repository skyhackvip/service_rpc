package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"github.com/skyhackvip/service_rpc/consumer"
)

func main() {
	gob.Register(User{})
	gob.Register(Order{})
	cli := consumer.NewClientProxy(consumer.DefaultOption)
	ctx := context.Background()

	var GetUserById func(id int) (User, error)
	cli.Call(ctx, "UserService.User.GetUserById", &GetUserById)
	u, err := GetUserById(2)
	fmt.Println("result:", u, err)

	var Hello func() string
	r, err := cli.Call(ctx, "UserService.Test.Hello", &Hello)
	fmt.Println("result:", r, err)

	var Add func(a, b int) int
	cli.Call(ctx, "UserService.Test.Add", &Add)
	r = Add(1, 2)
	fmt.Println("result:", r)

	var GetOrder func(int) Order
	cli.Call(ctx, "UserService.Order.GetOrder", &GetOrder)
	r = GetOrder(1)
	fmt.Println("result:", r)

	var Login func(string, string) bool
	cli.Call(ctx, "UserService.User.Login", &Login)
	r = Login("kavin", "123456")
	fmt.Println("result:", r)

}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type Order struct {
	OrderNo string
	Amount  float32
	Uid     int
}

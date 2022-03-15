package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"github.com/skyhackvip/service_rpc/consumer"
	"github.com/skyhackvip/service_rpc/user"
)

func main() {
	cli := consumer.NewClientProxy(consumer.DefaultOption)
	ctx := context.Background()
	var LocalTest func() string
	r, err := cli.Call(ctx, "user.Test", &LocalTest)
	fmt.Println(r, err)

	gob.Register(user.User{})
	var LocalQueryUser func(id int) (user.User, error)
	cli.Call(ctx, "user.QueryUser", &LocalQueryUser)
	u, err := LocalQueryUser(2)
	fmt.Println(u, err)
}

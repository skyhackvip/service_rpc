package main

import (
	"context"
	"encoding/gob"
	"github.com/skyhackvip/service_rpc/consumer"
	"log"
	"reflect"
)

func main() {
	gob.Register(User{})
	var GetUserById func(id int) (User, error)
	client := consumer.NewClient(consumer.DefaultOption)
	err := client.Connect("localhost:8898")
	if err != nil {
		panic(err)
	}

	service, err := consumer.NewService("User.GetUserById")
	if err != nil {
		panic(err)
	}

	//wrap call
	ret, err := client.Invoke(context.Background(), service, &GetUserById, 1)
	if err != nil {
		log.Println("call error:", err)
	} else {
		val := ret.([]reflect.Value)
		user := val[0].Interface().(User)
		log.Println("rpc return result:", user)
	}

}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

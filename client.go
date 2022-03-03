package main

import (
	//"encoding/gob"
	"fmt"
	"github.com/skyhackvip/service_rpc/consumer"
	//"github.com/skyhackvip/service_rpc/user"
	//"net"
)

func main() {

	var LocalTest func() string
	r, err := consumer.Call("user.Test", &LocalTest)
	fmt.Println(r)
	fmt.Println(err)
	consumer.ReflectCall(map[string]interface{}{"Test": LocalTest}, "Test")

	var LocalTest2 func(n string) (string, error)
	r, err = consumer.Call("user.TestInt", &LocalTest2, "1")
	fmt.Println(r)
	fmt.Println(err)
	consumer.ReflectCall(map[string]interface{}{"TestInt": LocalTest2}, "TestInt", "1")

	/*var LocalTest3 func(a, b string) (string, error)
	r, err = consumer.Call("user.Test3", &LocalTest3, "1", "2")
	fmt.Println(r)
	fmt.Println(err)
	*/

	/*conn, err := net.Dial("tcp", "localhost:3332")
	if err != nil {
		panic(err)
	}
	cli := consumer.NewClient(conn)

	var LocalTest func() string
	cli.Call("Test", &LocalTest)
	s := LocalTest()
	fmt.Println(s)

	fmt.Println("**************")

	var LocalTest2 func(n string) (string, error)
	cli.Call("TestInt", &LocalTest2)
	fmt.Println(LocalTest2("1"))

	fmt.Println("**************")

	gob.Register(user.User{})
	var LocalQueryUser func(id int) (user.User, error)
	cli.Call("QueryUser", &LocalQueryUser)
	u, err := LocalQueryUser(1)
	fmt.Println(u, err)

	fmt.Println("**************")

	var LocalQueryUser1 func(string) (user.User, error)
	cli.Call("QueryUser1", &LocalQueryUser1) //赋值到本地
	u, err = LocalQueryUser1("a")
	fmt.Println(u, err)
	*/

}

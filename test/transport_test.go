package main

import (
	"github.com/skyhackvip/service_rpc/transport"
	"net"
	"sync"
	"testing"
)

func TestTransport(t *testing.T) {
	addr := "localhost:2222"
	data := "sdlkjdsjklljkjdsf"
	var wg sync.WaitGroup
	wg.Add(2)
	//server
	go func() {
		defer wg.Done()
		l, err := net.Listen("tcp", addr)
		if err != nil {
			t.Fatal(err)
		}

		defer l.Close()
		conn, _ := l.Accept()
		trans := transport.NewTransport(conn)
		trans.Send([]byte(data))

	}()
	//client
	go func() {
		defer wg.Done()
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			t.Fatal(err)
		}
		trans := transport.NewTransport(conn)
		resData, err := trans.Read()
		if err != nil {
			t.Fatal(err)
		}
		t.Log("-------received:", string(data))
		if string(resData) != data {
			t.Fatal("not equal")
		}
	}()
	wg.Wait()
}

package main

import (
	"github.com/skyhackvip/service_rpc/data"
	"testing"
)

func TestData(t *testing.T) {
	d := data.RPCData{
		Name: "sdfjlksdfcesshi",
		Args: []interface{}{"aa", "bb"},
	}
	en, err := data.Encode(d)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("encode:", string(en))
	}
	de, err := data.Decode(en)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log("decode:", de.Name)
	}

}

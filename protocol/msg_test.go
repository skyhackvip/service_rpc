package protocol

import (
	"bytes"
	"github.com/skyhackvip/service_rpc/codec"
	"testing"
)

func TestMsg(t *testing.T) {
	req := NewRPCMsg()
	req.SetVersion(0) //0-255
	req.SetMsgType(Response)
	req.SetCompressType(Gzip)
	req.SetSerializeType(Gob)
	//	req.SetMsgId(888888888)

	/*
		meta := make(map[string]string)
		meta["color"] = "red"
		req.Metadata = meta
	*/

	req.ServiceClass = "User"
	req.ServiceMethod = "getUserById"
	data := map[string]int{"a": 1, "b": 2}
	coder := codec.GobCodec{}
	payload, _ := coder.Encode(data)
	req.Payload = payload

	var buf bytes.Buffer
	a, err := req.Send(&buf)
	t.Log(a, err)

	res, err := Read(&buf)
	if err != nil {
		t.Error(err)
	}
	t.Log(res.Header.Version())
	t.Log(res.Header.CompressType())
	t.Log(res.Header.MsgType())
	t.Log(res.Header.SerializeType())
	t.Log(res.ServiceClass)
	t.Log(res.ServiceMethod)
	rs := make(map[string]int, 0)
	err = coder.Decode(res.Payload, &rs)
	t.Log(rs, err)
}

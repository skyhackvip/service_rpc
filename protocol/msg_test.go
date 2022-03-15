package protocol

import (
	"bytes"
	"github.com/skyhackvip/service_rpc/codec"
	"github.com/skyhackvip/service_rpc/config"
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

	//req.Payload = []byte(`{"a":1, "b":2,}`)
	payload := map[string]int{"a": 1, "b": 2}
	coder := codec.New(config.CODEC_GOB)
	encodePayload, _ := coder.Encode(payload)
	req.Payload = encodePayload

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

/*
func TestConvert(t *testing.T) {
	s := "{aabbcdsfdsf}"
	t.Log(stringToByte(s))
	t.Log([]byte(s))
}*/

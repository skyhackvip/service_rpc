package global

import (
	"github.com/skyhackvip/service_rpc/codec"
	"github.com/skyhackvip/service_rpc/protocol"
)

var Codecs = map[protocol.SerializeType]codec.Codec{
	protocol.JSON: &codec.JSONCodec{},
	protocol.Gob:  &codec.GobCodec{},
}

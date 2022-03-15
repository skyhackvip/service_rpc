package protocol

const (
	HEADER_LEN = 5
)

const (
	magicNumber byte = 0x06
)

type MsgType byte

const (
	Request MsgType = iota
	Response
)

type CompressType byte

const (
	None CompressType = iota
	Gzip
)

type SerializeType byte

const (
	Gob SerializeType = iota
	JSON
)

type Header [HEADER_LEN]byte

func (h *Header) CheckMagicNumber() bool {
	return h[0] == magicNumber
}

func (h *Header) Version() byte {
	return h[1]
}

func (h *Header) SetVersion(version byte) {
	h[1] = version
}

func (h *Header) MsgType() MsgType {
	return MsgType(h[2])
}

func (h *Header) SetMsgType(msgType MsgType) {
	h[2] = byte(msgType)
}

func (h *Header) CompressType() CompressType {
	return CompressType(h[3])
}

func (h *Header) SetCompressType(compressType CompressType) {
	h[3] = byte(compressType)
}

func (h *Header) SerializeType() SerializeType {
	return SerializeType(h[4])
}

func (h *Header) SetSerializeType(serializeType SerializeType) {
	h[4] = byte(serializeType)
}

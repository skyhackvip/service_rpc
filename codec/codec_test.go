package codec

import (
	"github.com/skyhackvip/service_rpc/config"
	"testing"
)

type User struct {
	Id     int      `json:"id"`
	Name   string   `json:"name"`
	Habbit []string `json:"habbit"`
}

func TestJsonCodec(t *testing.T) {
	u := User{Id: 1,
		Name:   "Kavin",
		Habbit: []string{"coding", "touring"},
	}
	codec := JSONCodec{}
	b, err := codec.Encode(u)
	t.Log(b, err)

	ur := User{}
	err = codec.Decode(b, &ur)
	t.Log(ur, err)
}

func TestGobCodec(t *testing.T) {
	/*u := User{Id: 1,
		Name:   "Kavin",
		Habbit: []string{"coding", "touring"},
	}*/
	a := map[string]int{"a": 1, "b": 2}
	//codec := GobCodec{}
	coder := New(config.CODEC_GOB)
	b, err := coder.Encode(a)
	t.Log(b, err)

	//ur := User{}
	ur := make(map[string]int, 0)
	err = coder.Decode(b, &ur)
	t.Log(ur, err)
}

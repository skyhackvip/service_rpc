package codec

import (
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
	u := User{Id: 1,
		Name:   "Kavin",
		Habbit: []string{"coding", "touring"},
	}
	coder := GobCodec{}
	b, err := coder.Encode(u)
	t.Log(b, err)

	ur := User{}
	err = coder.Decode(b, &ur)
	t.Log(ur, err)
}

package util

import (
	"unsafe"
)

func StringToByte(s string) []byte {
	r := (*[2]uintptr)(unsafe.Pointer(&s))
	k := [3]uintptr{r[0], r[1], r[1]}
	return *(*[]byte)(unsafe.Pointer(&k))
}

func ByteToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

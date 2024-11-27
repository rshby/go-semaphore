package utils

import "encoding/json"

func Dump(i interface{}) string {
	return string(ToByte(i))
}

func ToByte(i interface{}) []byte {
	byte_, _ := json.Marshal(i)
	return byte_
}

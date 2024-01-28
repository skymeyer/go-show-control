package common

import (
	"encoding/binary"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadFromFile(file string, target interface{}) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, target)
}

func StringsContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func RemoveStringFromSlice(in []string, remove string) []string {
	for i, v := range in {
		if v == remove {
			return append(in[:i], in[i+1:]...)
		}
	}
	return in
}

func Uint16ToUint8(in uint16) (hi, lo uint8) {
	return uint8(in >> 8), uint8(in)
}

func Uint8ToUint16(hi, lo uint8) (out uint16) {
	return uint16(hi) << 8 & uint16(lo)
}

func PadToSize(data []byte, v interface{}) []byte {
	var (
		have = len(data)
		want = binary.Size(v)
	)
	if have > want {
		return data[:want-1]
	}
	if have < want {
		return append(data, make([]byte, want-have)...)
	}
	return data
}

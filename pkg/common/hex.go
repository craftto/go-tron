package common

import "encoding/hex"

// Has0xPrefix validates str begins with '0x' or '0X'.
func Has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

// Bytes2Hex encodes bytes as a hex string.
func Bytes2Hex(bytes []byte) string {
	encode := make([]byte, len(bytes)*2)
	hex.Encode(encode, bytes)
	return "0x" + string(encode)
}

// Hex2Bytes returns the bytes represented by the hexadecimal string s.
// s may be prefixed with "0x".
func Hex2Bytes(s string) ([]byte, error) {
	if Has0xPrefix(s) {
		s = s[2:]
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return hex2Bytes(s)
}

func bytes2Hex(d []byte) string {
	return hex.EncodeToString(d)
}

func hex2Bytes(str string) ([]byte, error) {
	return hex.DecodeString(str)
}

package aarinit

import (
	cryptrand "crypto/rand"
	"encoding/hex"
	"encoding/json"
)

type SrvRet struct {
	Ok bool //if Ok  is false, Value contains the error string
	//Value string //json string of data
	Value interface{}
}

func getSrvRetErrStr(err error) string {
	sr := SrvRet{
		false,
		err.Error(),
	}

	b, _ := json.Marshal(sr)
	return string(b)
}

func getSrvRetSuccessStr(v interface{}) string {
	sr := SrvRet{
		true,
		v,
	}

	b, _ := json.Marshal(sr)
	return string(b)
}

func getNewUuid() (string, error) {
	uuid := make([]byte, 16)
	n, err := cryptrand.Read(uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// TODO: verify the two lines implement RFC 4122 correctly
	uuid[8] = 0x80 // variant bits see page 5
	uuid[4] = 0x40 // version 4 Pseudo Random, see page 7

	return hex.EncodeToString(uuid), nil
}

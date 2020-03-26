package common

import (
	"fmt"
	"encoding/json"
)

// chaincode response结构
type chaincodeRet struct {
	Code int    // 0 success otherwise 1
	Description  string //description
}

// response 格式化消息
func GetRetByte(code int, des string) []byte {
	var r chaincodeRet
	r.Code = code
	r.Description = des

	b, err := json.Marshal(r)

	if err != nil {
		fmt.Println("marshal Ret failed")
		return nil
	}
	return b
}

// response
func GetRetString(code int, des string) string {
	var r chaincodeRet
	r.Code = code
	r.Description = des

	b, err := json.Marshal(r)

	if err != nil {
		fmt.Println("marshal Ret failed")
		return ""
	}
	return string(b[:])
}

// pObj: pointer to individual object
func NewObjectFromJsonString(jsonStr string, pObj interface{}) error {
	err := json.Unmarshal([]byte(jsonStr), pObj)

	if err != nil {
		return err
	}
	
	return nil
}

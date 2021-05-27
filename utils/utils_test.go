package utils

import (
	"errors"
	"fmt"
	"proxy-channel/utils"
	"testing"
	// "github.com/pquerna/ffjson/ffjson"
)

func Test_byte_replace(t *testing.T) {
	all := []byte("\r\nHost: 192.168.6.129:4008 \r\nConnection: Keep-Alive \r\n")
	old := []byte("Host: 192.168.6.129:4008")
	new := []byte("Host: 192.168.6.129:3008")
	ret := utils.Replace(all, old, new, -1)
	fmt.Printf("ret=[%s]\n", string(ret))
}

func Test_stringEqual_noCase(t *testing.T) {
	{
		str1 := "abcde"
		str2 := "aBCde"
		fmt.Println("case1=", IsStringEqual(str1, str2))
	}
}

func Test_ParseMultiPartFileBody(t *testing.T) {

	fileName := "multipartbody.dat"
	size, filecontent, err := GetFileContent(fileName)
	if err != nil {
		fmt.Printf("open file failed! file=%s err=%+v\n", fileName, err)
		return
	}
	_ = size
	_ = filecontent

	tmpboundary := "d9b7c92132248447f30673fcc068a48ee4c27b4dae220e5121ba00b90c35"
	retMapKv := make(map[string]string)
	ParseMultiPartFile(filecontent, tmpboundary, &retMapKv)

	fmt.Printf("last kv=%+v\n", retMapKv)

}
func Test_ParseMultiPartFile_Error(t *testing.T) {

	err := errors.New("create a error")
	fmt.Printf("errs=%s\n", err.Error()) //err:     create a error
	// fmt.Println("Error():", err.Error()) //Error(): create a error
	return

	// fileName := "error.dat"
	// size, filecontent, err := GetFileContent(fileName)
	// if err != nil {
	// 	fmt.Printf("open file failed! file=%s err=%+v\n", fileName, err)
	// 	return
	// }
	// _ = size
	// _ = filecontent

	// tmpboundary := "jcQoBazZFu9XPduYmFMM26YRGCYB9HKz"
	// retMapKv := make(map[string]string)
	// ParseMultiPartFile(filecontent, tmpboundary, &retMapKv)

	// fmt.Printf("last kv=%+v\n", retMapKv)

}

package utils

import (
	"fmt"
	"testing"
	"utils"
)

func init() {
}

func TestSerial(t *testing.T) {

	var ver uint16 = 0X18
	var pkgType uint8 = 13
	var appID uint32 = 3003521
	var ftip uint32 = 1234123
	var port uint16 = 1082
	var connID uint64 = 19971111

	clientBytes := []byte("1234")

	allbytes := SerialClientBytes(ver, pkgType, appID, ftip, port, connID, clientBytes)

	retobj := DeserialClientBytes(allbytes)
	fmt.Printf("retobj=%+v\n", *retobj)

	uip := utils.ExtractIpFromIpPort("10.0.0.1:1009")
	fmt.Println("uip=", uip)
	fmt.Println("strip=", utils.IntToIP(uip))
}

func TestBytesChannel(t *testing.T) {
	fmt.Println("aaaaaaaaaaaaaaaa")

	pool := NewBytesPool(300) // set the pool length is 3
	buf := []struct {
		b []byte
	}{
		{b: []byte("abcdef")},
		{b: []byte("123455")},
		{b: []byte("aaaabbbcccddd")},
	}

	fmt.Println("beg")
	for _, b := range buf {
		pool.Put(b.b)
	}
	fmt.Println("end")
	_ = buf
	for {
		fmt.Println(string(pool.Get()))
	}
	// for i := 0; i < len(buf); i++ {
	// 	fmt.Println(string(pool.Get()))
	// }
	// return
}

func TestExtractPort(t *testing.T) {

	ipport := "10.0.0.1:10069"
	fmt.Println("retPort=", utils.ExtractPortFromIpPort(ipport))

}

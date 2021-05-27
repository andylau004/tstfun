package utils

import (
	"encoding/binary"
	"fmt"
	"unsafe"
	// blg4go "vrv/blog4go"
)

// 目前协议头长度=21字节，不排除以后扩展

// 序列化 client 数据包
// | Ver + Type + AppID + ftip + Port + ConnID | Client Data |
func SerialClientBytes(ver uint16,
	pkgType uint8, appID uint32, ftip uint32, port uint16,
	connID uint64, clientBytes []byte) []byte {

	allsize := unsafe.Sizeof(ver) + unsafe.Sizeof(pkgType) + unsafe.Sizeof(appID) + unsafe.Sizeof(ftip) + unsafe.Sizeof(port) + unsafe.Sizeof(connID)

	allsize += (uintptr)(len(clientBytes))
	// fmt.Println("allsize2=", allsize)

	allbytes := make([]byte, allsize)

	binary.LittleEndian.PutUint16(allbytes[0:], ver)
	allbytes[2] = pkgType
	binary.LittleEndian.PutUint32(allbytes[3:], appID)
	binary.LittleEndian.PutUint32(allbytes[7:], ftip)
	binary.LittleEndian.PutUint16(allbytes[11:], port)
	binary.LittleEndian.PutUint64(allbytes[13:], connID)

	copy(allbytes[21:], clientBytes)

	return allbytes
}

// 构造ftserver ---> daodan ----> bdserver ----> linkdood server心跳包
func BuildHeartBeatPkg(ver uint16,
	pkgType uint8, appID uint32, ftip uint32, port uint16,
	connID uint64) []byte {

	allsize := unsafe.Sizeof(ver) + unsafe.Sizeof(pkgType) + unsafe.Sizeof(appID) + unsafe.Sizeof(ftip) + unsafe.Sizeof(port) + unsafe.Sizeof(connID)
	// fmt.Println("allsize2=", allsize)

	allbytes := make([]byte, allsize)

	binary.LittleEndian.PutUint16(allbytes[0:], ver)
	allbytes[2] = pkgType
	binary.LittleEndian.PutUint32(allbytes[3:], appID)
	binary.LittleEndian.PutUint32(allbytes[7:], ftip)
	binary.LittleEndian.PutUint16(allbytes[11:], port)
	binary.LittleEndian.PutUint64(allbytes[13:], connID)

	return allbytes
}

type ProtocolPKg struct {
	Ver         uint16
	PkgType     uint8
	AppID       uint32
	Ftip        uint32 // 前置服务ip，传给后端，用来标识从哪个前置服务ip传入
	Port        uint16 // 前置服务端口，传给后端，用来标识从哪个前置服务port传入
	ConnID      uint64
	ClientBytes []byte
}

func (this *ProtocolPKg) GetPkgInfo() string {
	return fmt.Sprintf("ver=%d PkgType=%d AppID=%d Ftip=%s:%d ConnID=%d len(Bytes)=%d",
		this.Ver, this.PkgType, this.AppID, IntToIP(this.Ftip), this.Port, this.ConnID, len(this.ClientBytes))
}
func (this *ProtocolPKg) GetConnKey() string {
	key := fmt.Sprintf("%d_%s_%d_%d",
		this.AppID, IntToIP(this.Ftip), this.Port, this.ConnID)
	return key
}

func DeserialClientBytes(allbytes []byte) *ProtocolPKg {
	retObj := &ProtocolPKg{}

	retObj.Ver = binary.LittleEndian.Uint16(allbytes[0:])
	fmt.Println("ver=", retObj.Ver)

	retObj.PkgType = (uint8)(allbytes[2])
	fmt.Println("pkgType=", retObj.PkgType)

	retObj.AppID = binary.LittleEndian.Uint32(allbytes[3:])
	fmt.Println("appID=", retObj.AppID)

	retObj.Ftip = binary.LittleEndian.Uint32(allbytes[7:])
	fmt.Println("ftip=", retObj.Ftip)

	retObj.Port = binary.LittleEndian.Uint16(allbytes[11:])
	fmt.Println("port=", retObj.Port)

	retObj.ConnID = binary.LittleEndian.Uint64(allbytes[13:])
	fmt.Println("connID=", retObj.ConnID)

	retObj.ClientBytes = allbytes[21:]
	fmt.Println("ClientBytes=", string(retObj.ClientBytes))

	return retObj
}

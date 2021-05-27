package main

import (
	"fmt"
)

func dummy_defs() {
	fmt.Println("")
}

// 基础扩展字段
type BaseExtend struct {
	SDKID string `json:"sdkID"`
}

type Derive1 struct {
	BaseExtend
	Number int64
}

type UnknownType struct {
	SDKID string
	Age   int32
	Val   int
}


type QueryGroupParamEx struct {
	ChangeVersion int64 `thrift:"changeVersion,1" json:"changeVersion"`
	IsDeleted     int16 `thrift:"isDeleted,2" json:"isDeleted"`
	PageSize      int16 `thrift:"pageSize,3" json:"pageSize"`
	MemberType    int8  `thrift:"memberType,4" json:"memberType,omitempty"`
}
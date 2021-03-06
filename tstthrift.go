package main

import (
	"Sample"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	// "Sample"

	"github.com/apache/thrift/lib/go/thrift"
)

func Usage() {
	fmt.Fprint(os.Stderr, "Usage of ", os.Args[0], ":\n")
	flag.PrintDefaults()
	fmt.Fprint(os.Stderr, "\n")
}

//定义服务
type Greeter struct {
}

//SayHello
func (this *Greeter) SayHello(ctx context.Context, u *Sample.User) (r *Sample.Response, err error) {
	fmt.Println("call func ,,,,,,,,,,,,,,,,,,,,say hello")
	strJson, _ := json.Marshal(u)
	return &Sample.Response{ErrCode: 0, ErrMsg: "success", Data: map[string]string{"User": string(strJson)}}, nil
}

//GetUser
func (this *Greeter) GetUser(ctx context.Context, uid int32) (r *Sample.Response, err error) {
	fmt.Println("call func ,,,,,,,,,,,,,,,,,,,,get user")
	return &Sample.Response{ErrCode: 1, ErrMsg: "user not exist."}, nil
}

func tst_thrift() {
	//命令行参数
	flag.Usage = Usage
	protocol := flag.String("P", "binary", "Specify the protocol (binary, compact, json, simplejson)")
	framed := flag.Bool("framed", false, "Use framed transport")
	buffered := flag.Bool("buffered", false, "Use buffered transport")
	addr := flag.String("addr", "localhost:9090", "Address to listen to")

	flag.Parse()

	fmt.Println("protocol=", *protocol, ", framed=", *framed,
		", buffered=", *buffered, ", addr=", *addr)

	var protocolFactory thrift.TProtocolFactory //protocol
	switch *protocol {
	case "compact":
		protocolFactory = thrift.NewTCompactProtocolFactory()
	case "simplejson":
		protocolFactory = thrift.NewTSimpleJSONProtocolFactory()
	case "json":
		protocolFactory = thrift.NewTJSONProtocolFactory()
	case "binary", "":
		protocolFactory = thrift.NewTBinaryProtocolFactoryDefault()
	default:
		fmt.Fprint(os.Stderr, "Invalid protocol specified", protocol, "\n")
		Usage()
		os.Exit(1)
	}

	var transportFactory thrift.TTransportFactory //buffered
	if *buffered {
		transportFactory = thrift.NewTBufferedTransportFactory(8192)
	} else {
		transportFactory = thrift.NewTTransportFactory()
	}

	if *framed { //framed
		transportFactory = thrift.NewTFramedTransportFactory(transportFactory)
	}

	handler := &Greeter{} //handler

	var err error
	var transport thrift.TServerTransport
	transport, err = thrift.NewTServerSocket(*addr) //transport,no secure
	if err != nil {
		fmt.Println("error running server:", err)
	}

	processor := Sample.NewGreeterProcessor(handler) //processor

	fmt.Println("Starting the simple server... on ", *addr)

	//start tcp server
	server := thrift.NewTSimpleServer4(processor, transport, transportFactory, protocolFactory)
	err = server.Serve()

	if err != nil {
		fmt.Println("error running server:", err)
	}
}

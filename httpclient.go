package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"mime/multipart"
	"net"
	"sync/atomic"
	"time"
	"utils"

	"github.com/valyala/fasthttp"
)

func dummy_httpclient() {
	fmt.Println("")
}

var gClientSocketCount int32

func GetCurTm() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func BuildFastHttpClient(hostPort string, timeout int) *fasthttp.Client {
	newClient := &fasthttp.Client{
		MaxConnsPerHost: 5000000,
		ReadTimeout:     time.Duration(timeout) * time.Second,
		Dial: func(addr string) (net.Conn, error) {
			fmt.Println("tm=", GetCurTm(), "before dail, host=", hostPort)
			tmpclient, err := fasthttp.DialTimeout(hostPort, time.Duration(timeout)*time.Second)
			fmt.Println("tm=", GetCurTm(), "after  dail, host=", hostPort)

			curSocketCount := atomic.AddInt32(&gClientSocketCount, 1)
			fmt.Println("tm=", GetCurTm(), "cur socket count=", curSocketCount)
			return tmpclient, err
		},
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	}
	fmt.Println("MaxIdleConnDuration=", newClient.MaxIdleConnDuration)
	fmt.Println("MaxConnDuration=", newClient.MaxConnDuration)
	fmt.Printf("newclient=%+v\n", *newClient)
	return newClient
}

// var gIdx int

var gChClients chan *fasthttp.Client
var gHostPort string = "127.0.0.1:7890"

func buildReqHttpPkg(inUserData string, hostport, urlpath string) (string, []byte) {
	// accToken := ""
	fileContents := []byte(inUserData)

	var bodyBuffer bytes.Buffer
	bodyWriter := multipart.NewWriter(&bodyBuffer) //返回一个设定了一个随机boundary的Writer w，并将数据写入&b
	{
		part, err := bodyWriter.CreateFormFile("file", "request") //使用给出的属性名（对应name）和文件名（对应filename）创建一个新的form-data头，part为io.Writer类型
		if err != nil {
			fmt.Println("create form file failed!", (err))
			return "", []byte("")
		}
		part.Write(fileContents) //然后将文件的内容添加到form-data头中

		err = bodyWriter.WriteField("dataType", "file") //WriteField方法调用CreateFormField，设置属性名（对应name）为"key",并在下一行写入该属性值对应的value = "val"
		if err != nil {
			fmt.Println("write field failed!", (err))
			return "", []byte("")
		}
		err = bodyWriter.Close()
		if err != nil {
			fmt.Println("body write failed!", (err))
			return "", []byte("")
		}

		bodyString := bodyBuffer.String()
		if len(bodyString) == 0 {
			fmt.Println("body buffer unexpected empty result")
			return "", []byte("")
		}
		if bodyString[0] == '\r' || bodyString[0] == '\n' {
			fmt.Println("body buffer unexpected newline")
			return "", []byte("")
		}
		// fmt.Printf("bodyString=[%s]\n", bodyString)
		return bodyWriter.FormDataContentType(), bodyBuffer.Bytes()
	}
	// fmt.Println("")
	// fmt.Printf("bodyWriter.Boundary=[%s]\n", bodyWriter.Boundary())
	// fmt.Println("")

	{ // try to parse data
		// r := multipart.NewReader(&bodyBufer, bodyWriter.Boundary())

		// part, err := r.NextPart()
		// if err != nil {
		// 	fmt.Printf("part 1: %v\n", err)
		// }
		// if g, e := part.FormName(), "myfile"; g != e {
		// 	fmt.Printf("part 1: want form name %q, got %q\n", e, g)
		// } else {
		// 	fmt.Printf("part 1: want form name %q\n", e)
		// }
		// slurp, err := ioutil.ReadAll(part)
		// if err != nil {
		// 	fmt.Printf("part 1: ReadAll: %v\n", err)
		// }
		// if e, g := string(fileContents), string(slurp); e != g {
		// 	fmt.Printf("part 1: want contents %q, got %q\n", e, g)
		// } else {
		// 	fmt.Printf("part 1: want contents %q\n", e)
		// }

		// part, err = r.NextPart()
		// if err != nil {
		// 	fmt.Printf("part 2: %v\n", err)
		// }
		// if g, e := part.FormName(), "key"; g != e {
		// 	fmt.Printf("part 2: want form name %q, got %q\n", e, g)
		// } else {
		// 	fmt.Printf("part 2: want form name %q\n", e)
		// }
		// slurp, err = ioutil.ReadAll(part)
		// if err != nil {
		// 	fmt.Printf("part 2: ReadAll: %v\n", err)
		// }
		// if e, g := "val", string(slurp); e != g {
		// 	fmt.Printf("part 2: want contents %q, got %q\n", e, g)
		// } else {
		// 	fmt.Printf("part 1: want contents %q\n", e)
		// }

		// part, err = r.NextPart() //上面的例子只有两个part
		// if part != nil || err == nil {
		// 	fmt.Printf("client expected end of parts; got %v, %v\n", part, err)
		// }
	}
	{
		// tmpurl := fmt.Sprintf("http://%s/%s", hostPort, AgentDataPath)

		// newClient := BuildFastHttpClient(hostport, 5)

		// contentType := bodyWriter.FormDataContentType()
		// request := fasthttp.AcquireRequest()
		// fmt.Println("contentType=", contentType)

		// request.Header.SetContentType(contentType)
		// request.SetBody(bodyBuffer.Bytes())
		// request.Header.SetMethod("POST")
		// request.Header.Set("token", accToken)
		// request.SetRequestURI(url)

		// resp := fasthttp.AcquireResponse()
		// err := newClient.Do(request, resp)
		// if err != nil {
		// 	errMsg := fmt.Sprintf("client request agent path failed! err=%+v url=%s accToken=%s contentType=%s body=%s\n",
		// 		err.Error(), url, accToken, contentType, string(bodyBuffer.Bytes()))
		// 	return errCheckFunc(errMsg)
		// }
		// if resp.StatusCode() != fasthttp.StatusOK {
		// 	errMsg := fmt.Sprintf("status code invalid, [err=%v] code=[%d]", err, resp.StatusCode())
		// 	return errCheckFunc(errMsg)
		// }
	}

	{ // self build http header + body
		// tmpBoundary := RandStringRunes(60)
		// httpHeader := fmt.Sprintf("POST %s HTTP/1.1\r\nUser-Agent: fasthttp\r\nHost: %s\r\n", urlpath, hostport)
		// contentType := fmt.Sprintf("Content-Type: multipart/form-data; boundary=%s\r\n", tmpBoundary)
		// token := fmt.Sprintf("token: %s\r\n\r\n", accToken)

		// begBoundary := fmt.Sprintf("--%s\r\n", tmpBoundary)
		// endBoundary := fmt.Sprintf("--%s--\r\n", tmpBoundary)

		// var httpBody, allHttp string

		// content_Disp2 := fmt.Sprintf("Content-Disposition: form-data; name=\"dataType\"\r\n\r\n")
		// httpBody += begBoundary
		// httpBody += content_Disp2
		// httpBody += "file\r\n"

		// content_Disp := fmt.Sprintf("Content-Disposition: form-data; name=\"file\"; filename=\"request\"\r\n")
		// content_Type1 := fmt.Sprintf("Content-Type: application/octet-stream\r\n\r\n")

		// httpBody += (begBoundary + content_Disp + content_Type1)
		// httpBody += (string(protoPkg) + string(clientPkg) + ("\r\n"))
		// httpBody += endBoundary

		// {
		// 	contentLength := fmt.Sprintf("Content-Length: %d\r\n", len(httpBody))
		// 	httpHeader += contentType
		// 	httpHeader += contentLength
		// 	httpHeader += token
		// }

		// allHttp = httpHeader + httpBody
		// logger.Info("send to dandao httpPkg", zap.String("allHttp", allHttp))
	}
	return "", []byte("")
}

func socketclient_workImpl(thrdIdx int) {
	newSocket := gArrSocketClients[thrdIdx]
	fmt.Printf("thrd [%d] get sock=%d\n", thrdIdx, newSocket)

	var err error
	path := "/mixin/v1/pingpong"
	accToken := ""
	workWheel := 1
	for {
		fmt.Println("")
		fmt.Println(" -----------------------------------beg----------------------------------- ")
		reqString := fmt.Sprintf("client [%d] req 666", thrdIdx)

		tmpBoundary := RandStringRunes(60)
		httpHeader := fmt.Sprintf("POST %s HTTP/1.1\r\nUser-Agent: fasthttp\r\nHost: %s\r\n", path, gHostPort)
		contentType := fmt.Sprintf("Content-Type: multipart/form-data; boundary=%s\r\n", tmpBoundary)
		token := fmt.Sprintf("token: %s\r\n\r\n", accToken)

		begBoundary := fmt.Sprintf("--%s\r\n", tmpBoundary)
		endBoundary := fmt.Sprintf("--%s--\r\n", tmpBoundary)

		var httpBody, allHttp string

		content_Disp2 := fmt.Sprintf("Content-Disposition: form-data; name=\"dataType\"\r\n\r\n")
		httpBody += begBoundary
		httpBody += content_Disp2
		httpBody += "file\r\n"

		content_Disp := fmt.Sprintf("Content-Disposition: form-data; name=\"file\"; filename=\"request\"\r\n")
		content_Type1 := fmt.Sprintf("Content-Type: application/octet-stream\r\n\r\n")

		httpBody += (begBoundary + content_Disp + content_Type1)
		httpBody += (reqString + ("\r\n"))
		httpBody += endBoundary

		{
			contentLength := fmt.Sprintf("Content-Length: %d\r\n", len(httpBody))
			httpHeader += contentType
			httpHeader += contentLength
			httpHeader += token
		}

		allHttp = httpHeader + httpBody
		fmt.Printf("send to dandao httpPkg, allHttp=%s\n", allHttp)

		err = Send(newSocket, utils.Str2bytes(allHttp), "send dandao")
		if err != nil {
			fmt.Printf("conn send dandao failed! err=%+v\n", err)
			return
		}

		var recvBuff [64 * utils.KB]byte
		recvSize, err := newSocket.Read(recvBuff[:])
		if err != nil { // 跟upstream连接关闭，需要清理资源
		} else {
			respBody := recvBuff[0:recvSize]
			fmt.Printf("thrd [%d] resp[%d] respBody=%s\n", thrdIdx, workWheel, string(respBody))

			workWheel++
			fmt.Println(" -----------------------------------end----------------------------------- ")
		}

		time.Sleep(time.Second)
	}

}

// 使用 fast http client 请求服务器方式
func httpclient_workImpl(thrdIdx int) {
	newClient := gArrHttpClients[thrdIdx]
	fmt.Printf("thrd [%d] get clt=%p\n", thrdIdx, newClient)

	// tr := &http.Transport{
	// 	MaxIdleConns:        1024,
	// 	MaxIdleConnsPerHost: 1024,
	// 	TLSHandshakeTimeout: 10 * time.Second,
	// 	Dial: (&net.Dialer{
	// 		Timeout:   75 * time.Second,
	// 		KeepAlive: 75 * time.Second,
	// 	}).Dial,
	// }

	// client = &http.Client{Transport: tr}

	accToken := ""
	workWheel := 1
	for {
		fmt.Println("")
		fmt.Println(" -----------------------------------beg----------------------------------- ")

		url := fmt.Sprintf("http://%s%s", gHostPort, "/mixin/v1/pingpong")

		// newClient := BuildFastHttpClient(hostPort, timeout)
		// newClient := <-gChClients

		clientPkg := fmt.Sprintf("client [%d] req 666", thrdIdx)

		contentType, bodyBytes := buildReqHttpPkg(clientPkg, gHostPort, "/mixin/v1/pingpong")
		// fmt.Printf("thrd [%d] build body=%s\n", thrdIdx, string(bodyBytes))

		request := fasthttp.AcquireRequest()
		request.SetRequestURI(url)

		resp := fasthttp.AcquireResponse()

		request.Header.SetContentType(contentType)
		// fmt.Println("contentType=", contentType)
		request.Header.SetMethod("POST")
		request.Header.Set("token", accToken)
		request.SetBody(bodyBytes)

		// tmBeg := time.Now()
		// fmt.Println("tm=", tmBeg, " before do req")
		// fmt.Println("req.URI().Host()=", string(request.URI().Host()), ", hasBodyStream=", request.IsBodyStream())
		var err error
		// err := newClient.DoVrv(thrdIdx, request, resp)
		// fmt.Println("tm=", time.Now(), " after do req")
		if err != nil {
			errMsg := fmt.Sprintf("client request agent path failed! err=%+v url=%s accToken=%s contentType=%s body=%s\n",
				err.Error(), url, accToken, contentType, string(bodyBytes))
			fmt.Printf("thrd [%d] do failed! err=%s\n", thrdIdx, errMsg)
		}
		if resp.StatusCode() != fasthttp.StatusOK {
			errMsg := fmt.Sprintf("status code invalid, [err=%v] code=[%d]", err, resp.StatusCode())
			fmt.Printf("thrd [%d] status failed! err=%s\n", thrdIdx, errMsg)
		}
		fmt.Printf("tm= %+v, thrd [%d] Got resp[%d] respCode=%d respBody=%s\n", GetCurTm(), thrdIdx, workWheel, resp.StatusCode(), string(resp.Body()))

		// byteBody := resp.Body()
		// fmt.Printf("[thrd:%d] resp body=%+v\n", thrdIdx, string(byteBody))
		// gChClients <- newClient

		workWheel++
		fmt.Println(" -----------------------------------en----------------------------------- ")
		time.Sleep(1 * time.Second)
	}

}

var gArrHttpClients []*fasthttp.Client
var gArrSocketClients []net.Conn
var gWorkCount int

func tstHttpClient() {
	var timeout int = 20
	var hostPort string = "127.0.0.1:7890"
	gWorkCount = 2

	// gChClients = make(chan *fasthttp.Client, 1024)
	tag := 0
	if 1 == tag {

		for i := 0; i < gWorkCount; i++ {
			conn, err := net.Dial("tcp", hostPort)
			if err != nil {
				fmt.Printf("dial host: %s failed! err=%+v\n", hostPort, err)
			} else {
				gArrSocketClients = append(gArrSocketClients, conn)
			}

		}

		for i := 0; i < gWorkCount; i++ {
			go socketclient_workImpl(i)
			time.Sleep(time.Microsecond * 500)
		}

	} else {

		for i := 0; i < gWorkCount; i++ {
			newClient := BuildFastHttpClient(hostPort, timeout)

			gArrHttpClients = append(gArrHttpClients, newClient)

			fmt.Printf("new http client=%p\n", newClient)
		}

		for i := 0; i < gWorkCount; i++ {
			go httpclient_workImpl(i)
			time.Sleep(time.Microsecond * 500)
		}

	}

	select {}
}

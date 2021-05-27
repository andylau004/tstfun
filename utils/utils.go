package utils

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
	"unsafe"

	blg4go "vrv/blog4go"

	toml "github.com/pelletier/go-toml"

	"github.com/pquerna/ffjson/ffjson"
	"github.com/valyala/fasthttp"
	// blg4go "github.com/YoungPioneers/blog4go"
)

const (
	ONE_PIECE_SIZE = 4 * 1024
)

const (
	FIRST_PKG = iota
	DATA_TYPE
	CLOSE_TYPE
	UNKNOWN_TYPE
)

var (
	errBlog4GoInit = errors.New("error blog4go init")
	nullBytes      = make([]byte, 0)
)

type ByteSize uint64

const (
	_           = iota             // ignore first value by assigning to blank identifier
	KB ByteSize = 1 << (10 * iota) // 1 << (10*1)
	MB                             // 1 << (10*2)
	GB                             // 1 << (10*3)
	TB                             // 1 << (10*4)
	PB                             // 1 << (10*5)
	EB                             // 1 << (10*6)
	// ZB                             // 1 << (10*7)
	// YB                             // 1 << (10*8)
)

var g_u64MockId uint64

// func CountElapse(msg string) func() {

// 	start := time.Now()
// 	blg4go.Infof("%s Begin... ", msg)

// 	return func() {
// 		blg4go.Infof("%s End... Elapse:(%s)", msg, time.Since(start))
// 	}
// }

// --------------------------------------------
// for calculate random 60 bytes,, for boundary
func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// 计算随机http boundary
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// --------------------------------------------

func Str2bytes(s string) []byte {

	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}

	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func GetGoId_Str() string {
	sRet := ("GOId=" + fmt.Sprint(GoID()) + " ")
	return sRet
}

func GoID() int {
	return 0

	// var buf [64]byte
	// n := runtime.Stack(buf[:], false)

	// idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	// GoroutineId, err := strconv.Atoi(idField)
	// if err != nil {
	// 	panic(fmt.Sprintf("cannot get goroutine id: %v", err))
	// }

	// return GoroutineId
}

func GetGoId_StrEx() string {
	sRet := ("GOId=" + fmt.Sprint(GoIDEx()) + " ")
	return sRet
}
func GoIDEx() uint64 {
	return atomic.AddUint64(&g_u64MockId, uint64(1))
}

func IsFileExist(sFileName string) bool {
	var _, err = os.Stat(sFileName)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func GetCurStrTime() string {
	return (time.Now().Format("2006-01-02 15:04:05") + " ")
}

func Int2String(val int) string {
	return strconv.Itoa(val)
}
func Int64String(val int64) string {
	return strconv.FormatInt(val, 10)
}

func Str2int64(str string) (int64, error) {
	i64, err := strconv.ParseInt(str, 10, 64)
	return i64, err
}

func Str2int(str string) (int, error) {
	i, err := strconv.Atoi(str)
	return i, err
}

// 分片拼接完整文件；
// func BuildSliceFile(userFolder string, fileId int64, Burst int) (string, bool) {

// 	// lastfile := userFolder + "/" + utils.Int64String(reqObj.FileId)
// 	lastfile := userFolder + "/" + Int64String(fileId)
// 	// fmt.Println("userfolder=", userFolder, ", fileId=", fileId, ", Burst=", Burst)

// 	var bok bool = true
// 	for idx := 1; idx <= Burst; idx++ {
// 		tmpname := lastfile + "-" + Int2String(idx)
// 		// fmt.Println("tmpname=", tmpname)

// 		// read tmpname, append lastfile
// 		// fmt.Println("before get content")
// 		_, sliceContent, err := GetFileContent(tmpname)
// 		// fmt.Println("after get content")

// 		if err != nil {
// 			fmt.Println("get slice file bytes failed! file=", tmpname, " err=", err)
// 			blg4go.Error("get slice file bytes failed! file=", tmpname, " err=", err)
// 			bok = false
// 			break
// 		} else {
// 			_, err = AppendFile(lastfile, sliceContent)
// 			if err != nil {
// 				fmt.Println("apend lastfile failed! file=", lastfile, " err=", err)
// 				blg4go.Error("apend lastfile failed! file=", lastfile, " err=", err)
// 				bok = false
// 				break
// 			}
// 		}
// 	} // traverse slice file

// 	if bok {
// 		fmt.Println("lastfile build successed! lastfile=", lastfile)
// 		blg4go.Info("lastfile build successed! lastfile=", lastfile)
// 		return lastfile, bok
// 	}
// 	return "", false
// }

// If the file doesn't exist, create it, or append to the file
func AppendFile(fileName string, wrBytes []byte) (int, error) {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("open file failed! err=", err)
		return -1, err
	}
	if _, err := f.Write(wrBytes); err != nil {
		fmt.Println("write file failed! err=", err)
		return -1, err
	}
	// fmt.Println("write file success! wrBytes=", string(wrBytes))

	if err := f.Close(); err != nil {
		fmt.Println("close file failed! err=", err)
		return -1, err
	}
	return len(wrBytes), nil
}

func SaveFile(fileName string, wrBytes []byte) (int, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println("open file failed. err=", err.Error())
		return 0, err
	}
	defer file.Close()
	return file.Write(wrBytes)
}

func SaveZeroFile(fileName string, wrBytes []byte) (int, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Println("open file failed. err=", err.Error())
		return 0, err
	}
	defer file.Close()
	return file.Write(wrBytes)
}

func GetFileSize(fileFullPath string) int64 {
	fileInfo, err := os.Stat(fileFullPath)
	if err != nil {
		return -1
	}

	filesize := fileInfo.Size()
	// fmt.Println("fileName=", fileFullPath, ", fileSize=", filesize) //返回的是字节
	return filesize
}

func GetFileContent(fileName string) (int, []byte, error) {
	fp, err := os.Open(fileName) // 获取文件指针
	if err != nil {
		fmt.Println("open file faild!!! filename=", fileName, ", err=", err)
		return -1, nil, err
	}
	defer fp.Close()

	filesize := GetFileSize(fileName)
	if filesize <= 0 {
		return -1, nil, err
	}
	readBuffer := make([]byte, filesize)

	var allsize int
	for {
		// 注意这里要取bytesRead, 否则有问题
		bytesRead, err := fp.Read(readBuffer) // 文件内容读取到buffer中
		// fmt.Println("bytesRead=", bytesRead, ", err=", err)

		// nSendLen, err = conn.Write(buffer[:bytesRead])
		// checkError(err)
		// fmt.Println("nSendLen=", nSendLen, ", err=", err)
		allsize += bytesRead

		if err != nil {
			if err == io.EOF {
				err = nil
				break
			} else {
				fmt.Println("some error happened")
				return allsize, nil, err
			}
		}

	}
	return allsize, readBuffer[0:allsize], nil
}

//调用os.MkdirAll递归创建文件夹
func CreateFullFolder(filePath string) error {
	if !IsExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		return err
	}
	return nil
}

// 判断所给路径文件/文件夹是否存在(返回true是存在)
func IsExist(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func TrimString(s, sep string) []string {
	arrstr := strings.Split(s, sep)
	if len(arrstr) > 0 {

	}
	return arrstr
}

func SizeToMB(size int64) int {
	return (int)(size / (int64)(MB))
}
func SizeToGB(size int64) int {
	return (int)(size / (int64)(GB))
}

func GetSerial(obj interface{}) string {
	retBuffer, _ := ffjson.Marshal(obj)
	return string(retBuffer)
}

type LogCfgBase interface {
	GetBlgCfgFile() string
}

func ConfigDecode(fname string, cfgobj interface{}) error {
	if fname == "" {
		return nil
	}
	fd, err := os.Open(fname)
	if err != nil {
		fmt.Println("open config file error, ", err)
		return err
	}
	defer fd.Close()

	decoder := toml.NewDecoder(fd)
	if decoder == nil {
		return errors.New("decode failed!")
	}
	err = decoder.Decode(cfgobj)
	if err != nil {
		fmt.Println("decode config error ", err)
		return err
	}

	t, ok := cfgobj.(LogCfgBase)
	if ok && t != nil {

		cfgFile := t.GetBlgCfgFile()
		errT := initBlog4Go(cfgFile)
		// fmt.Println("blgcfgfile=", cfgFile)

		if errT != nil {
			tmpname := fmt.Sprintf("blgCfg:%s", cfgFile)
			return errors.New("log init failed, name:" + tmpname)
		}
		return nil
	} else {
		return errors.New("log init failed, maybe no impl LogCfgBase interface!")
	}
	// return nil
}

func initBlog4Go(fname string) error {
	err := blg4go.NewWriterFromConfigAsFile(fname)
	if nil != err {
		fmt.Println("init log fun failed! err=", err.Error(), ", fname=", fname)
		return errBlog4GoInit
	}
	return nil
}

func IsGetMethod(ctx *fasthttp.RequestCtx) bool {
	return strings.EqualFold(string(ctx.Method()), "GET")
}
func IsPostMethod(ctx *fasthttp.RequestCtx) bool {
	return strings.EqualFold(string(ctx.Method()), "POST")
}

func Md5(str string) []byte {
	has := md5.Sum([]byte(str))
	strMd5 := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return []byte(strMd5)
}

func String2int64(str string) int64 {
	if n, err := strconv.Atoi(str); err == nil {
		return int64(n)
	}
	return -1
}

func Int2string(num int64) string {
	return strconv.FormatInt(num, 10)
}

func BuildFastHttpClient(hostPort string, timeout int) *fasthttp.Client {
	newClient := &fasthttp.Client{
		MaxConnsPerHost: 5000000,
		ReadTimeout:     time.Duration(timeout) * time.Second,
		Dial: func(addr string) (net.Conn, error) {
			return fasthttp.DialTimeout(hostPort, time.Duration(timeout)*time.Second)
		},
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return newClient
}

// 封装执行POST请求，给指定url，返回 < response, err, status code >
func ExecPostHttp(hostPort, url, reqHttpBody, contentType string, extraHeader map[string]string, timeout int) ([]byte, error, int) {
	var err error
	newClient := BuildFastHttpClient(hostPort, timeout)

	request := fasthttp.AcquireRequest()
	request.SetRequestURI(url)
	request.Header.SetMethod("POST")
	request.Header.Set("Content-Type", contentType)
	for k, v := range extraHeader {
		request.Header.Set(k, v)
	}
	request.SetBodyString(reqHttpBody)

	resp := fasthttp.AcquireResponse()
	if err = newClient.Do(request, resp); err != nil {
		fmt.Printf("do request failed. [err=%v]\n", err)
		blg4go.Errorf("do request failed. [err=%v] [url=%v]", err, url)
		return nullBytes, err, -1
	}
	if resp.StatusCode() != fasthttp.StatusOK {
		blg4go.Errorf("status code invalid, code=%d. [err=%v]", resp.StatusCode(), err)
		return nullBytes, errors.New("status invalid, statuscode=" + strconv.Itoa(resp.StatusCode())), -1
	}
	byteBody := resp.Body()
	// blg4go.Infof("resp body=%+v", string(byteBody))

	return (byteBody), nil, 0
}

func errCheckFunc(errMsg string) ([]byte, error) {
	fmt.Println(errMsg)
	blg4go.Error(errMsg)
	return []byte(""), errors.New(errMsg)
}

// socket 实现http post 给单导
// 连接对象保存，需要复用 传数据给内网
func SocketPostMultipartFile(inConn net.Conn, hostPort, path, accToken string, tcpPkg []byte, timeout int) (net.Conn, []byte, error) {
	var conn net.Conn
	var err error

	if inConn == nil { // 没有，则创建upload连接；有，则复用；
		conn, err = net.Dial("tcp", hostPort)
		if err != nil {
			fmt.Printf("dial host: %s failed! err=%+v\n", hostPort, err)
			blg4go.Errorf("dial host: %s failed! err=%+v", hostPort, err)
			return nil, []byte(""), err
		}
		inConn = conn
	}

	tmpBoundary := RandStringRunes(60)
	httpHeader := fmt.Sprintf("POST %s HTTP/1.1\r\nUser-Agent: fasthttp\r\nHost: %s\r\n", path, hostPort)
	contentType := fmt.Sprintf("Content-Type: multipart/form-data; boundary=%s\r\n", tmpBoundary)
	token := fmt.Sprintf("Token: %s\r\n\r\n", "")

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
	httpBody += (string(tcpPkg) + ("\r\n"))

	httpBody += endBoundary

	{
		contentLength := fmt.Sprintf("Content-Length: %d\r\n", len(httpBody))
		httpHeader += contentType
		httpHeader += contentLength
		httpHeader += token
	}

	allHttp = httpHeader + httpBody
	blg4go.Infof("allHttp=%s", allHttp)

	err = Send(inConn, Str2bytes(allHttp), "send dandao")
	if err != nil {
		fmt.Printf("conn send dandao failed! err=%+v\n", err)
		blg4go.Errorf("conn send dandao failed! err=%+v\n", err)
		return nil, []byte(""), err
	}

	return inConn, []byte(""), err
}

func PostMulitiPartFile(hostPort, url, accToken string, fileContents []byte, timeout int) ([]byte, error) {

	var bodyBuffer bytes.Buffer
	bodyWriter := multipart.NewWriter(&bodyBuffer) //返回一个设定了一个随机boundary的Writer w，并将数据写入&b
	{
		part, err := bodyWriter.CreateFormFile("file", "request") //使用给出的属性名（对应name）和文件名（对应filename）创建一个新的form-data头，part为io.Writer类型
		if err != nil {
			blg4go.Errorf("create form file failed, err=%+v", err)
			return nullBytes, err
		}
		part.Write(fileContents) //然后将文件的内容添加到form-data头中

		err = bodyWriter.WriteField("dataType", "file") //WriteField方法调用CreateFormField，设置属性名（对应name）为"key",并在下一行写入该属性值对应的value = "val"
		if err != nil {
			blg4go.Errorf("write field failed, err=%+v", err)
			return nullBytes, err
		}
		err = bodyWriter.Close()
		if err != nil {
			blg4go.Errorf("body write failed, err=%+v", err)
			return nullBytes, err
		}

		bodyString := bodyBuffer.String()
		if len(bodyString) == 0 {
			return errCheckFunc("body buffer unexpected empty result")
		}
		if bodyString[0] == '\r' || bodyString[0] == '\n' {
			return errCheckFunc("body buffer unexpected newline")
		}
		fmt.Printf("bodyString=[\n%s]\n", bodyString)
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

	newClient := BuildFastHttpClient(hostPort, timeout)

	contentType := bodyWriter.FormDataContentType()
	request := fasthttp.AcquireRequest()
	// fmt.Println("contentType=", contentType)

	request.Header.SetContentType(contentType)
	request.SetBody(bodyBuffer.Bytes())
	request.Header.SetMethod("POST")
	request.Header.Set("token", accToken)
	request.SetRequestURI(url)

	resp := fasthttp.AcquireResponse()
	err := newClient.Do(request, resp)
	if err != nil {
		errMsg := fmt.Sprintf("client request agent path failed! err=%+v url=%s accToken=%s contentType=%s body=%s\n",
			err.Error(), url, accToken, contentType, string(bodyBuffer.Bytes()))
		return errCheckFunc(errMsg)
	}
	if resp.StatusCode() != fasthttp.StatusOK {
		errMsg := fmt.Sprintf("status code invalid, [err=%v] code=[%d]", err, resp.StatusCode())
		return errCheckFunc(errMsg)
	}
	// fmt.Println("client resp statusCode=", resp.StatusCode())

	return resp.Body(), nil
}

func GetAllHeaderKv(ctx *fasthttp.RequestCtx) map[string]string {
	retMap := make(map[string]string)

	ctx.Request.Header.VisitAll(func(key, value []byte) {
		fmt.Println("key=", string(key), ", value=", string(value))
		retMap[string(key)] = string(value)
	})

	return retMap
}

func GetHeaderValByKey(ctx *fasthttp.RequestCtx, targetKey string) string {
	var retVal string

	ctx.Request.Header.VisitAll(func(key, value []byte) {
		if IsStringEqual(targetKey, string(key)) {
			retVal = string(value)
		}
	})

	return retVal
}

func IsStringEqual(str1, str2 string) bool {
	return strings.EqualFold(str1, str2)
}

func ParseMultiPartFile(byteBody []byte, boundary string, retMapKv *map[string]string) error {
	var bodyBuffer bytes.Buffer
	bodyBuffer.Write(byteBody)

	r := multipart.NewReader(&bodyBuffer, boundary)

	var part *multipart.Part
	var err error
	for {

		part, err = r.NextPart()
		if err != nil {
			errMsg := fmt.Sprintf("%s", err.Error())
			if errMsg == "EOF" {
				return nil
			}
			fmt.Printf("next part failed! err=%+v\n", err)
			blg4go.Errorf("next part failed! err=%+v", err)
			break
		}
		if part == nil {
			break
		}
		fmt.Printf("part=%+v\n", *part)

		formName := part.FormName()

		slurp, err := ioutil.ReadAll(part)
		if err != nil {
			fmt.Printf("part 1: ReadAll: %v\n", err)
		}
		// fmt.Println("datatype value=", string(slurp))

		(*retMapKv)[string(formName)] = string(slurp)
		_ = formName
		// if formName == "dataType" {
		// 	fmt.Println("datatype value=", string(slurp))
		// 	blg4go.Infof("datatype value=%s", string(slurp))
		// }

		// if formName == "file" {
		// 	fmt.Println("file value=", string(slurp))
		// 	blg4go.Infof("file value=%s", string(slurp))
		// 	return slurp
		// }
	}

	return err
}

func GetCost(tStart time.Time) time.Duration {
	return time.Now().Sub(tStart)
}

func Ip2long(ipAddr string) (uint32, error) {
	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return 0, errors.New("wrong ipAddr format")
	}
	ip = ip.To4()
	// return binary.BigEndian.Uint32(ip), nil
	return binary.LittleEndian.Uint32(ip), nil
}

// func Long2ip(ipLong uint32) string {
// 	ipByte := make([]byte, 4)
// 	binary.BigEndian.PutUint32(ipByte, ipLong)
// 	ip := net.IP(ipByte)
// 	return ip.String()
// }

func IntToIP(ip uint32) string {
	result := make(net.IP, 4)
	result[0] = byte(ip)
	result[1] = byte(ip >> 8)
	result[2] = byte(ip >> 16)
	result[3] = byte(ip >> 24)
	return result.String()
}

// "ip:port"  ---> (uint32)ip
func ExtractIpFromIpPort(ipport string) uint32 {
	findPos := strings.Index(ipport, ":")
	if findPos > 0 {
		tmpip := ipport[0:findPos]
		retIp, _ := Ip2long(string(tmpip))
		return retIp
	}
	return 0
}

// "ip:port"  ---> (uint16)port
func ExtractPortFromIpPort(ipport string) uint16 {
	findPos := strings.Index(ipport, ":")
	if findPos > 0 {
		tmpPort := ipport[findPos+1:]

		// retPort, _ := strconv.ParseUint(tmpPort, 10, 32)
		retPort, err := strconv.ParseInt(tmpPort, 10, 32)
		if err != nil {
			fmt.Println("parse int failed! err=", err)
			return 0
		}
		return (uint16)(retPort)
	}
	fmt.Println("parse int failed! findPos invalid, findPos=", findPos, ", ipport=", ipport)
	return 0
}

// caseInsensitiveCompare does a case insensitive equality comparison of
// two []byte. Assumes only letters need to be matched.
func CaseInsensitiveCompare(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i]|0x20 != b[i]|0x20 {
			return false
		}
	}
	return true
}

func Send(conn net.Conn, buf []byte, descri string) error {
	for send := 0; send < len(buf); {
		s, err := conn.Write(buf[send:])
		if err != nil {
			if buf != nil {
				fmt.Printf("tcp.go send failed! err=%+v descri=%s\n", err, descri)
			}
			return err
		}
		send += s
	}
	return nil
}

func Recv(conn net.Conn, buf []byte, description string) error {
	length := len(buf)
	for recv := 0; recv < length; {
		n, err := conn.Read(buf[recv:])
		if err != nil {
			fmt.Println("tcp.go recv failed, err=", err)
			// log.Error("tcp.go recv error ", err, " length=", length, " description=", description)
			return err
		}
		recv += n
	}
	return nil
}

// func AppendFile(fileName string, content []byte) error {
// 	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0666)
// 	if err != nil {
// 		fmt.Printf("append file: %s failed! err=%+v\n", fileName, err)
// 		return err
// 	}
// 	// 查找文件末尾的偏移量
// 	n, _ := f.Seek(0, os.SEEK_END)
// 	// 从末尾的偏移量开始写入内容
// 	_, err = f.WriteAt([]byte(content), n)
// 	defer f.Close()
// }

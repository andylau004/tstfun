package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"reflect"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
	"utils"

	"github.com/pquerna/ffjson/ffjson"
	"github.com/samuel/go-zookeeper/zk"

	blg4go "vrv/blog4go"
)

const (
	TmpBuf = `{
		"Name": "big brother",
		"Number": 0,
		"sdkID": "1943"
		}`

	QueryGroupJson = `{
			"Name": "big brother",
			"Number": 0,
			"sdkID": "1943"
			}`

	MaxTcpCount = 30000
	HttpPort    = 10077
)

var (
	gParamAssignMap map[int][]interface{}
)

type Att struct {
	Message *string
}

type CallBack interface {
	getName() string
	BaseCall
}

type BaseCall interface {
	doSomething()
}

type User struct {
	name string
}

func (user User) doSomething() {
	fmt.Println("something to do, name=", user.name)
}

func (user User) getName() string {
	return user.name
}

func hasinterface(intf interface{}) {
	switch t := intf.(type) {
	case User:
		fmt.Println("find interface = User, name=", t.name)
	default:
		fmt.Println("unknown interface")
	}
}

type TestInfo struct {
	Value *int64
	Name  *string
}

func GetNewInfo() *TestInfo {
	retObj := &TestInfo{}

	{
		var tmpval int64
		tmpval = 9981
		retObj.Value = &tmpval

		fmt.Printf("Value addr=%p\n", retObj.Value)
	}
	{
		var tmpname string
		tmpname = "kan tou bu shuo tou..."
		retObj.Name = &tmpname

		fmt.Printf("Name addr=%p\n", retObj.Name)
	}
	return retObj
}

type QueryGroupParam struct {
	ChangeVersion int64 `thrift:"changeVersion,1" json:"changeVersion"`
	IsDeleted     int16 `thrift:"isDeleted,2" json:"isDeleted"`
	PageSize      int16 `thrift:"pageSize,3" json:"pageSize"`
	MemberType    int8  `thrift:"memberType,4" json:"memberType,omitempty"`
}

func initBlog4Go() {
	err := blg4go.NewWriterFromConfigAsFile("blog4go_config.xml")
	if nil != err {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("init blg4go successfully.")
	blg4go.Info("init blg4go successfully.")
}

func init() {
	rand.Seed(time.Now().UnixNano())
	initBlog4Go()
}

type UserEx struct {
	Name  string
	Email string
}

func (u *UserEx) Notify() error {
	log.Printf("UserEx: Sending User Email To %s<%s>\n",
		u.Name,
		u.Email)
	return nil
}

type Notifier interface {
	Notify() error
}

func SendNotification(notify Notifier) error {
	return notify.Notify()
}

func Test3() {
	user := &UserEx{
		Name:  "AriesDevil",
		Email: "ariesdevil@xxoo.com",
	}

	SendNotification(user)
}
func TestZk() {
	defer fmt.Println("testzk done ...")
	c, _, err := zk.Connect([]string{"172.16.8.155:11100"}, time.Second) //*10)
	if err != nil {
		panic(err)
	}
	children, stat, ch, err := c.ChildrenW("/vrv/im/service/dbconfigServer/2.0")
	if err != nil {
		fmt.Printf("childrenw ret, err=%+v\n", err)
		return
		panic(err)
	}
	fmt.Printf("children=%+v stat=%+v\n", children, stat)
	e := <-ch
	fmt.Printf("e=%+v\n", e)
}

func Tst1() {
	TestZk()
	return

	Test3()
	return

	blg4go.Errorf("test out errorf...")

	blg4go.Info("test out info")
	blg4go.Warn("test out warn")

	time.Sleep(2 * time.Second)
	return

	{
		groupIDs := []string{TmpBuf, QueryGroupJson}
		fmt.Printf("groupIDs=%+v\n", groupIDs)
		return
	}
	{
		var tstObj QueryGroupParamEx
		tstObj.ChangeVersion = 1943

		v := reflect.ValueOf(tstObj)
		fmt.Println("v=", v)
		fmt.Println("v Type:", v.Type())
		fmt.Println("v CanSet:", v.CanSet())

		return
	}
	{
		var tstObj QueryGroupParam
		tstObj.ChangeVersion = 1943

		fmt.Printf("tstObj=%+v\n", tstObj)

		err1 := ffjson.Unmarshal([]byte(QueryGroupJson), &tstObj)
		fmt.Printf("err1=%+v tstObj=%+v\n", err1, tstObj)

		return
	}
	{
		tmpobj := GetNewInfo()

		fmt.Printf("val=%d addr=%p\n", *tmpobj.Value, tmpobj.Value)
		fmt.Printf("name=%s addr=%p\n", *tmpobj.Name, tmpobj.Name)

		return
	}
	{
		var ListGroups []*User

		u1 := &User{name: "usa no3"}
		u2 := &User{name: "china no1"}
		ListGroups = append(ListGroups, u1)
		ListGroups = append(ListGroups, u2)
		fmt.Printf("listgourps = %+v\n", ListGroups)

		for _, obj := range ListGroups {
			fmt.Printf("obj=%+v\n", *obj)
		}
		return
	}
	{
		// var tmpVal *User
		var arrInts []*User
		fmt.Println("len=", len(arrInts))
		for _, obj := range arrInts {
			fmt.Println("aaaaaaaaaa, obj=", obj)
		}
		// fmt.Println("len(nil)=", len(tmpVal))
		return
	}

	{
		var basObj BaseExtend

		err1 := ffjson.Unmarshal([]byte(TmpBuf), &basObj)
		fmt.Printf("err1=%+v basObj=%+v\n", err1, basObj)
		return

		fmt.Println("tmpbuf=", TmpBuf)

		var fileStat UnknownType
		err := ffjson.Unmarshal([]byte(TmpBuf), &fileStat)

		fmt.Printf("err=%+v fileStat=%+v\n", err, fileStat)

		return

		var d1 Derive1
		retBuffer, _ := ffjson.Marshal(&d1)
		fmt.Println(" retBuffer =", string(retBuffer))
		// fmt.Println(d1)
		return
	}

	{
		arrInts := make([]int, 0, 1)
		arrInts = append(arrInts, 1)
		arrInts = append(arrInts, 3)
		arrInts = append(arrInts, 5)
		arrInts = append(arrInts, 7)
		fmt.Printf("arrData=%+v\n", arrInts)
		return
	}
	{
		arrData := make([]User, 3, 10)
		fmt.Printf("arrData=%+v\n", arrData)
	}
	return

	obj := User{name: "key le"}

	var base BaseCall
	base = BaseCall(obj)
	base.doSomething()

	hasinterface(base)
}

func bar(count, size int) string {
	str := ""
	for i := 0; i < size; i++ {
		if i < count {
			str += "="
		} else {
			str += " "
		}
	}
	return str
}

func tstparseJson() {
	// gParamAssignMap

	gParamAssignMap = make(map[int][]interface{})

	fmt.Printf("gParamAssginMap=%+v\n", gParamAssignMap)
}
func timefun() {
	var timeBegin time.Time
	timeBegin = time.Now()

	for timeBegin.IsZero() {
		fmt.Println("is zero no...")
	}

}

func test(value interface{}) {
	switch value.(type) {
	case string:
		// 将interface转为string字符串类型
		op, ok := value.(string)
		fmt.Println(op, ok)
	case int32:
		// 将interface转为int32类型
		op, ok := value.(int32)
		fmt.Println(op, ok)
	case int64:
		// 将interface转为int64类型
		op, ok := value.(int64)
		fmt.Println(op, ok)
	case User:
		// 将interface转为User struct类型，并使用其Name对象
		// op, ok := value.(User)
		// fmt.Println(op.Name, ok)
	case []int:
		// 将interface转为切片类型
		op := make([]int, 0)
		op = value.([]int)
		fmt.Println("[]int arrr=", op)
	default:
		fmt.Println("unknown")
	}
}

func tstSlice() {

	any5 := []int{1, 2, 3, 4, 5, 888, 9981}
	test(any5)

}

func GetParamByName(param map[string]interface{}, name string, isInt bool) interface{} {
	if pv, ok := param[name]; ok {
		if isInt {
			v := int64(pv.(float64))
			return v
		}

		return pv
	}

	return nil
}

func GetData(param map[string]interface{}) []interface{} {
	ret := make([]interface{}, 0)

	if uid := GetParamByName(param, "uid", true); uid != nil {
		ret = append(ret, uid)
	} else {
		ret = append(ret, int64(0))
	}

	if message := GetParamByName(param, "message", false); message != nil {
		ret = append(ret, message)
	} else {
		ret = append(ret, "")
	}

	return ret
}

func GetDataEx(retSlice *[]interface{}, param map[string]interface{}) {

	if uid := GetParamByName(param, "uid", true); uid != nil {
		*retSlice = append(*retSlice, uid)
		fmt.Println("111t")
	} else {
		*retSlice = append(*retSlice, int64(0))
		fmt.Println("111")
	}

	if message := GetParamByName(param, "message", false); message != nil {
		*retSlice = append(*retSlice, message)
		fmt.Println("222t")
	} else {
		*retSlice = append(*retSlice, "")
		fmt.Println("222")
	}

}

func tstSlice1() {

	{
		var idx int = 0
		fmt.Println("idx=", idx)
		// fmt.Println("idx=", idx++)
		return
	}
	{
		result := make([]interface{}, 0)

		mapParam := make(map[string]interface{})

		var ttt float64 = 20180815
		mapParam["uid"] = ttt
		mapParam["message"] = "message2"

		GetDataEx(&result, mapParam)
		fmt.Printf("result=%+v\n", result)
		return
	}

	mapInputParam := make(map[string]interface{})

	var ttt float64 = 20180815
	// var ttt int = 20180815
	mapInputParam["uid"] = ttt
	mapInputParam["message"] = "tst map data"

	ret := GetData(mapInputParam)
	fmt.Printf("ret=%+v\n", ret)
}

var host = flag.String("host", "localhost", "host")
var port = flag.Int("port", 10077, "port")
var maxcount = flag.Int("maxcount", 30000, "max connect count")

func tstTcpClient() {

	flag.Parse()

	fmt.Printf("host: %s port: %d maxconn: %d\n", *host, *port, *maxcount)
	barrier := NewBarrier(*maxcount)

	var wg sync.WaitGroup
	wg.Add(*maxcount)

	time.Sleep(2 * time.Second)
	for i := 0; i < *maxcount; i++ {

		go func(idx int) {
			defer wg.Done()

			// fmt.Println("before work idx: ", idx)
			barrier.BarrierWait()
			// fmt.Println("idx: ", idx, " begin work")

			dstAddr := *host + ":" + int2string(*port)
			conn, err := net.Dial("tcp", dstAddr)
			if err != nil {
				fmt.Printf("idx: %d err info: %+v\n", idx, err)
				return
			}
			defer conn.Close()
			time.Sleep(time.Millisecond * 200)
			fmt.Printf("idx: %d successed! Connecting to %s\n", idx, dstAddr)
		}(i)

	}

	wg.Wait()
	fmt.Println("all tcp client done")
}

var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func AsyncRecvUpstream() {

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {

		time.Sleep(2 * time.Second)

		wg.Done()
	}()

	wg.Wait()
}

func generator(n int) <-chan int {
	outCh := make(chan int)
	go func() {
		for i := 0; i < n; i++ {
			outCh <- i
		}
		close(outCh)
	}()
	return outCh
}
func do(inCh <-chan int, outCh chan<- int, wg *sync.WaitGroup) {
	for v := range inCh {
		outCh <- v * v
	}
	wg.Done()
}
func testChannelData() {
	inCh := generator(10)
	outCh := make(chan int, 10)

	var wg sync.WaitGroup
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go do(inCh, outCh, &wg)
	}
	go func() {
		wg.Wait()
		close(outCh)
	}()
	for r := range outCh {
		fmt.Println(r)
	}
}

const (
	//MAX_READ_BYTE = 4096
	MAX_READ_BYTE = 4096
)

func readThread(idx uint64, wg *sync.WaitGroup, conn net.Conn, chResp chan []byte, ctx context.Context, cancel context.CancelFunc) {
	prefixStr := fmt.Sprintf("idx=%d ", idx)
	defer wg.Done()
	for {
		var buf [MAX_READ_BYTE]byte
		n, err := conn.Read(buf[:])

		if err != nil || n == 0 {
			cancel()
			fmt.Println(prefixStr, "recv thread exit! recv ctx done, err=", err)
			return
		} else {
			fmt.Println(prefixStr, "recv thread, recv size:", n)

			respStr := prefixStr + "server resp: "
			respBuf := make([]byte, n+len(respStr))

			copy(respBuf, []byte(respStr))
			copy(respBuf[len(respStr):], buf[0:n])

			chResp <- respBuf
		}
	}

}
func writeThread(idx uint64, wg *sync.WaitGroup, conn net.Conn, chResp chan []byte, ctx context.Context, cancel context.CancelFunc) {
	prefixStr := fmt.Sprintf("idx=%d ", idx)
	defer wg.Done()
	for {
		select {
		case buff := <-chResp:
			err := Send(conn, buff, "write thread go")
			if err != nil {
				fmt.Println(prefixStr, "write thread failed! err=", err)
				cancel()
				return
			}
			fmt.Println(prefixStr, "write thread response:", string(buff))
		case <-ctx.Done():
			fmt.Println(prefixStr, "write thread exit! recv ctx done")
			return
		}
	}

}

var (
	g_respCount uint64
)

func handleClientConnection(wg *sync.WaitGroup, conn net.Conn, ctx context.Context, cancel context.CancelFunc) {
	fmt.Println("handle client beg")
	defer fmt.Println("handle client end")
	wg.Add(2)

	chResp := make(chan []byte)
	newidx := atomic.AddUint64(&g_respCount, 1)

	go readThread(newidx, wg, conn, chResp, ctx, cancel)
	go writeThread(newidx, wg, conn, chResp, ctx, cancel)

	wg.Wait()
}

func tstListen() {
	address := ":8686"
	ln, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Println("listen client server error", address)
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Accept client server error", err, "address", address)
			return
		}

		// 接收/处理，跟client连接
		ctx, cancel := context.WithCancel(context.Background())

		var wg sync.WaitGroup

		go handleClientConnection(&wg, conn, ctx, cancel)
	}
}

func tstSelect() {
	{
		chData := make(chan int)
		go func() {
			time.Sleep(time.Millisecond * 2)
			chData <- 77
		}()
		for {
			select {
			case val := <-chData:
				fmt.Println("val=", val)
				break
			default:
				fmt.Println("default")
				time.Sleep(time.Millisecond * 500)
				break
			}
			fmt.Println(" for ----- execute ")
		}

		fmt.Println("for is done")
		return
	}
	{
		tstListen()
		return
	}
	{
		var id int = 1100
		var count int = 0
		var res int
		res = id % count
		fmt.Println("res=", res)
		return
	}

	{
		testChannelData()
		return
	}

	{
		recvBuff := make([]byte, 16)
		cpysize := copy(recvBuff, []byte("123456"))
		fmt.Println("cpysize=", cpysize, ", len = ", len(recvBuff), ", cap=", cap(recvBuff))

		recvPkg := recvBuff[0:3]
		fmt.Println("len=", len(recvPkg), ", len(recvBuff)=", len(recvBuff), ", recvPkg=", string(recvPkg))

		copy(recvPkg, []byte("xyz"))
		fmt.Println("recvBuff=", string(recvBuff))

		fmt.Println("recvPkg=", string(recvPkg))
		return
	}

	{
		recvBuff := make([]byte, 16)
		cpysize := copy(recvBuff, []byte("123"))
		fmt.Println("cpysize=", cpysize, ", len = ", len(recvBuff), ", cap=", cap(recvBuff))
		for {

			newRecvBuff := make([]byte, 128)
			recvBuff = newRecvBuff

			fmt.Println("cpysize=", cpysize, ", len = ", len(recvBuff), ", cap=", cap(recvBuff))

			time.Sleep(time.Second * 1)
		}
		return
	}

	{
		var respBytes []byte
		fmt.Printf("len=%d addr=%p\n", len(respBytes), respBytes)

		// var begPos int = 0
		// tmpBytes := make([]byte, 100)

		for i := 0; i < 10; i++ {
			tmpstr := fmt.Sprintf("%d", i)

			respBytes = append(respBytes, []byte(tmpstr)...)

			fmt.Printf("len=%d\n", len(respBytes))
			// cpysize := copy(tmpBytes[begPos:], []byte(tmpstr))
			// begPos += cpysize
		}

		fmt.Printf("len=%s\n", string(respBytes))
		return
	}
	{
		fmt.Println("tm=", time.Now(), ", beg")
		AsyncRecvUpstream()
		fmt.Println("tm=", time.Now(), ", end")
		return
	}

	{

		chData := make(chan []byte, 100)
		fmt.Println("tm=", time.Now(), ", beg")

		go func() {
			time.Sleep(time.Second)
			chData <- []byte("123")
		}()

		var heartTime int = 4
		for {
			select {
			case data := <-chData:
				{
					fmt.Println("tm=", time.Now(), ", recv data=", string(data))
				}
			case <-time.After(time.Duration(heartTime) * time.Second):
				fmt.Println("tm=", time.Now(), ", time out trig")
				break
			}
		}
		fmt.Println("tm=", time.Now(), ", end")
		select {}
		return
	}

	{
		recvBuff := make([]byte, 16)
		size := copy(recvBuff[0:], []byte("1234567"))
		fmt.Println("size=", size)
		begPos := size

		size = copy(recvBuff[begPos:], []byte("89"))
		endPos := begPos + size
		fmt.Println("endPos=", endPos)

		cpyPkg := recvBuff[:endPos]
		fmt.Println("cpyPkg=", string(cpyPkg))
	}
	return

	tmpErr := errors.New("no found upstream obj")
	fmt.Println(tmpErr)

	if tmpErr != nil {
		fmt.Println("tmperr != nil")
	}
	return

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

func tstMultiThread_Client() {
	utils.CreateFullFolder("savedir")

	upHostPort := "192.168.6.129:4008"

	for i := 0; i < 5; i++ {

		go func(i int) {
			// thrdId := fmt.Sprintf("%d", i)

			upstreamTcp, err := net.Dial("tcp", upHostPort)
			if err != nil {
				fmt.Printf("dial dd host: %s failed! err=%+v\n", upHostPort, err)
				return
			}

			sendData := fmt.Sprintf("%d send data %d%d%d%d%d", i, i, i, i, i, i)

			Send(upstreamTcp, []byte(sendData), "abc")

			recvBuff := make([]byte, 100)
			Recv(upstreamTcp, recvBuff, "abc")

			fmt.Printf("thrdId: %d recv resp:[%s]\n", i, string(recvBuff))

			time.Sleep(time.Second)
		}(i)

	}

	select {}

}

func tstwriteFile() {
	utils.AppendFile("123", []byte("123"))
	utils.AppendFile("123", []byte("456"))
	utils.AppendFile("123", []byte("67890"))
}

// 测试fast http client发送mulitfile请求，给httpserver
// httpserver返回chunked回包，要求client能正确解析；
func tstMultiPartFileReq() {

}

func createFixedBody(bodySize int) []byte {
	var b []byte
	for i := 0; i < bodySize; i++ {
		b = append(b, byte(i%10)+'0')
	}
	return b
}

func TstCh() {
	tstHttpClient1()
	return

	tstGcFun()
	return

	tstHttpClient()
	return

	bodySize := 12
	body := createFixedBody(bodySize)
	fmt.Println(body)
	return
	// tstCtxFun()
	// return

	tstSelect()
	return

	tstwriteFile()
	return

	tstMultiThread_Client()
	return

	for i := 0; i < 10; i++ {
		fmt.Println("b=", RandStringRunes(60))
	}

	tstTcpClient()
	return

	tstSlice1()
	return

	// tstCxxGo()
	return

	UseZapLog()
	return

	timefun()
	return
	ExtractPkg()
	return

	BarrFun()
	return

	{
		fmt.Println(time.Duration(8) * time.Second)
		return
	}
	{
		buf := make([]byte, 10)
		copy(buf, []byte("123"))
		fmt.Println("buf=", buf, ", len=", len(buf), ", cap=", cap(buf))

		buf = buf[:0]
		fmt.Println("buf=", buf, ", len=", len(buf), ", cap=", cap(buf))
	}
	{
		tstparseJson()
		return
		tstsomeCtx()
		return
	}
	{
		var i int
		done := make(chan struct{})
		go func() {
			i = 5
			done <- struct{}{}
		}()
		tmpval := <-done
		fmt.Println("tmpval=", tmpval, ", i=", i)
		return
	}
	{
		tmpBytes := []byte("atrings all")
		n := tmpBytes[0:1]
		fmt.Println("n=", n[0])
		return

		var b bytes.Buffer
		b.Grow(64)
		b.Write([]byte("abcde"))
		fmt.Printf("%d", b.Len())
		return

		// fmt.Println(math.Round(4.2))   //四舍五入
		// fmt.Println(math.Round(41.22)) //四舍五入
		// fmt.Println(math.Round(41.51)) //四舍五入

		// for {
		// 	select
		// }
	}

	{
		var succ float32 = 4075272.00
		var fail float32 = 8714.00

		percent := int(float32(succ) * 100.0000 / float32(succ+fail))

		// tmpval := (float64)(succ / (succ + fail))
		fmt.Println("percent===", percent)

		// var partData int = 73
		// var totalProgress int
		// (double)FilesProcessed / TotalFilesToProcess * 100);
		return
	}
	{
		for i := 10; i <= 100; i += 10 {
			str := "[" + bar(i/10, 10) + "] " + strconv.Itoa(i) + "%"
			fmt.Printf("\r%s", str)
			time.Sleep(1 * time.Second)
		}
		fmt.Println()
		return
	}
	{
		// var  correct float = 25
		// var  questionNum float = 100
		// var percent float
		// percent = (float)(correct * 100.0f / questionNum)
		// fmt.Println(percent)
		return
	}
	{
		timeBegin := time.Now()
		fmt.Println(time.Since(timeBegin).Microseconds())
		return
	}

	{
		chInts := make(chan int, 1)
		aaa, ok := <-chInts
		fmt.Println("aaa=", aaa, ", ok=", ok)
		return
		go func() {
			for i := 1; i <= 22; i++ {
				chInts <- i
			}
			close(chInts)
			time.Sleep(3 * time.Second)
			fmt.Println("close ch")
		}()

		fmt.Println("before work")
		for task := range chInts {
			fmt.Println("task=", task)
		}
		fmt.Println("after work")

		time.Sleep(1 * time.Second)
		return
	}
	{
		chInts := make(chan int)

		go func() {
			for i := 1; i <= 22; i++ {
				chInts <- i
			}
			close(chInts)
			time.Sleep(3 * time.Second)
			fmt.Println("close ch")
		}()

		fmt.Println("before work")
		for task := range chInts {
			fmt.Println("task=", task)
		}
		fmt.Println("after work")

		time.Sleep(1 * time.Second)
		return
	}

	chInts := make(chan int, 51)
	for i := 0; i < 10; i++ {
		chInts <- i
	}
	close(chInts)

	for {
		tmpval := <-chInts
		fmt.Println("tmpval=", tmpval)
	}
	fmt.Println("tstch over")
}
func tstTimeCost() {
	fmt.Println("beg update group info, tm=", time.Now())
	defer func() {
		fmt.Println("end update group info, tm=", time.Now())
	}()

	time.Sleep(time.Second * 2)

	// tmnow := time.Now()
	// defer fmt.Println("end update group info, tm=", tmnow)
}
func httpGet__1() {
	resp, err := http.Get("http://10.10.11.23:9981/pkg")
	if err != nil {
		fmt.Println("get failed, err=", err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("read all failed, err=", err)
		return
	}

	fmt.Println("resp=", string(body))
}
func tst_st() {
	type AckReq struct {
		LstAckID  uint64
		LackIDs   []uint64
		KwConnKey string
	}

	var lackIDs []uint64
	rspAck := &AckReq{ // notify sender, re send pkgs in arr;
		LstAckID:  1,
		LackIDs:   lackIDs,
		KwConnKey: "111",
	}
	// retJson := utils.Jsonize(&rspAck)

	js, _ := ffjson.Marshal(&rspAck)
	fmt.Println("js=", string(js))
}

func tst_while() {

	chExit := make(chan struct{})
	for {
		select {
		case <-chExit:
			fmt.Println("1111111111111111111111111")
			return
		case <-time.After(time.Second * time.Duration(2)):
			fmt.Println("bNotify---------------2")
			break
		case <-time.After(time.Second * time.Duration(3)):
			fmt.Println("bNotify---------------3")
			break
			// default:
			// 	fmt.Println("default ,,,,,,,,,,,,,,,,,,,,,,,")
		} // --- end select ---

	} // --- end for ---
}

var gMtxSendIdx sync.Mutex
var gSendPkgIdx int32

func tst_a() {

	for {

		count := 1
		mode := 1
		if 1 == mode {

			pfnClose := func() {
				fmt.Println("inner ,,, count=", count)
			}

			defer pfnClose()

			count++
		}

		fmt.Println("outer ,,, count=", count)
		time.Sleep(2 * time.Second)
	}
	return

	go func() {
		for {
			gMtxSendIdx.Lock()
			retIdx := gSendPkgIdx
			gSendPkgIdx = (gSendPkgIdx + 1) % (int32)(5)
			gMtxSendIdx.Unlock()

			fmt.Println(retIdx)
			time.Sleep(time.Millisecond * 300)
		}
	}()

	select {}
}

var gPoolUser *sync.Pool

func tstSyncPool() {

	gPoolUser = &sync.Pool{
		New: func() interface{} {
			obj := &UserEx{}
			return obj
		},
	}

	go func() {
		defer fmt.Println("exit,,,,go thrd")
		obj := gPoolUser.Get()
		user := obj.(*UserEx)
		user.Name = "name2"
		user.Email = "511@qq.com"

		fmt.Printf("addr=%p\n", user)

		gPoolUser.Put(user)
		// gPoolUser.Put(user)
		// gPoolUser.Put(user)
	}()

	time.Sleep(2 * time.Second)

	{
		tmp1 := gPoolUser.Get()
		obj1 := tmp1.(*UserEx)
		fmt.Printf("obj1=%p, data=%+v\n", obj1, *obj1)
	}
	{
		tmp2 := gPoolUser.Get()
		obj2 := tmp2.(*UserEx)
		fmt.Printf("obj2=%p, data=%+v\n", obj2, *obj2)
	}
	{
		tmp3 := gPoolUser.Get()
		obj3 := tmp3.(*UserEx)
		fmt.Printf("obj3=%p, data=%+v\n", obj3, *obj3)
	}
}

func main() {
	tstSyncPool()
	return
	tst_a()
	return
	tst_while()
	return
	httpGet__1()
	return

	TstCh()
	return

	// _ := syncx.NewImmutableResource(func() (nil, nil) {}, nil)
	// return
	tstGcFun()
	return

	tstTimeCost()
	return
	// fmt.Println("abcd")

	Tst1()
	return

	var obj Att

	go func(obj *Att) {
		defer fmt.Println("thread exit...")
		var info string
		info = "hello pointer..."

		obj.Message = &info

		fmt.Printf("message addr=%p\n", &info)
	}(&obj)

	time.Sleep(time.Second * 2)

	fmt.Printf("obj=%+v\n", obj)

	time.Sleep(time.Second * 3)
}

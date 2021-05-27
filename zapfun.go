package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	// "vlog"

	"github.com/pquerna/ffjson/ffjson"
	"gopkg.in/natefinch/lumberjack.v2"

	// "app/doodagent"

	// "github.com/apache/thrift/lib/go/thrift"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	// svchub "servicehub"
	// "transport"
)

func dummy_zapfunc() {
	fmt.Println("")
}

type ConfigParams struct {
	// ServerType     ServerType     `json:"server_type"`
	// ServiceHubType ServiceHubType `json:"service_hub_type"`
	Port          int32  `json:"port"`
	NumberOfConns int32  `json:"number_of_conns"`
	Name          string `json:"name"`
	IP            string `json:"ip"`
	Version       string `json:"version"`
	SID           string `json:"sid"`
	LogFile       string `json:"log_file"`
}

var (
	// theAgent      doodagent.DoodAgent
	gConfigParams ConfigParams
	gLogger       *zap.Logger
	errParam      = errors.New("param is error")
	errNullClient = errors.New("null client obj")
)

func NewZapLogger(hook *lumberjack.Logger, level zapcore.Level) (*zap.Logger, error) {
	if hook == nil {
		return nil, errParam
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	atomicLevel := zap.NewAtomicLevelAt(level)
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(hook), atomicLevel)
	caller := zap.AddCaller()

	return zap.New(core, caller), nil
}

func initZapLogger() {
	hook := &lumberjack.Logger{Filename: gConfigParams.LogFile, MaxBackups: 10, MaxAge: 30, Compress: false}
	var err error
	gLogger, err = NewZapLogger(hook, zapcore.InfoLevel)
	if err != nil {
		panic("Failed to init logger:" + err.Error())
	}
}
func initConfigParams() {
	// read config params
	data, err := ioutil.ReadFile("./ddagent.json")
	if err != nil {
		fmt.Println("Failed to read config file, quit.")
		panic("Failed to read config file.")
	}

	err = json.Unmarshal(data, &gConfigParams)
	if err != nil {
		fmt.Println("Failed to parse config param")
		panic("Failed to parse config param.")
	}
}

type NameInfo struct {
	Name string
	Age  int
}

func GetSerial(obj interface{}) string {
	retBuffer, _ := ffjson.Marshal(obj)
	return string(retBuffer)
}

func UseZapLog() {
	fmt.Println("beg ---------------------- usezaplog")
	defer fmt.Println("end ---------------------- usezaplog")

	initConfigParams()
	initZapLogger()

	url := "www.baidu.com"
	var id int64 = 123456
	info := &NameInfo{
		Name: "username_123",
		Age:  32,
	}
	retstr := GetSerial(info)

	fmt.Println("aaaaaaaaaaaaaaaaaaaa")
	// gLogger.Info("get single batch,", zap.Int64("retid", retID), zap.Error(err), zap.String("lastRet", GetSerial(resObj)))
	// gLogger.Error("Failed to create simple server", zap.Error(nil))
	var tmpErr error
	gLogger.Info("add group failed", zap.Error(errParam),
		zap.Int64("id", id),
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
		zap.String("retstr=", retstr),
		zap.Error(tmpErr))
	fmt.Println("bbbbbbbbbbbbbbbbbbbb")

	fmt.Println(strings.Trim("¡¡¡Hello, Gophers!!!", "!¡"))
	fmt.Println(strings.Trim("21256832650 ", " "))

	arrayIDs := make([]int64, 0)
	arrayIDs = append(arrayIDs, 123456)
	arrayIDs = append(arrayIDs, 998123)
	fmt.Println("arr=", arrayIDs)

	time.Sleep(time.Second * 1)
}

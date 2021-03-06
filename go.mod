module tstfun

go 1.16

require (
	Sample v0.0.0-00010101000000-000000000000
	github.com/VictoriaMetrics/fastcache v1.5.8
	github.com/allegro/bigcache v1.2.1
	github.com/apache/thrift v0.12.0
	github.com/gopherjs/gopherjs v0.0.0-20210530221158-9a6984c0bd90 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pquerna/ffjson v0.0.0-20190930134022-aa0246cd15f7
	github.com/samuel/go-zookeeper v0.0.0-20190923202752-2cc03de413da
	github.com/smartystreets/assertions v1.2.0 // indirect
	github.com/smartystreets/goconvey v1.6.4
	github.com/stretchr/testify v1.7.0 // indirect
	github.com/valyala/fasthttp v1.15.1
	go.uber.org/zap v1.15.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	utils v0.0.0-00010101000000-000000000000

// Sample v0.0.0-00010101000000-000000000000
)

replace Sample => ./gen-go/Sample

replace utils => ./utils

replace syncx => ./syncx

replace timex => ./timex

replace taskpool => ./taskpool

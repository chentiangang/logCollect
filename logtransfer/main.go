package main

import (
	"fmt"
	"project/logCollect/com"
	"project/logCollect/logtransfer/es"
	"project/logCollect/logtransfer/kafka"
	"strings"

	"github.com/chentiangang/oconfig"
	"github.com/chentiangang/xlog"
)

var app com.Transfer

func init() {
	var (
		err error
	)
	if err = oconfig.UnMarshalFile("./conf/config.ini", &app); err != nil {
		xlog.LogError("unmarshal configure failed, err:%v", err)
		panic(err)
	}
	if err = xlog.Init(app.OutPut.Type, app.OutPut.Level, app.OutPut.File, app.OutPut.Module); err != nil {
		xlog.LogDebug("init log failed, err: %v", err)
		panic(err)
	}
}

func main() {
	var (
		err        error
		topic      []string
		topicLimit int = 4
		n          int
	)
	topic = strings.Split(app.Consumer.TopicList, ",")
	for _, index := range topic {
		n++
		if n > topicLimit {
			err = fmt.Errorf("Exceed Topic the maximum limit.")
			xlog.LogError("%s", err)
			panic(err)
		}

		go func() {
			if err = es.Init(app.Es.Addr, index, app.Es.ThreadNum, app.Es.QueueSize); err != nil {
				xlog.LogError("init es failled, err:%v", err)
				panic(fmt.Sprintf("init es failed, err:%v\n", err))
				return
			}
			if err = kafka.Init(app.Consumer.Addr, index); err != nil {
				xlog.LogError("init kafka failled, err:%v", err)
				panic(err)
				return
			}
		}()
	}
	// 让主线程一直阻塞
	select {}
}

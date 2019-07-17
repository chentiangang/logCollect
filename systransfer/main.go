package main

import (
	"project/logCollect/systransfer/common"
	"project/logCollect/systransfer/influx"
	"project/logCollect/systransfer/kafka"
	"project/tools/packctl/oconfig"

	"github.com/chentiangang/xlog"
)

var appConf *common.AppConfig

func initLog() (err error) {
	var level int
	switch appConf.LogConf.LogLevel {
	case "debug":
		level = xlog.XLogLevelDebug
	case "trace":
		level = xlog.XLogLevelTrace
	case "info":
		level = xlog.XLogLevelInfo
	case "warn":
		level = xlog.XLogLevelWarn
	case "error":
		level = xlog.XLogLevelError
	default:
		level = xlog.XLogLevelDebug
	}

	err = xlog.Init(appConf.LogConf.LogType, level, appConf.LogConf.Filename, appConf.LogConf.Module)
	return
}

func main() {
	err := oconfig.UnMarshalFile("./conf/systransfer.ini", &appConf)
	if err != nil {
		xlog.LogError("unmarshal file failed, err:%v\n", err)
		return
	}

	xlog.LogDebug("app conf:%#v\n", appConf)
	err = initLog()
	if err != nil {
		xlog.LogDebug("init log failed, err:%v", err)
		return
	}

	err = influx.Init(appConf.InfluxConf.Addr)
	if err != nil {
		xlog.LogError("init influx db failed, err:%v", err)
		return
	}

	err = kafka.Init(appConf.KafkaConf.Addr, appConf.KafkaConf.Topic)
	if err != nil {
		xlog.LogError("init kafka failed, err:%v", err)
	}

	select {}
}

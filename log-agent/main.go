package main

import (
	"fmt"
	"project/logCollect/com"
	"project/logCollect/log-agent/etcd"
	"project/logCollect/log-agent/kafka"
	"project/logCollect/log-agent/tailf"

	"os"
	"strings"

	"github.com/chentiangang/oconfig"
	"github.com/chentiangang/xlog"
)

var (
	app      com.Agent
	etcdAddr = []string{"10.0.0.100:2379"}
	etcdKey  string
)

func init() {
	var (
		err      error
		hostname string
	)
	if err = oconfig.UnMarshalFile("./conf/log-agent.ini", &app); err != nil {
		return
	}
	if err = xlog.Init(app.OutPut.Type, app.OutPut.Level, app.OutPut.File, app.OutPut.Module); err != nil {
		panic(fmt.Sprintf("init logs failed, err:%v", err))
	}
	if hostname, err = os.Hostname(); err != nil {
		panic(err)
	}
	etcdKey = fmt.Sprintf("/logagent/%s", hostname)
}

func main() {
	var (
		err          error
		addrs        []string
		producerList []*com.Producer
		ProducerInCh <-chan []*com.Producer
	)
	addrs = strings.Split(app.Kafka.Addr, ",")
	if err = kafka.Init(addrs, app.Kafka.QueueSize, "user-web-api"); err != nil {
		panic(fmt.Sprintf("init kafka client failed, err:%v", err))
	}
	xlog.LogDebug("init kafka success")

	if err = etcd.Init(etcdAddr, etcdKey); err != nil {
		panic(fmt.Sprintf("init etcd client failed, err:%v", err))
	}
	xlog.LogDebug("init etcd success, addrs:%v", addrs)

	if producerList, err = etcd.GetCfg(etcdAddr, etcdKey); err != nil {
		panic(fmt.Sprintf("get etcd logagent conf failed, err:%v\n", err))
	}
	xlog.LogDebug("etcd conf:%#v", producerList)

	ProducerInCh = etcd.Watch()

	if err = tailf.Init(producerList, ProducerInCh); err != nil {
		panic(fmt.Sprintf("init tailf client failed, err:%v", err))
	}
	xlog.LogDebug("init tailf success")

	go tailf.Main()

	xlog.LogDebug("run finished.")
	select {}
}

package es

import (
	"context"

	"github.com/chentiangang/xlog"
	"github.com/olivere/elastic"
)

type Cli struct {
	client    *elastic.Client
	index     string
	threadNum int
	queue     chan interface{}
}

var (
	cli *Cli = &Cli{}
)

func Init(addr string, index string, threadNum, queueSize int) (err error) {
	var (
		client *elastic.Client
	)
	xlog.LogDebug("init es addr:%s index:%s thread_num:%d queue_size:%d",
		addr, index, threadNum, queueSize)

	if client, err = elastic.NewClient(elastic.SetURL(addr)); err != nil {
		xlog.LogError("init es client failed, err: %v", err)
		return
	}

	cli.client = client
	cli.index = index
	cli.threadNum = threadNum
	cli.queue = make(chan interface{}, queueSize)

	// 执行写入es的线程数，根据系统cpu进行配置
	for i := 0; i < threadNum; i++ {
		go writeEs()
	}
	return
}

func AcceptMsg(msg interface{}) {
	xlog.LogDebug("append msg to es queue, msg:%#v", msg)
	cli.queue <- msg
}

func writeEs() {
	for data := range cli.queue {
		_, err := cli.client.Index().
			Index(cli.index).
			Type(cli.index).BodyJson(data).Do(context.Background())
		if err != nil {
			xlog.LogError("do insert es failed, err:%v", err)
			continue
		}
	}
}

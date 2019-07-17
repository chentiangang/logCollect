package etcd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chentiangang/xlog"
	"go.etcd.io/etcd/clientv3"
)

func pushConf(addr []string) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   addr,
		DialTimeout: 3 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	defer client.Close()

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	key := fmt.Sprintf("/logagent/%s", hostname)
	topic := hostname[:len(hostname)-3]
	path := fmt.Sprintf("/app/logs/%s/%s.log", topic, topic)
	putData := fmt.Sprintf("[{\"Topic\":\"%s\",\"Path\":\"%s\"}]", topic, path)

	xlog.LogInfo("put data: %s", putData)
	if _, err = client.Put(context.Background(), key, putData); err != nil {
		panic(err)
		xlog.LogError("put failed, err:%v\n", err)
	}
}

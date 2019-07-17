package etcd

import (
	"context"
	"encoding/json"
	"project/logCollect/com"
	"time"

	"github.com/chentiangang/xlog"
	"go.etcd.io/etcd/clientv3"
)

type ClientCfg struct {
	client           *clientv3.Client
	addr             []string
	watchKey         string
	ProducerListInCh chan []*com.Producer
}

var CliCfg *ClientCfg

func Init(addr []string, watchKey string) (err error) {
	CliCfg = &ClientCfg{
		addr:             addr,
		watchKey:         watchKey,
		ProducerListInCh: make(chan []*com.Producer),
	}
	CliCfg.client, err = clientv3.New(clientv3.Config{
		Endpoints:   addr,
		DialTimeout: time.Second,
	})

	if err != nil {
		xlog.LogError("create etcd client failed, err:%v, address:%v", err, addr)
		return
	}

	go CliCfg.watch()
	return
}

func (c *ClientCfg) watch() {
	for {
		resultInCh := c.client.Watch(context.Background(), c.watchKey)
		for v := range resultInCh {
			if v.Err() != nil {
				xlog.LogError("watch:%s failed, err:%v", v.Err())
				continue
			}
			xlog.LogDebug("wacth result, value:%#v", v.Events)

			for _, event := range v.Events {
				xlog.LogDebug("event_type:%v key: %v val: %v\n", event.Type, string(event.Kv.Key), string(event.Kv.Value))
				var producer []*com.Producer
				if event.Type == clientv3.EventTypeDelete {
					continue
				}
				err := json.Unmarshal(event.Kv.Value, &producer)
				if err != nil {
					xlog.LogWarn("unmarshal failed, key: %s val:%s", c.watchKey, string(event.Kv.Value))
					continue
				}
				c.ProducerListInCh <- producer
			}
		}
	}
}

func Watch() <-chan []*com.Producer {
	return CliCfg.ProducerListInCh
}

func GetCfg(addr []string, key string) (producerList []*com.Producer, err error) {
	var (
		resp *clientv3.GetResponse
	)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if resp, err = CliCfg.client.Get(ctx, key); err != nil {
		xlog.LogError("get key:%s from etcd failed, err:%v", key, err)
		return
	}

	if len(resp.Kvs) == 0 {
		pushConf(addr)
		if resp, err = CliCfg.client.Get(ctx, key); err != nil {
			xlog.LogError("get key:%s from etcd failed, err:%v", key, err)
			return
		}
		xlog.LogInfo("push default config to etcd cluster success.")
	}

	for _, i := range resp.Kvs {
		xlog.LogInfo("etcd response key: %s", i.Key)
		xlog.LogInfo("etcd response value: %s", i.Value)
	}

	keyVals := resp.Kvs[0]
	xlog.LogDebug("get key:%s from etcd succ, key:%v val:%v", key, string(keyVals.Key), string(keyVals.Value))

	if err = json.Unmarshal(keyVals.Value, &producerList); err != nil {
		xlog.LogError("unmarshal failed, data: %v", string(keyVals.Value))
		return
	}

	xlog.LogDebug("get config from etcd succ, conf:%#v", &producerList)
	return
}

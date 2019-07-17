package kafka

import (
	"encoding/json"
	"project/logCollect/logtransfer/es"

	"github.com/Shopify/sarama"
	"github.com/chentiangang/xlog"
)

var (
	consumer sarama.Consumer
)

func Init(addr string, topic string) (err error) {
	// 1.new consumer
	consumer, err = sarama.NewConsumer([]string{addr}, nil)
	if err != nil {
		xlog.LogError("new kafka consumer failed, err:%v", err)
		return
	}

	// 2.get partitions info
	xlog.LogDebug("connect to kafka succ")
	partitions, err := consumer.Partitions(topic)
	if err != nil {
		xlog.LogError("get partition failed, err:%v", err)
		return
	}
	xlog.LogDebug("get partition succ, partition:%#v", partitions)

	// 3. reader partitionsi count
	for _, p := range partitions {
		// 4.创建一个分区消费者
		pc, err := consumer.ConsumePartition(topic, p, sarama.OffsetNewest)
		if err != nil {
			xlog.LogError("consumer partition failed, err:%v", err)
			continue
		}

		go func() {
			messageChan := pc.Messages()
			for m := range messageChan {
				xlog.LogInfo("recv from kafka, text:%v\n", string(m.Value))
				var msg map[string]interface{} = make(map[string]interface{}, 16)

				err = json.Unmarshal(m.Value, &msg)
				if err != nil {
					xlog.LogError("unmarshal failed, err:%v", err)
					continue
				}
				es.AcceptMsg(msg)
				xlog.LogInfo("append msg to es succ, msg:%#v", msg)
			}
		}()
	}
	return
}

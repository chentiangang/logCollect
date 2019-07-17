package kafka

import (
	"encoding/json"
	"project/logCollect/com"
	"regexp"
	"time"

	"github.com/Shopify/sarama"
	"github.com/chentiangang/xlog"
)

var (
	syncProducer sarama.SyncProducer
	msgCh        chan *Message
	ip           string
)

type Message struct {
	Data  string
	Topic string
}

const fmtRe = `.*\[(\d\d:\d\d:\d\d.\d\d\d)\]\[([A-Z]+).*\]\[([a-z0-9-]+)\]\[([A-Z]+)\]\[([/A-Za-z0-9]+)\] (.*)`
const errRe = `.*\[(\d\d:\d\d:\d\d.\d\d\d)\]\[([A-Z]+)[ ]*\]([\s\S]+)`

func Init(addr []string, queueSize int, topic string) (err error) {
	cfg := sarama.NewConfig()
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	cfg.Producer.Partitioner = sarama.NewRandomPartitioner
	cfg.Producer.Return.Successes = true

	if syncProducer, err = sarama.NewSyncProducer(addr, cfg); err != nil {
		xlog.LogError("product close, err:%v", err)
		xlog.LogError("address: %s, queuesize: %d", addr[0], queueSize)
		return
	}
	if ip, err = com.GetLocalIP(); err != nil {
		panic(err)
	}
	msgCh = make(chan *Message, queueSize)
	go sendToKafka()
	return
}

func sendToKafka() {
	for msg := range msgCh {
		kafkaMsg := &sarama.ProducerMessage{}
		kafkaMsg.Topic = msg.Topic
		kafkaMsg.Value = sarama.StringEncoder(msg.Data)
		partition, offset, err := syncProducer.SendMessage(kafkaMsg)
		if err != nil {
			xlog.LogError("send message to kafka failed, err:%v", err)
			return
		}
		xlog.LogDebug("send to kafka success, partition: %v offset: %v", partition, offset)
	}

}

func FmtLog(re *regexp.Regexp, s string) string {
	var (
		javaFmt  *com.JavaFmt
		response []byte
		err      error
	)
	result := re.FindAllStringSubmatch(s, -1)
	if len(result) == 0 {
		xlog.LogDebug("lines: %q", s)
		r, err := regexp.Compile(errRe)
		if err != nil {
			panic(err)
		}
		res := r.FindAllStringSubmatch(s, -1)
		sub := res[0]
		if len(sub) < 4 {
			return ""
		}
		javaFmt = &com.JavaFmt{Time: sub[1], Level: sub[2], Data: sub[3], IP: ip}
		if response, err = json.Marshal(javaFmt); err != nil {
			xlog.LogError("marshal line failed, err: %s", err)
			return ""
		}
		return string(response)
	}
	submatch := result[0]
	javaFmt = &com.JavaFmt{Time: submatch[1], Level: submatch[2], Uid: submatch[3], Method: submatch[4], ApiName: submatch[5], Data: submatch[6], IP: ip}
	if response, err = json.Marshal(javaFmt); err != nil {
		xlog.LogError("marshal line failed, err: %s", err)
		return ""
	}
	return string(response)
}

func SendToChannel(line string, topic string) (ok bool) {
	fmtRe, err := regexp.Compile(fmtRe)
	if err != nil {
		panic(err)
	}
	line = FmtLog(fmtRe, line)

	var msg *Message
	msg = &Message{
		Data:  line,
		Topic: topic,
	}
	select {
	case msgCh <- msg:
		xlog.LogDebug("send to channle success")
		return true
	default:
		xlog.LogWarn("channle is full,waiting 500ms")
		time.Sleep(time.Millisecond * 500)
		return false
	}
}

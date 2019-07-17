package tailf

import (
	"context"
	"project/logCollect/log-agent/kafka"
	"regexp"

	"github.com/chentiangang/xlog"
	"github.com/hpcloud/tail"
)

type TailObj struct {
	Path   string
	Topic  string
	tail   *tail.Tail
	ctx    context.Context
	cancel context.CancelFunc
}

func NewTailObj(path, topic string) (tailObj *TailObj, err error) {
	tailObj = &TailObj{
		Path:  path,
		Topic: topic,
	}
	// WithCancel返回带有新完成通道的父进程副本。
	// 当调用返回的cancel函数或父上下文的Done通道被关闭时，返回上下文的Done通道将被关闭，无论哪个通道先发生。
	// 取消此上下文 将释放与之关联的资源，因此代码应该在此上下文中运行的操作完成后立即调用cancel。
	tailObj.ctx, tailObj.cancel = context.WithCancel(context.Background())

	tailObj.tail, err = tail.TailFile(path, tail.Config{
		ReOpen:    true,
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll:      true,
	})
	if err != nil {
		xlog.LogError("tail file err: %v", err)
		return
	}
	return
}

func (t *TailObj) Run() {
	var lines string
	var n int
	head, err := regexp.Compile(colorRe)
	matchErr, err := regexp.Compile("ERROR")
	if err != nil {
		panic(err)
	}
	for {
		select {
		case <-t.ctx.Done():
			xlog.LogWarn("task path:%s topic: %s is exit", t.Path, t.Topic)
			return
		case line, ok := <-t.tail.Lines:
			// 如果管道里面有数据，在这个位置也可以对收集对的每一行进行格式处理
			if !ok {
				xlog.LogWarn("get message from tailf failed")
				continue
			}
			if line.Text == "" {
				continue
			}

			if head.MatchString(line.Text) || matchErr.MatchString(line.Text) {
				if n == 0 {
					lines = line.Text
					n++
					continue
				}
				n++
				if n == 2 {
					if ok := kafka.SendToChannel(lines, t.Topic); !ok {
						lines, n = line.Text, 1
						continue
					}
					lines, n = line.Text, 1
					continue
				}
			}
			lines = lines + "\n" + line.Text
		}
	}
}

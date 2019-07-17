package tailf

import (
	"fmt"
	"project/logCollect/com"

	"github.com/chentiangang/xlog"
)

type Task struct {
	TaskInMap        map[string]*TailObj
	ProducerList     []*com.Producer
	ProducerListInCh <-chan []*com.Producer
}

var (
	task *Task
	ip   string
)

func Init(producerList []*com.Producer, producerListInCh <-chan []*com.Producer) (err error) {
	if ip, err = com.GetLocalIP(); err != nil {
		panic(fmt.Sprintf("get local ip failed, err:%v\n", err))
	}
	task = &Task{
		ProducerList:     producerList,
		ProducerListInCh: producerListInCh,
		TaskInMap:        make(map[string]*TailObj, 16),
	}

	for _, cfg := range producerList {
		// 检测里是否有相同的配置
		if task.exists(cfg) {
			xlog.LogWarn("init tail task failed, configure: %#v is exists", cfg)
			continue
		}
		tailObj, err := NewTailObj(cfg.Path, cfg.Topic)
		if err != nil {
			xlog.LogWarn("init tail Task failed, conf: %#v, err:%v", cfg, err)
			continue
		}
		go tailObj.Run()
		key := fmt.Sprintf("%s+%s", tailObj.Path, tailObj.Topic)
		task.TaskInMap[key] = tailObj
	}
	return
}

func (t *Task) exists(producer *com.Producer) bool {
	for _, task := range t.TaskInMap {
		if task.Path == producer.Path && task.Topic == producer.Topic {
			return true
		}
	}
	return false
}

func (t *Task) listTask() {
	for key, topic := range t.TaskInMap {
		xlog.LogInfo("---Key: %s  Topic: %s---is running", key, topic.Topic)
	}
}

func (t *Task) run() {
	for {
		t.listTask()
		tmpProducerList := <-t.ProducerListInCh
		xlog.LogInfo("发现新的配置, data: %#v", tmpProducerList)
		for _, producer := range tmpProducerList {
			if t.exists(producer) {
				xlog.LogDebug("the task of conf: %#v is already run", producer)
				continue
			}
			xlog.LogDebug("new task of conf:%#v is running", producer)
			// 如果不存在，说明是一个新的日志收集配置.那么起现代战争新的日志收集任务
			tailObj, err := NewTailObj(producer.Path, producer.Topic)
			if err != nil {
				xlog.LogError("init tail task failed, conf: %#v, err:%v", producer, err)
				continue
			}
			go tailObj.Run()
			key := fmt.Sprintf("%s+%s", tailObj.Path, tailObj.Topic)
			task.TaskInMap[key] = tailObj
		}
		// 从已经运行的任务里面判断是否存在最新的配置当中，如果不存在的话
		// 说明这个任务的配置已经被删除
		// 检查etcd里面不存在的任务,并停止掉
		for key, task := range t.TaskInMap {
			flag := false
			for _, producer := range tmpProducerList {
				if t.exists(producer) {
					flag = true
					break
				}
			}
			if !flag {
				// 把这个任务取消,并且从tailTaskList删除
				task.cancel()
				delete(t.TaskInMap, key)
			}
		}
	}
}

func Main() {
	// task_manager主要做一些任务管理,日志收集的的工作
	task.run()
}

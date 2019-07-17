package com

type Consumer struct {
	Addr      string `ini:"address"`
	TopicList string `ini:"topic_list"`
}

type Es struct {
	Addr      string `ini:"address"`
	ThreadNum int    `ini:"thread_num"`
	QueueSize int    `ini:"queue_size"`
}

type Transfer struct {
	Consumer Consumer `ini:"consumer"`
	OutPut   OutPut   `ini:"output"`
	Es       Es       `ini:"es"`
}

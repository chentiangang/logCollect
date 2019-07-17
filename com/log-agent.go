package com

type Producer struct {
	Path  string `ini:"path"`
	Topic string `ini:"topic"`
}

type Kafka struct {
	QueueSize int    `ini:"queue_size"`
	Addr      string `ini:"address"`
}

type Agent struct {
	OutPut OutPut `ini:"output"`
	Kafka  Kafka  `ini:"kafka"`
}

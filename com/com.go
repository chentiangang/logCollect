package com

type JavaFmt struct {
	Date    string `json:"date"`
	Time    string `json:"time"`
	Level   string `json:"level""`
	Uid     string `json:"uid"`
	Method  string `json:"method"`
	ApiName string `json:"apiName"`
	Data    string `json:"data"`
	IP      string `json:"ip"`
}

type OutPut struct {
	Level  string `ini:"level"`
	File   string `ini:"file"`
	Type   string `ini:"type"`
	Module string `ini:"module"`
}

/*
type Key struct {
	Public   Public   `json:"public"`
	LogAgent LogAgent `json:"logagent"`
}

type LogAgent struct {
	HostName *[]Producer `json:"hostname"`
}

type Public struct {
	Logs struct {
		Level  string `json:"level"`
		File   string `json:"file"`
		Type   string `json:"type"`
		Module string `json:"module"`
	}

	Kafka struct {
		QueueSize int    `json:"queue_size"`
		Address   string `json"address"`
	}
}
*/

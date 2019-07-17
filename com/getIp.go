package com

import (
	"net"

	"github.com/chentiangang/xlog"
)

func GetLocalIP() (ip string, err error) {
	var (
		addrs []net.Addr
	)
	if addrs, err = net.InterfaceAddrs(); err != nil {
		return
	}
	for _, addr := range addrs {
		ipAddr, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		if ipAddr.IP.IsLoopback() || !ipAddr.IP.IsGlobalUnicast() {
			continue
		}
		xlog.LogDebug("get local ip:%#v\n", ipAddr.IP.String())
		ip = ipAddr.IP.String()
		return
	}
	return
}

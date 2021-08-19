package tcping

import (
	"fmt"
	"net"
	"strconv"
	"time"
)

func Tcping(timeout int, ports, host string) bool {
	port, err := strconv.Atoi(ports)
	if err != nil || port < 1 || port > 65535 {
		return false
	}
	_, err = net.LookupIP(host)
	if err != nil {
		return false
	}

	addr := fmt.Sprintf("%s:%d", host, port)

	_, err = net.DialTimeout("tcp", addr, time.Second*time.Duration(timeout))
	if err != nil {
		return false
	}
	return true
}

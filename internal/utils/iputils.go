package utils

import (
	"fmt"
	"net"
)

func GetLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		ip := ipNet.IP
		if ok && !ip.IsLoopback() && ip.To4() != nil {
			return ip.String(), nil
		}
	}

	return "", fmt.Errorf("cannot find non-loopback IP address")
}

func IsValidIPv4(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && parsedIP.To4() != nil
}

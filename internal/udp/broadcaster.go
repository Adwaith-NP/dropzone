package udp

import (
	"fmt"
	"net"
	"strings"
	"time"
)

func getLocalIP() (string, error) {
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

func StartBroadcast(name string, port int) {
	addr := fmt.Sprintf("255.255.255.255:%d", port)
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Println("Error dialing UDP : ", err)
		return
	}
	defer conn.Close()
	localIP, err := getLocalIP()
	if err != nil {
		fmt.Println("Error to get localIP : ", err)
		return
	}
	fmt.Println("Broadcast on UDP , port : ", port)
	for {
		message := fmt.Sprintf("%s|%s", strings.TrimSpace(name), localIP)
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error in sending UDP broadcast", err)
		} else {
			fmt.Println("Broadcasted : ", message)
		}
		time.Sleep(2 * time.Second)
	}
}

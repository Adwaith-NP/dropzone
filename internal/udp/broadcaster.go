package udp

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/Adwaith-NP/dropzone/internal/utils"
)

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
	localIP, err := utils.GetLocalIP()
	if err != nil {
		fmt.Println("Error to get localIP : ", err)
		return
	}
	fmt.Println("Broadcast on port ", port)
	for {
		message := fmt.Sprintf("%s|%s", strings.TrimSpace(name), localIP)
		_, err := conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error in sending UDP broadcast", err)
		}
		time.Sleep(2 * time.Second)
	}
}

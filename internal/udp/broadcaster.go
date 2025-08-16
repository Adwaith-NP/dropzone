package udp

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

// Send a UDP message ("name : ip") every 2 seconds
func StartBroadcast(name string, port int, localIP string, startStopSignel chan bool) {
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
	running := true // used to stop or start broadcast
	for {
		select {
		case cmd := <-startStopSignel: //Signal sent from tcp.Receive; sends false when download starts to stop broadcasting
			running = cmd
		default:
			if running {
				message := fmt.Sprintf("%s|%s", strings.TrimSpace(name), localIP)
				_, err := conn.Write([]byte(message))
				if err != nil {
					fmt.Fprintln(os.Stderr, "Error in sending broadcast", err)
				}
				time.Sleep(2 * time.Second) // set a sleeping time 2 second to terminate too much udp signels
			}
		}
	}
}

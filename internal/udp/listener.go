package udp

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type receivers struct {
	receiverMap map[string][]string
}

var store receivers

func StartListening(listenPort int) (string, error) {
	count := 0
	store = receivers{
		receiverMap: make(map[string][]string),
	}
	stopChan := make(chan struct{})
	addr := net.UDPAddr{
		Port: listenPort,
		IP:   net.IPv4zero,
	}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		return "", fmt.Errorf("failed to listen on UDP port %d: %w", listenPort, err)
	}
	defer conn.Close()

	fmt.Printf("\nListening for receivers on UDP port %d \nType x and hit enter to stop listening\n\n", listenPort)
	buf := make([]byte, 1024)
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			text := scanner.Text()
			text = strings.ToLower(text)
			if text == "x" {
				stopChan <- struct{}{}
			}
		}
	}()
	for {
		// A thread that scan user input
		conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			netErr, ok := err.(net.Error)
			if ok && netErr.Timeout() {
				select {
				case <-stopChan:
					fmt.Println("Exiting main loop.")
					break
				default:
					continue
				}

			} else {
				fmt.Printf("Error reading from UDP : %s", err)
				continue
			}
		}

		message := string(buf[:n])
		parts := strings.Split(message, "|")
		if len(parts) == 2 {
			name := parts[0]
			ip := parts[1]
			if _, exist := store.receiverMap[ip]; !exist {
				count++
				fmt.Printf("%d -> Received from %s: %s\n", count, remoteAddr.String(), message)
				store.receiverMap[ip] = []string{name, strconv.Itoa(count)}
			}
		}
	}

	// var input string
	// fmt.Print("\nSelect a receiver : ")
	// fmt.Scanln(&input)

	return "", nil
}

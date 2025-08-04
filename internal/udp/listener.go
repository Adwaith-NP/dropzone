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
	receiverMap map[string]string
}

var store receivers

// Display all broadcast and ask user to select a receiver
func StartListening(listenPort int) (string, error) {
	count := 0
	store = receivers{
		receiverMap: make(map[string]string),
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

	fmt.Printf("\nListening for receivers on UDP port %d \nHit enter to stop listening\n\n", listenPort)
	buf := make([]byte, 1024)
	//Look for user input that use to stop scanning
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			stopChan <- struct{}{}
			break
		}
	}()

loop:
	for {
		conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			netErr, ok := err.(net.Error)
			if ok && netErr.Timeout() {
				select {
				case <-stopChan:
					break loop
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
			ip := parts[1]
			if _, exist := store.receiverMap[ip]; !exist {
				count++
				fmt.Printf("%d -> Received from %s: %s\n", count, remoteAddr.String(), message)
				store.receiverMap[ip] = strconv.Itoa(count)
			}
		}
	}
	if len(store.receiverMap) != 0 {
		fmt.Println("Type x for exit")
		fmt.Print("\nChoose receiver number : ")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			choice := scanner.Text()
			choice = strings.ToLower(choice)
			if choice == "x" {
				return "", nil
			}
			num, err := strconv.Atoi(choice)
			if err == nil && num > 0 && num <= len(store.receiverMap) {
				for ip, index := range store.receiverMap {
					if choice == index {
						return ip, nil
					}
				}
			} else {
				fmt.Println("invalid receiver selection")
				fmt.Print("\nChoose receiver number : ")
			}
		}
	}

	return "", nil
}

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
	receiverMap map[string]int
}

var store receivers

// Display all broadcast and ask user to select a receiver
func StartListening(listenPort int) (string, error) {
	count := 0
	store = receivers{
		receiverMap: make(map[string]int),
	}
	stopChan := make(chan struct{}) // chanel to communicate with user input tread function that stop the UDP scanning loop
	addr := net.UDPAddr{
		Port: listenPort,
		IP:   net.IPv4zero,
	}
	conn, err := net.ListenUDP("udp", &addr) // Yes udp don't have a connection but we can use the conn variable to close the process and auther things
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
			stopChan <- struct{}{} //Pass a signel to the scanning loop if user hit enter
			break
		}
	}()
	// Its the loop that look up for UDP signels in every 2 second
loop:
	for {
		conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		n, remoteAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			netErr, ok := err.(net.Error)
			if ok && netErr.Timeout() {
				select {
				case <-stopChan: // The signel from the user input will stop the loop
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

		// If the parts len is 2 then the ip has collected and use it as a key to store in receiverMap by the value of count map[string]int {"192.32.54.12":1}
		// count represent the index of the use , so the user can select by given index
		// Key IP help to store unique user , meas the UDP signel broadcasted from every 2 second from every receiver , so we use IP to identify and display
		if len(parts) == 2 {
			ip := parts[1]
			if _, exist := store.receiverMap[ip]; !exist {
				count++                                                                         // Used to set the index for the IP
				fmt.Printf("%d -> Received from %s: %s\n", count, remoteAddr.String(), message) // Display index, connection address, username and ip if ip not precent in receiverMap
				store.receiverMap[ip] = count
			}
		}
	}
	// After ending the loop it ask to select a receiver by given index
	if len(store.receiverMap) != 0 {
		fmt.Println("Type x for exit")
		var choice string
		for {
			fmt.Print("\nChoose receiver number : ")
			fmt.Scan(&choice)
			choice = strings.ToLower(choice)
			if choice == "x" { // exit when user enter x,X
				os.Exit(0)
			}

			num, err := strconv.Atoi(choice) // convert string to int , to doing this we can use the value to comapre with index
			if err == nil && num > 0 && num <= len(store.receiverMap) {
				for ip, index := range store.receiverMap {
					if num == index {
						return ip, nil
					}
				}
			} else {
				fmt.Println("invalid receiver selection")
			}
		}
	}

	return "", nil
}

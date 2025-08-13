package tcp

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/Adwaith-NP/dropzone/cmd"
)

// get the download directory according to the os
func getDownloadDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	switch runtime.GOOS {
	case "windows":
		return filepath.Join(home, "Downloads")

	case "darwin":
		return filepath.Join(home, "Downloads")

	case "linux":
		// Check XDG user dirs config
		configFile := filepath.Join(home, ".config", "user-dirs.dirs")
		if file, err := os.Open(configFile); err == nil {
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "XDG_DOWNLOAD_DIR") {
					parts := strings.SplitN(line, "=", 2)
					if len(parts) == 2 {
						dir := strings.Trim(parts[1], `"`)
						dir = strings.ReplaceAll(dir, "$HOME", home)
						return dir
					}
				}
			}
		}
		// Fallback
		return filepath.Join(home, "Downloads")
	}

	return filepath.Join(home, "Downloads")
}

// heare is the logic to receive files
func getAllFile(conn net.Conn) error {
	baseDir := getDownloadDir()
	fmt.Println("Downloading ......")
	for {
		var metaLen int64
		if err := binary.Read(conn, binary.LittleEndian, &metaLen); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		metaBytes := make([]byte, metaLen)
		if _, err := io.ReadFull(conn, metaBytes); err != nil {
			return err
		}

		var meta map[string]string
		if err := json.Unmarshal(metaBytes, &meta); err != nil {
			return err
		}

		destPath := filepath.Join(baseDir, meta["path"])
		fmt.Println(meta["path"])

		// Its create directory if directory not exixt - 0755 is the permission mask
		if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
			return err
		}

		size, _ := strconv.ParseInt(meta["size"], 10, 64)
		f, err := os.Create(destPath)
		if err != nil {
			return err
		}

		if _, err := io.CopyN(f, conn, size); err != nil {
			f.Close()
			return err
		}

		f.Close()

	}
	return nil
}

// verify the request from sender , if granded then the download start
func ReceiveFiles(post int) error {
	var meta map[string]any
	var choice string
	port := fmt.Sprintf(":%d", post)

	ln, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	conn, err := ln.Accept()
	if err != nil {
		return err
	}
	defer conn.Close()

	//Receive the leg of json
	lengthBytes := make([]byte, 8)
	if _, err := io.ReadFull(conn, lengthBytes); err != nil {
		return err
	}

	var length int64
	for i := 0; i < 8; i++ {
		length |= int64(lengthBytes[i]) << (8 * i)
	}

	//Receive json by using the len
	jsonBytes := make([]byte, length)
	if _, err := io.ReadFull(conn, jsonBytes); err != nil {
		return err
	}

	if err := json.Unmarshal(jsonBytes, &meta); err != nil {
		return err
	}

	choice = cmd.RequestInquiry(meta)

	if choice == "y" {
		conn.Write([]byte("accepted"))
		err = getAllFile(conn)
		if err != nil {
			return err
		}
	} else {
		conn.Write([]byte("rejected"))
	}
	return nil
}

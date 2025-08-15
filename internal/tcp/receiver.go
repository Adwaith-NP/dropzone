package tcp

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func printDirTree(prefix string, tree map[string]any) {
	for key := range tree {
		value := tree[key]
		switch v := value.(type) {
		case map[string]any:
			newLine := prefix + "|--" + key
			nextPrefix := prefix + "   "
			fmt.Println(newLine)
			printDirTree(nextPrefix, v)
		default:
			fmt.Println(prefix + "|--" + key)
		}
	}
}

func requestInquiry(meta map[string]any) string {
	var choice string
	fmt.Println("File transaction request (Accept or decline):")
	if meta["Type"] == "directory" {
		name, _ := meta["Name"].(string)
		totalSize, _ := meta["TotalSize"].(float64)
		fileCount, _ := meta["FileCount"].(float64)
		treeStructure, _ := meta["TreeStructure"].(map[string]any)
		fmt.Printf("\nName: %s\nTotal Size: %.0f Bytes\nFile Count: %.0f\n", name, totalSize, fileCount)
		fmt.Print("Do you want to display file tree structure (Y - yes/enter - skip) : ")
		fmt.Scan(&choice)
		choice = strings.TrimSpace(strings.ToLower(choice))
		if choice == "y" {
			fmt.Print("\n\n")
			printDirTree("", treeStructure)
		}

	} else {
		name, _ := meta["Name"].(string)
		size, _ := meta["Size"].(float64)
		fmt.Printf("\nName: %s\nSize: %.0f Bytes\n", name, size)
	}
	fmt.Printf("\n\n* Start download (y/enter) : ")
	fmt.Scan(&choice)
	choice = strings.TrimSpace(strings.ToLower(choice))
	return choice
}

// heare is the logic to receive files
func downloadAllFile(conn net.Conn, baseDir string) error {
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
func ReceiveFiles(post int, baseDir string) error {
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

	choice = requestInquiry(meta)

	if choice == "y" {
		conn.Write([]byte("accepted"))
		err = downloadAllFile(conn, baseDir)
		if err != nil {
			return err
		}
	} else {
		conn.Write([]byte("rejected"))
	}
	return nil
}

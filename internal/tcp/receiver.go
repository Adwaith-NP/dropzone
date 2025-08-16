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

	"github.com/Adwaith-NP/dropzone/internal/utils"
)

// Print directory tree structure by using metadata
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

func requestInquiry(meta map[string]any) bool {
	var choice string
	fmt.Println("File transaction request (Accept or decline):")
	if meta["Type"] == "directory" {
		//Get directory info from meta data
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
		//Get file info from meta data
		name, _ := meta["Name"].(string)
		size, _ := meta["Size"].(float64)
		fmt.Printf("\nName: %s\nSize: %.0f Bytes\n", name, size)
	}

	fmt.Printf("\n\n* Start download (y/enter) : ")
	fmt.Scan(&choice)
	choice = strings.TrimSpace(strings.ToLower(choice))
	return choice == "y"
}

// heare is the logic to receive files
func downloadAllFile(conn net.Conn, baseDir string) error {
	fmt.Println("Downloading ......")
	for {
		var metaLen uint64
		if err := binary.Read(conn, binary.LittleEndian, &metaLen); err != nil { //get json len
			if err == io.EOF {
				break
			}
			return err
		}

		metaBytes := make([]byte, metaLen)
		if _, err := io.ReadFull(conn, metaBytes); err != nil { //by using json len its store json
			return err
		}

		var meta map[string]string
		if err := json.Unmarshal(metaBytes, &meta); err != nil { //get the meta data
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
func ReceiveFiles(post int, baseDir string, startStopSignel chan bool) error {
	var meta map[string]any // collect meta data
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
	// look up for meta data

	var length uint64
	if err := binary.Read(conn, binary.LittleEndian, &length); err != nil { //get json len
		return err
	}

	jsonBytes := make([]byte, length)
	if _, err := io.ReadFull(conn, jsonBytes); err != nil { //Receive json by using the unit64 len
		return err
	}

	if err := json.Unmarshal(jsonBytes, &meta); err != nil { //convert json data to map[string]:any
		return err
	}

	choice := requestInquiry(meta) //User interaction for avoid or accept file or directory

	if choice {
		conn.Write([]byte(utils.StatusAccepted))
		startStopSignel <- false             // pass false signel to broadcast function to stop broadcasting
		err = downloadAllFile(conn, baseDir) // start downloading the file or directory
		startStopSignel <- true              // pass true signel to broadcast function, to start broadcasting
		if err != nil {
			return err
		}
	} else {
		conn.Write([]byte(utils.StatusRejected))
	}
	return nil
}

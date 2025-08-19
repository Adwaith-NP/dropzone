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
	"strconv"
	"strings"
	"time"

	"github.com/Adwaith-NP/dropzone/internal/utils"
)

// Print directory tree structure by using metadata
func printDirTree(prefix string, tree map[string]any) {
	for key, value := range tree {
		switch v := value.(type) {
		case map[string]any:
			// Directory
			fmt.Printf("%s|-- %s\n", prefix, key)
			// Add indentation for children
			printDirTree(prefix+"   ", v)

		case float64:
			// File with size
			if v < 1024*1024 {
				kb := v / 1024.0
				fmt.Printf("%s|-- %s (%.2f KB)\n", prefix, key, kb)
			} else {
				mb := v / (1024.0 * 1024.0)
				fmt.Printf("%s|-- %s (%.2f MB)\n", prefix, key, mb)
			}
		}
	}
}

func requestInquiry(meta map[string]any) bool {
	fmt.Println("File transaction request (Accept or decline):")
	if meta["Type"] == "directory" {
		//Get directory info from meta data
		name, _ := meta["Name"].(string)
		totalSize, _ := meta["TotalSize"].(float64)
		fileCount, _ := meta["FileCount"].(float64)
		treeStructure, _ := meta["TreeStructure"].(map[string]any)

		if totalSize < 1024*1024 {
			kb := totalSize / 1024.0
			fmt.Printf("\nName: %s\nTotal Size: %.0f KB\nFile Count: %.0f\n", name, kb, fileCount)
		} else {
			mb := totalSize / (1024.0 * 1024.0)
			fmt.Printf("\nName: %s\nTotal Size: %.0f MB\nFile Count: %.0f\n", name, mb, fileCount)
		}

		reader := bufio.NewReader(os.Stdin)

		fmt.Print("Do you want to display file tree structure (Y - yes / Enter - skip) : ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(strings.ToLower(choice))

		if choice == "y" {
			fmt.Println()
			printDirTree("", treeStructure)
		}
	} else {
		//Get file info from meta data
		name, _ := meta["Name"].(string)
		size, _ := meta["Size"].(float64)
		if size < 1024*1024 {
			kb := size / 1024.0
			fmt.Printf("\nName: %s\nSize: %.2f KB", name, kb)
		} else {
			mb := size / (1024.0 * 1024.0)
			fmt.Printf("\nName: %s\nSize: %.2f MB", name, mb)
		}
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n\n* Start download (y/enter) : ")
	input, _ := reader.ReadString('\n') // read until Enter
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y"
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

		wd := &utils.DropData{
			Writer:        f,
			LastTime:      time.Now(),
			TotalFileSize: size,
		}

		if _, err := io.CopyN(wd, conn, size); err != nil {
			f.Close()
			return err
		}
		f.Close()
		fmt.Printf("\r[%s] 100%% (Done)\n", strings.Repeat("|", utils.BARWIDTH))
	}
	fmt.Print("\n\nDownload complete")
	return nil
}

// verify the request from sender , if granded then the download start
func ReceiveFiles(port int, baseDir string, startStopSignel chan bool) error {
	var meta map[string]any // collect meta data
	lnAdrr := fmt.Sprintf("0.0.0.0:%d", port)
	ln, err := net.Listen("tcp", lnAdrr)
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

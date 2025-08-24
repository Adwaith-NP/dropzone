package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const BARWIDTH = 50
const DisplayInterval = 100 * time.Millisecond
const SpeedCalInterval = time.Second

type DropData struct {
	Writer        io.Writer
	Current       int64
	LastTime      time.Time
	TotalFileSize int64
	LastBytes     int64
	Speed         string
	LastTimeSpeed time.Time
}

func (wd *DropData) Write(data []byte) (int, error) {
	n, err := wd.Writer.Write(data)
	wd.Current += int64(n)
	if time.Since(wd.LastTime) >= DisplayInterval {
		percent := float64(wd.Current) / float64(wd.TotalFileSize) * 100
		filled := int(percent * float64(BARWIDTH) / 100)

		bytesThisSecond := wd.Current - wd.LastBytes
		speed := float64(bytesThisSecond) / time.Since(wd.LastTime).Seconds()

		if time.Since(wd.LastTimeSpeed) >= SpeedCalInterval {
			if speed > 1024*1024 {
				wd.Speed = fmt.Sprintf("%.2f MB/s", speed/(1024*1024))
			} else {
				wd.Speed = fmt.Sprintf("%.2f KB/s", speed/1024)
			}
			wd.LastTimeSpeed = time.Now()
		}

		bar := "[" + strings.Repeat("\033[32m|\033[0m", filled) + strings.Repeat(" ", BARWIDTH-filled) + "]"
		fmt.Printf("\r%s %.0f%% (%s)", bar, percent, wd.Speed)
		wd.LastTime = time.Now()
		wd.LastBytes = wd.Current
	}
	return n, err
}

type DirectoryMeta struct {
	Type          string
	Name          string
	TotalSize     int64
	FileCount     int
	TreeStructure map[string]any
}

type FileMeta struct {
	Type string
	Name string
	Size int64
}

// Return true if the path point a directory
func IsDirectory(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// Return the file same as string
func GetAllFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		name := info.Name()
		if strings.HasPrefix(name, ".") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, fmt.Errorf("no files found in directory")
	}
	return files, nil
}

func BuildFileListMeta(fileList map[string]any, totalSize int64) DirectoryMeta {
	meta := DirectoryMeta{
		Type:          "directory",
		Name:          "DropZomeFiles",
		TotalSize:     totalSize,
		FileCount:     len(fileList),
		TreeStructure: fileList,
	}
	return meta
}

func BuildFileMeta(path string) (FileMeta, error) {
	info, err := os.Stat(path)
	if err != nil {
		return FileMeta{}, nil
	}
	meta := FileMeta{
		Type: "file",
		Name: info.Name(),
		Size: info.Size(),
	}
	return meta, nil
}

func BuildDirectoryMeta(rootPath string) (DirectoryMeta, error) {
	tree := make(map[string]any)
	var totalSize int64
	var fileCount int

	err := buildTreeRecursive(rootPath, tree, &totalSize, &fileCount)
	if err != nil {
		return DirectoryMeta{}, err
	}
	meta := DirectoryMeta{
		Type:          "directory",
		Name:          filepath.Base(rootPath),
		TotalSize:     totalSize,
		FileCount:     fileCount,
		TreeStructure: tree,
	}
	return meta, nil
}

func buildTreeRecursive(path string, tree map[string]any, totalSize *int64, fileCount *int) error {
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}
		fullPath := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			subTree := make(map[string]any)
			tree[entry.Name()] = subTree

			err := buildTreeRecursive(fullPath, subTree, totalSize, fileCount)
			if err != nil {
				return err
			}
		} else {
			info, err := entry.Info()
			if err != nil {
				return err
			}
			tree[info.Name()] = info.Size()
			*totalSize += info.Size()
			*fileCount++
		}
	}
	return nil
}

// get the download directory according to the os
func GetDownloadDir() string {
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

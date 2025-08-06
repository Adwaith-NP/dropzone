package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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

// Return true if the file path exist
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
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

package cmd

import (
	"fmt"
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

func RequestInquiry(meta map[string]any) string {
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

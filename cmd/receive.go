package cmd

import (
	"fmt"
	"strings"
)

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
			for key := range treeStructure {
				value := treeStructure[key]
				fmt.Println(value)
				fmt.Printf("Type of a: %T\n\n", value)
				// switch v := value.(type) {
				// case string:
				// 	fmt.Println(v, "str")
				// case int64:
				// 	fmt.Println(v, "int")
				// case map[string]any:
				// 	fmt.Println(v, "{map}")
				// default:
				// 	fmt.Println("error : ", v)
				// }

			}
		}

	} else {
		name, _ := meta["Name"].(string)
		size, _ := meta["Size"].(float64)
		fmt.Printf("\nName: %s\nSize: %.0f Bytes\n", name, size)
	}
	fmt.Printf("\n\n* Type n and hit enter for reject\n* Type y and hit enter to accept\nInput : ")
	fmt.Scan(&choice)
	choice = strings.TrimSpace(strings.ToLower(choice))
	return choice
}

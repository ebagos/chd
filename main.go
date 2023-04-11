package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	username := "username"
	defaultDir := filepath.Join("C:", "Users", username, "Downloads")

	dirPtr := flag.String("dir", defaultDir, "Directory to process")
	lengthPtr := flag.Int("length", 15, "Length of characters to compare")

	flag.Parse()

	processDirectory(*dirPtr, *lengthPtr)
}

func processDirectory(curDir string, length int) {
	count := 0
	chked := 0

	absPath, err := filepath.Abs(curDir)
	if err != nil {
		fmt.Println("Error converting path to absolute:", err)
		return
	}

	all, err := os.ReadDir(absPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	os.Chdir(absPath)

	for index, f := range all {
		end := len(all) - 1
		info1, err := f.Info()
		if err != nil {
			fmt.Println("Error getting file info:", err)
			continue
		}

		for i := index + 1; i < end; i++ {
			chked++

			if len(f.Name()) < length {
				continue
			}
			if len(all[i].Name()) < length {
				continue
			}

			str := f.Name()[0:length]
			cmp := all[i].Name()[0:length]
			if str == cmp {
				info2, err := all[i].Info()
				if err != nil {
					fmt.Println("Error getting file info:", err)
					continue
				}

				if info1.Size() <= info2.Size() {
					os.Remove(f.Name())
				} else {
					os.Remove(all[i].Name())
				}
				count++
			}
		}
	}
	fmt.Println("count =", count, "checked =", chked)
}

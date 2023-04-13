package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type Config struct {
	Dir    string `json:"dir"`
	Length int    `json:"length"`
}

func main() {
	configFile := flag.String("config", "", "Path to the configuration file")
	dirPtr := flag.String("dir", "", "Directory to process")
	lengthPtr := flag.String("length", "", "Length of characters to compare")
	flag.Parse()

	config, err := readConfig(*configFile)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}

	dir := getSetting(*dirPtr, config.Dir, "PROCESS_DIR", "C:\\Users")
	length := getSettingInt(*lengthPtr, config.Length, "PROCESS_LENGTH", 15)

	processDirectory(dir, length)
}

func readConfig(configFile string) (*Config, error) {
	if configFile == "" {
		return &Config{}, nil
	}

	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func getSetting(cmdVal, fileVal, envKey, defaultVal string) string {
	if cmdVal != "" {
		return cmdVal
	}

	if fileVal != "" {
		return fileVal
	}

	envVal := os.Getenv(envKey)
	if envVal != "" {
		return envVal
	}

	return defaultVal
}

func getSettingInt(cmdVal string, fileVal int, envKey string, defaultVal int) int {
	if cmdVal != "" {
		intVal, err := strconv.Atoi(cmdVal)
		if err == nil {
			return intVal
		}
	}

	if fileVal != 0 {
		return fileVal
	}

	envVal := os.Getenv(envKey)
	if envVal != "" {
		intVal, err := strconv.Atoi(envVal)
		if err == nil {
			return intVal
		}
	}

	return defaultVal
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
		end := len(all) // 修正箇所
		info1, err := f.Info()
		if err != nil {
			fmt.Println("Error getting file info:", err)
			continue
		}

		for i := index + 1; i < end; i++ { // 修正箇所
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

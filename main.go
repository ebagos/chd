package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

type Config struct {
	Dir      string `json:"dir"`
	Length   int    `json:"length"`
	HashAlgo string `json:"hash_algo"`
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
	hashAlgo := getSetting(config.HashAlgo, "MD5", "HASH_ALGO", "")

	processDirectory(dir, length, hashAlgo)
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

func getFileHash(filePath, hashAlgo string) (string, error) {
	var hash []byte
	var err error

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	if hashAlgo == "SHA-1" {
		h := sha1.New()
		if _, err := io.Copy(h, file); err != nil {
			return "", err
		}
		hash = h.Sum(nil)
	} else if hashAlgo == "SHA-256" {
		h := sha256.New()
		if _, err := io.Copy(h, file); err != nil {
			return "", err
		}
		hash = h.Sum(nil)
	} else { // MD5 by default
		h := md5.New()
		if _, err := io.Copy(h, file); err != nil {
			return "", err
		}
		hash = h.Sum(nil)
	}

	return hex.EncodeToString(hash), nil
}

func processDirectory(curDir string, length int, hashAlgo string) {
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

	// Map to store file hash values and their corresponding file paths
	hashMap := make(map[string]string)

	for _, f := range all {
		info, err := f.Info()
		if err != nil {
			fmt.Println("Error getting file info:", err)
			continue
		}

		if info.IsDir() {
			continue
		}

		filePath := f.Name()

		// Ignore files with name length less than specified length
		if len(filePath) < length {
			continue
		}

		hash, err := getFileHash(filePath, hashAlgo)
		if err != nil {
			fmt.Println("Error computing file hash:", err)
			continue
		}

		// If hash already exists in the map, delete the current file
		if existingFile, ok := hashMap[hash]; ok {
			if info.Size() <= getFileSize(existingFile) {
				os.Remove(filePath)
			} else {
				os.Remove(existingFile)
				hashMap[hash] = filePath
			}
			count++
		} else {
			hashMap[hash] = filePath
		}

		chked++
	}

	fmt.Println("count =", count, "checked =", chked)
}

func getFileSize(filePath string) int64 {
	info, err := os.Stat(filePath)
	if err != nil {
		fmt.Println("Error getting file info:", err)
		return 0
	}
	return info.Size()
}

package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestReadConfig(t *testing.T) {
	tempFile, err := os.CreateTemp("", "temp_config")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	tempConfig := Config{
		Dir:    "C:\\TestDir",
		Length: 10,
	}

	encoder := json.NewEncoder(tempFile)
	err = encoder.Encode(tempConfig)
	if err != nil {
		t.Fatal(err)
	}
	tempFile.Close()

	config, err := readConfig(tempFile.Name())
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if config.Dir != tempConfig.Dir {
		t.Errorf("Expected dir: %s, got: %s", tempConfig.Dir, config.Dir)
	}

	if config.Length != tempConfig.Length {
		t.Errorf("Expected length: %d, got: %d", tempConfig.Length, config.Length)
	}
}

func TestGetSetting(t *testing.T) {
	cmdVal := "cmdVal"
	fileVal := "fileVal"
	envKey := "TEST_ENV"
	defaultVal := "defaultVal"
	os.Setenv(envKey, "envVal")
	defer os.Unsetenv(envKey)

	got := getSetting(cmdVal, fileVal, envKey, defaultVal)
	if got != cmdVal {
		t.Errorf("Expected: %s, got: %s", cmdVal, got)
	}

	got = getSetting("", fileVal, envKey, defaultVal)
	if got != fileVal {
		t.Errorf("Expected: %s, got: %s", fileVal, got)
	}

	got = getSetting("", "", envKey, defaultVal)
	if got != os.Getenv(envKey) {
		t.Errorf("Expected: %s, got: %s", os.Getenv(envKey), got)
	}

	got = getSetting("", "", "", defaultVal)
	if got != defaultVal {
		t.Errorf("Expected: %s, got: %s", defaultVal, got)
	}
}

func TestGetSettingInt(t *testing.T) {
	cmdVal := "10"
	fileVal := 20
	envKey := "TEST_ENV_INT"
	defaultVal := 30
	os.Setenv(envKey, "40")
	defer os.Unsetenv(envKey)

	got := getSettingInt(cmdVal, fileVal, envKey, defaultVal)
	cmdInt, _ := strconv.Atoi(cmdVal)
	if got != cmdInt {
		t.Errorf("Expected: %d, got: %d", cmdInt, got)
	}

	got = getSettingInt("", fileVal, envKey, defaultVal)
	if got != fileVal {
		t.Errorf("Expected: %d, got: %d", fileVal, got)
	}

	got = getSettingInt("", 0, envKey, defaultVal)
	envInt, _ := strconv.Atoi(os.Getenv(envKey))
	if got != envInt {
		t.Errorf("Expected: %d, got: %d", envInt, got)
	}

	got = getSettingInt("", 0, "", defaultVal)
	if got != defaultVal {
		t.Errorf("Expected: %d, got: %d", defaultVal, got)
	}
}

func TestProcessDirectory(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "temp_dir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	fileNames := []string{
		"file_001_abc.txt",
		"file_001_def.txt",
		"file_002_abc.txt",
		"file_002_def.txt",
	}

	for _, fileName := range fileNames {
		err := os.WriteFile(filepath.Join(tempDir, fileName), []byte("test content"), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	length := 8
	processDirectory(tempDir, length)

	files, err := os.ReadDir(tempDir)
	if err != nil {
		t.Fatal(err)
	}

	expectedFiles := 2
	if len(files) != expectedFiles {
		t.Errorf("Expected %d files, got: %d", expectedFiles, len(files))
	}

	remainingFiles := make(map[string]bool)
	for _, f := range files {
		remainingFiles[f.Name()] = true
	}

	expectedRemovedFiles := []string{"file_001_abc.txt", "file_002_abc.txt"}
	for _, fileName := range expectedRemovedFiles {
		if remainingFiles[fileName] {
			t.Errorf("Expected file %s to be removed, but it was not", fileName)
		}
	}
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

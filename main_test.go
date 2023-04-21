/*
main.goのユニットテスト
*/
package main

import (
	"encoding/json"
	"io"
	"os"
	"testing"
)

func TestReadConfig(t *testing.T) {
	// テスト用の設定ファイルを作成
	configType := Config{
		Dir:      "C:\\Users",
		Length:   15,
		HashAlgo: "MD5",
	}
	configJson, _ := json.Marshal(configType)
	os.WriteFile("test.json", configJson, 0644)
	defer os.Remove("test.json")

	// テスト用の設定ファイルを読み込む
	config, err := readConfig("test.json")
	if err != nil {
		t.Errorf("Error reading config file: %s", err)
	}

	// 設定ファイルの内容を確認
	if config.Dir != "C:\\Users" {
		t.Errorf("Error reading config file: %s", err)
	}
	if config.Length != 15 {
		t.Errorf("Error reading config file: %s", err)
	}
	if config.HashAlgo != "MD5" {
		t.Errorf("Error reading config file: %s", err)
	}
}

func TestGetSetting(t *testing.T) {
	// テスト用の設定ファイルを作成
	configType := Config{
		Dir:      "C:\\Users",
		Length:   15,
		HashAlgo: "MD5",
	}
	configJson, _ := json.Marshal(configType)
	os.WriteFile("test.json", configJson, 0644)
	defer os.Remove("test.json")

	// テスト用の設定ファイルを読み込む
	config, err := readConfig("test.json")
	if err != nil {
		t.Errorf("Error reading config file: %s", err)
	}

	// テスト用の環境変数を設定
	os.Setenv("PROCESS_DIR", "C:\\Users")
	os.Setenv("PROCESS_LENGTH", "15")
	os.Setenv("HASH_ALGO", "MD5")

	// テスト用のコマンドライン引数を設定
	dirPtr := "C:\\Users"
	lengthPtr := "15"

	// 設定ファイルの内容を確認
	if getSetting(dirPtr, config.Dir, "PROCESS_DIR", "C:\\Users") != "C:\\Users" {
		t.Errorf("Error reading config file: %s", err)
	}
	if getSettingInt(lengthPtr, config.Length, "PROCESS_LENGTH", 15) != 15 {
		t.Errorf("Error reading config file: %s", err)
	}
	if getSetting(config.HashAlgo, "MD5", "HASH_ALGO", "") != "MD5" {
		t.Errorf("Error reading config file: %s", err)
	}
}

func TestGetFileHash(t *testing.T) {
	// テスト用のファイルを作成
	os.WriteFile("test.txt", []byte("test"), 0644)
	defer os.Remove("test.txt")

	// ハッシュ値を取得
	hash, err := getFileHash("test.txt", "MD5")
	if err != nil {
		t.Errorf("Error reading config file: %s", err)
	}

	// ハッシュ値を確認
	if hash != "098f6bcd4621d373cade4e832627b4f6" {
		t.Errorf("Error reading config file: %s", err)
	}

	// SHA-1のハッシュ値を取得
	hash, err = getFileHash("test.txt", "SHA-1")
	if err != nil {
		t.Errorf("Error reading config file: %s", err)
	}
	// SHA-1のハッシュ値を確認
	if hash != "a94a8fe5ccb19ba61c4c0873d391e987982fbbd3" {
		t.Errorf("Error reading config file: %s", err)
	}

	// SHA-256のハッシュ値を取得
	hash, err = getFileHash("test.txt", "SHA-256")
	if err != nil {
		t.Errorf("Error reading config file: %s", err)
	}
	// SHA-256のハッシュ値を確認
	if hash != "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08" {
		t.Errorf("Error reading config file: %s", err)
	}
}

func TestProcessDirectory(t *testing.T) {
	// テスト用のディレクトリ作成
	os.Mkdir("test", 0755)
	defer os.RemoveAll("test")
	// テスト用のファイル群作成
	file1, err := os.Create("test/test1.txt")
	if err != nil {
		t.Errorf("Error reading config file: %s", err)
	}
	defer file1.Close()
	io.WriteString(file1, "test1")
}

func TestGetFileSize(t *testing.T) {
	// テスト用のファイルを作成
	os.WriteFile("test.txt", []byte("test"), 0644)
	defer os.Remove("test.txt")

	// ファイルサイズを取得
	size := getFileSize("test.txt")

	// ファイルサイズを確認
	if size != 4 {
		t.Error("Error reading config file")
	}
}

func TestMain(m *testing.M) {
	// テスト実行
	exitCode := m.Run()

	// テスト終了
	os.Exit(exitCode)
}

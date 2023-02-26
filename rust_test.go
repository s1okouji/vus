package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func beforeEach() {
	root := "./tmp/"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			fmt.Println("Directory:", path)
		} else {
			fmt.Println("File:", path)
		}
		err = os.RemoveAll(path)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Println("File not exist")
			} else {
				return err
			}
		}
		fmt.Println("Removed:", path)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
}

func getToken() string {
	// var config RustConfig
	// config.Config = *LoadConfig()
	// return config.GetToken()
	return os.Getenv("token")
}

func check(except []string) {
	root := "./tmp" // 探索を開始するルートディレクトリ
	i := 0

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// ディレクトリの場合は何もしない
		if info.IsDir() {
			return nil
		}
		// ファイルの場合はファイル名を出力
		fmt.Println(path)
		// 期待値と確認
		if except[i] != path {
			fmt.Errorf("CheckError: Excepted %s ,but %s", except[i], path)
		}
		i += 1
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
}

func TestUnzip_normal1(t *testing.T) {
	beforeEach()
	unzip("./test_data/sample.zip", "./tmp/")
	var except []string
	except = append(except, "tmp\\sample\\sample.txt")
	check(except)
}

func TestUnzip_normal2(t *testing.T) {
	beforeEach()
	err := os.MkdirAll("./tmp/sample", os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = ioutil.WriteFile("./tmp/sample/example.txt", []byte("aaaa"), 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	unzip("./test_data/before_data/sample.zip", "./tmp/")
	unzip("./test_data/sample.zip", "./tmp/")
	var except []string
	except = append(except, "tmp\\sample\\example.txt")
	except = append(except, "tmp\\sample\\sample.txt")
	except = append(except, "tmp\\sample\\sample_b.txt")
	check(except)
}
func TestDownload(t *testing.T) {
	var config RustConfig
	config.Config = *LoadConfig()
	downloadLatestOxide(config.GetServerPath(), getToken())
}

func TestGetURL(t *testing.T) {
	fmt.Println(getLinuxAssetsUrl(getToken()))
}

func TestDeploy(t *testing.T) {
	DeployOxide()
}

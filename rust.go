package main

import (
	"archive/zip"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func DeployOxide() error {
	var config RustConfig
	config.Config = *LoadConfig()
	loadedVersion := config.LoadedVersion()
	latestVersion := getLatestVersion(config.GetToken())
	log.Println("Loaded Version:" + loadedVersion)
	log.Println("The latest Version:" + latestVersion)

	if latestVersion != loadedVersion {
		log.Println("Update to New Version" + latestVersion)
		log.Println("Stop Rust Process")
		// TODO: 潜在的にバグの可能性を秘めているので修正するべき
		// err := exec.Command("systemctl", "stop", "rust-server.service").Wait()
		// if err != nil {
		// 	return err
		// }
		log.Println("Update Rust Dedicated Version. Please Wait.")
		// script := "/usr/games/steamcmd +@sSteamCmdForcePlatformType linux +force_install_dir /opt/rust_server +login anonymous +app_update 258550 validate +quit"
		// cmd_array := strings.Split(script, " ")
		// out, err := exec.Command(cmd_array[0], cmd_array[1:]...).Output()
		// if err != nil {
		// 	return err
		// }

		// log.Println(string(out))

		err := downloadLatestOxide(config.GetServerPath())
		if err != nil {
			return err
		}
		// TODO 結構記憶があいまいなので要確認
		unzip(config.GetServerPath()+"/Oxide.zip", config.GetServerPath())
		log.Println("Finish to unzip Oxide Files")
		os.Remove(config.GetServerPath() + "/Oxide.zip")
		// out, err = exec.Command("systemctl", "start", "rust-server.service").Output()
		// if err != nil {
		// 	return err
		// }

		// log.Println(out)

	} else {
		log.Println("Need not to update")
	}

	return nil
}

func downloadLatestOxide(filepath string) error {
	resp, err := http.Get(getLinuxAssetsUrl())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath + "/Oxide.zip")
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// filepath = "~/hogehoge/sample.zip"
func unzip(filePath string, target string) error {
	archive, err := zip.OpenReader(filePath)
	if err != nil {
		return err
	}
	defer archive.Close()

	for _, file := range archive.File {
		unzippedFile, err := file.Open()
		if err != nil {
			return err
		}

		defer unzippedFile.Close()

		targetPath := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(targetPath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), os.ModePerm); err != nil {
			return err
		}

		// 出力先のファイルを作成
		// outputFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0775)
		outputFile, err := os.Create(targetPath)

		if err != nil {
			return err
		}

		defer outputFile.Close()

		// zipファイル内のファイルを出力先にコピーする
		if _, err := io.Copy(outputFile, unzippedFile); err != nil {
			return err
		}
	}

	return nil
}

func getLatestVersion(Token string) string {
	client := http.DefaultClient
	req, _ := http.NewRequest("GET", "https://api.github.com/repos/OxideMod/Oxide.Rust/releases/latest", nil)
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", "Bearer "+Token)
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	byteArray, _ := io.ReadAll(resp.Body)
	var data map[string]interface{}
	json.Unmarshal(byteArray, &data)

	return data["name"].(string)
}

func getLinuxAssetsUrl() string {
	var ret string
	client := http.DefaultClient
	req, _ := http.NewRequest("GET", "https://api.github.com/repos/OxideMod/Oxide.Rust/releases/latest", nil)
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("Authorization", "Bearer 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")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	byteArray, _ := io.ReadAll(resp.Body)
	var data map[string]interface{}
	json.Unmarshal(byteArray, &data)
	assets := data["assets"].([]interface{})
	for _, v := range assets {
		datum := v.(map[string]interface{})
		// For Linux
		if datum["name"].(string) == "Oxide.Rust-linux.zip" {
			// fmt.Println(datum["url"])
			ret = datum["browser_download_url"].(string)
		}

		// For Windows
		// if datum["name"].(string) == "Oxide.Rust.zip" {
		// 	// fmt.Println(datum["url"])
		// 	ret = datum["browser_download_url"].(string)
		// }
	}
	return ret
}

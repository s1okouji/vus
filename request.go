package main

import (
	"fmt"
)

func main() {
	// fmt.Println("Start request to Github API")
	// fmt.Println(getLatestVersion())
	// fmt.Println(getLinuxAssetsUrl())
	// fmt.Println(loadedVersion(loadConfig()))
	var config RustConfig
	config.Config = *LoadConfig()
	fmt.Println(config.LoadedVersion())
	DeployOxide()
}

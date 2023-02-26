package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Data *map[string]interface{}
}
type Loader interface {
	LoadedVersion()
}

func (c *Config) LoadedVersion() string {
	fmt.Println(c.Data)
	return ""
}

func LoadConfig() *Config {
	bytes, err := ioutil.ReadFile("packages.yml")
	if err != nil {
		panic(err)
	}

	var data map[string]interface{}
	yaml.Unmarshal(bytes, &data)
	var ret Config
	ret.Data = &data
	return &ret
}

func (c *Config) GetToken() string {
	// return (*c.Data)["github"].(map[string]interface{})["token"].(string)
	return os.Getenv("token")
}

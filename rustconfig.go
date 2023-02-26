package main

type RustConfig struct {
	Config
}

func (c *RustConfig) LoadedVersion() string {
	return (*c.Config.Data)["rust"].(map[string]interface{})["oxide"].(map[string]interface{})["version"].(string)
}

func (c *RustConfig) GetServerPath() string {
	return (*c.Config.Data)["rust"].(map[string]interface{})["path"].(string)
}

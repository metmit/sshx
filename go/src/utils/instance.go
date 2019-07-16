package utils

import "sync"

//Singleton File
var once sync.Once
var instance *File

func GetFileInstance() *File {
	once.Do(func() {
		instance = new(File)
	})
	return instance
}

//Singleton Str
var strOnce sync.Once
var strInstance *Str

func GetStrInstance() *Str {
	strOnce.Do(func() {
		strInstance = new(Str)
	})
	return strInstance
}

//Singleton config
var configOnce sync.Once
var configInstance *Config

func GetConfig() *Config {
	configOnce.Do(func() {
		configInstance = &Config{Version: "1.1.0", FolderName: ".sfsshx"}
	})
	return configInstance
}

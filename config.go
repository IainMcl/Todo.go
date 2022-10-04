package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type IConfig interface {
	ReadUserConfig() error
	View() error
	Set(key string, value string) error
	Delete(key string) error
}

type Config struct {
	UserConfig string
	ConfigFile string // file path and file name of config
	DbName     string
	TableName  string
}

type ConfigType int

const (
	defaultConfig ConfigType = iota
	userConfig
)

func readJson[T Config](path string, obj *T) error {
	// Read json file
	jsonFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(byteValue, obj)
	if err != nil {
		panic(err)
	}
	return nil
}

func (c *Config) ReadConfig(configType ConfigType) (Config, error) {
	var path string
	switch configType {
	case defaultConfig:
		path = c.ConfigFile
	case userConfig:
		path = c.UserConfig
	}

	err := readJson[Config](path, c)
	if err != nil {
		return Config{}, err
	}
	return *c, nil
}

func (c *Config) WriteConfig(configType ConfigType) error {
	return nil
}

func (c *Config) View() error {
	// Read config file
	config, err := c.ReadConfig(defaultConfig)
	if err != nil {
		return err
	}

	// Read user config file
	userConfig, err := c.ReadConfig(userConfig)
	if err != nil {
		return err
	}

	// update config with values from user config
	config.DbName = userConfig.DbName
	config.TableName = userConfig.TableName

	// print config
	config.Print()
	return nil
}

func (c *Config) Print() {
	fmt.Println("DbName: ", c.DbName)
	fmt.Println("TableName: ", c.TableName)
}

func (c *Config) Set(key string, value string) error {
	return nil
}

func (c *Config) Delete(key string) error {
	return nil
}

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type IConfig interface {
	ReadConfig() error
	View() error
	Set(key string, value string) error
	Delete(key string) error
	GetConfigPath() string
	GetDbName() string
	GetConfig() map[string]string
	GetTableName() string
	CreateDefaultConfig() error
}

type Config struct {
	ConfigPath string // file path and file name of config
	DbName     string
	TableName  string
	config     map[string]string
}

func (c *Config) GetConfig() map[string]string {
	return c.config
}

func (c *Config) GetDbName() string {
	return c.DbName
}

func (c *Config) GetTableName() string {
	return c.TableName
}

func (c *Config) GetConfigPath() string {
	return c.ConfigPath
}

func readJson[T interface{ Config }](path string, obj *T) error {
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

func (c *Config) ReadConfig() error {
	err := readJson[Config](c.ConfigPath, c)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) WriteConfig() error {
	return nil
}

func (c *Config) View() error {
	// Read config file
	err := c.ReadConfig()
	if err != nil {
		return err
	}
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

func (c *Config) CreateDefaultConfig() error {
	// Create a config file at config.ConfigPath
	// If the file already exists, return an error

	// Create config folder
	configFolder := c.ConfigPath[:len(c.ConfigPath)-len("config.json")] // TODO: This is a hack. Find a better way to get the folder path
	err := os.MkdirAll(configFolder, 0755)
	if err != nil {
		return err
	}
	// Create the config file
	file, err := os.Create(c.ConfigPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the default config to the file
	_, err = file.WriteString(`{
		"DbName": "~/.todo/todo.db",
		"TableName": "todo"
	}`)
	if err != nil {
		return err
	}

	return nil
}

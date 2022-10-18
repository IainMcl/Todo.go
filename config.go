package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
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
	WriteConfig() error
}

type Config struct {
	ConfigPath string // file path and file name of config
	DbName     string
	TableName  string
	config     map[string]interface{}
}

func NewConfig() *Config {
	home, _ := os.UserHomeDir()
	dotPath := filepath.Join(home, ".todo")

	return &Config{
		ConfigPath: filepath.Join(dotPath, "config.json"),
		DbName:     filepath.Join(dotPath, "todo.db"),
		TableName:  "todo",
		config:     map[string]interface{}{},
	}
}

func (c *Config) GetConfig() map[string]interface{} {
	return c.config
}

func (c *Config) GetDbName() string {
	return c.DbName
}

func (c *Config) GetTableName() string {
	return c.TableName
}

func (c *Config) GetConfigPath() string {
	if c.ConfigPath == "" {
		currentUser, err := user.Current()
		if err != nil {
			fmt.Println("Error getting current user: ", err)
			os.Exit(1)
		}

		c.ConfigPath = filepath.Join(currentUser.HomeDir, "/.todo/config.json")
	}
	return c.ConfigPath
}

func readJson(path string) (map[string]interface{}, error) {
	// Read json file
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	// Read json file into byte array
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	var x map[string]interface{}
	err = json.Unmarshal(byteValue, &x)

	if err != nil {
		return nil, err
	}
	return x, nil
}

func ReadConfig(filePath string) (*Config, error) {
	c := NewConfig()
	json, err := readJson(filePath)
	if err != nil {
		return nil, err
	}
	var found bool
	c.ConfigPath, found = json["ConfigPath"].(string)
	if !found {
		c.ConfigPath = filePath
	}
	c.DbName, found = json["DbName"].(string)
	if !found {
		c.DbName = filepath.Join(filepath.Dir(filePath), "todo.db")
	}
	c.TableName, found = json["TableName"].(string)
	if !found {
		c.TableName = "todo"
	}
	c.config, found = json["config"].(map[string]interface{})
	if !found {
		c.config = map[string]interface{}{}
	}
	return c, nil
}

func (c *Config) WriteConfig() error {
	jsonStr, err := json.Marshal(c)
	if err != nil {
		return err
	}

	// Write the config to the file
	// If the path in c.ConfigPath does not exist, create it
	if _, err := os.Stat(c.ConfigPath); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(c.ConfigPath), 0755)
		if err != nil {
			return err
		}
	}

	err = ioutil.WriteFile(c.ConfigPath, jsonStr, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) View() error {
	fmt.Println("Config file: ", c.ConfigPath)
	fmt.Println("DbName: ", c.DbName)
	fmt.Println("TableName: ", c.TableName)
	fmt.Println("Config: ", c.config)
	return nil
}

func (c *Config) Set(key string, value string) error {
	c.config[key] = value
	return nil
}

func (c *Config) Delete(key string) error {
	c.config[key] = ""
	return nil
}

func (c *Config) CreateDefaultConfig() error {
	currentUser, err := user.Current()
	if err != nil {
		return err
	}

	c.ConfigPath = filepath.Join(currentUser.HomeDir, "/.todo/config.json")
	// Create this filepath and file if it does not exist
	if _, err := os.Stat(c.ConfigPath); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(c.ConfigPath), 0755)
		if err != nil {
			return err
		}
		_, err = os.Create(c.ConfigPath)
		if err != nil {
			return err
		}
	} else {
		// Read the config file
		_, err := ReadConfig(c.ConfigPath)
		if err != nil {
			return err
		}
		return nil
	}

	// escaped string config path
	escapedConfigPath := filepath.Join(currentUser.HomeDir, "/.todo/config.json")
	// replace / with //
	escapedConfigPath = filepath.ToSlash(escapedConfigPath)
	escapedDbName := filepath.Join(currentUser.HomeDir, "/.todo/todo.db")
	// replace / with //
	escapedDbName = filepath.ToSlash(escapedDbName)

	defaultConfig := fmt.Sprintf(`{
		"ConfigPath": "%s",
		"DbName": "%s",
		"TableName": "todo",
		"config": {
			"key1": "value1",
			"key2": "value2"
		}
	}`, escapedConfigPath, escapedDbName)
	// Write the default config to the file
	err = os.WriteFile(c.ConfigPath, []byte(defaultConfig), 0644)
	if err != nil {
		return err
	}
	return nil
}

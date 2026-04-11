package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const configFileName = ".gatorconfig.json"

// represents the json file structure of the RSS
type Config struct {
	URL             string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func NewConfig() *Config {
	c := Config{}
	return &c
}

func getConfigFilePath() (string, error) {
	// returns the path to the config file
	homeDir, err := os.UserHomeDir()
	address := filepath.Join(homeDir, configFileName)
	return address, err
}

func Read() (Config, error) {
	// reads in the config file

	address, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	/*
		// TODO: check whether you need this or not, ask Boots when you're done

		file, err := os.Open(address)
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()
	*/
	body, err := os.ReadFile(address)
	if err != nil {
		return Config{}, err
	}

	// unmarshall the json data
	var config Config
	err = json.Unmarshal(body, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

//func Write(cfg Config) error {
// writes to the config file
//}

func (c *Config) SetUser(user string) error {
	// set the username in the config
	c.CurrentUserName = user

	address, err := getConfigFilePath()
	if err != nil {
		return err
	}

	// update the config file to reflect the username change
	// this marshalls config into json format
	jsonData, err := json.Marshal(*c)
	if err != nil {
		return err
	}

	err = os.WriteFile(address, jsonData, 0666)

	return err
}

func (c *Config) Prettify() (string, error) {
	jsonData, err := json.Marshal(*c)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil

}

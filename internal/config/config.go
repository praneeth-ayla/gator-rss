package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	config := Config{}

	filePath, err := getConfigFilePath()
	if err != nil {
		return config, fmt.Errorf("unable to get home dir")
	}
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return config, fmt.Errorf("unable to read file")
	}

	json.Unmarshal(fileData, &config)

	return config, nil
}

func (cnf *Config) SetUser(userName string) error {
	cnf.CurrentUserName = userName

	filePath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("unable to get config path")
	}

	data, err := json.Marshal(cnf)
	if err != nil {
		return fmt.Errorf("unable to marshal data")
	}
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("unable to update file data")
	}

	return nil

}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	filePath := path.Join(homeDir, ".gatorconfig.json")
	return filePath, nil

}

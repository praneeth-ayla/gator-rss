package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
)

// Config holds application configuration settings.
type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// Read reads the application configuration from a JSON file.
func Read() (Config, error) {
	config := Config{}

	// Get the path to the configuration file.
	filePath, err := getConfigFilePath()
	if err != nil {
		return config, fmt.Errorf("unable to get home dir")
	}
	// Read the file content.
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return config, fmt.Errorf("unable to read file")
	}

	// Unmarshal JSON data into the Config struct.
	json.Unmarshal(fileData, &config)

	return config, nil
}

// SetUser sets the current user name in the configuration and persists it to disk.
func (cnf *Config) SetUser(userName string) error {
	cnf.CurrentUserName = userName

	// Get the path to the configuration file.
	filePath, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("unable to get config path")
	}

	// Marshal the updated configuration to JSON.
	data, err := json.Marshal(cnf)
	if err != nil {
		return fmt.Errorf("unable to marshal data")
	}
	// Write the updated configuration back to the file.
	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("unable to update file data")
	}

	return nil

}

// getConfigFilePath returns the full path to the application's configuration file.
func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	filePath := path.Join(homeDir, ".gatorconfig.json")
	return filePath, nil

}

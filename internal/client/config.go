package client

import (
	"encoding/json"
	"os"
	"path/filepath"

	"efrc/internal/wg"
)

type Config struct {
	ServerURL  string `json:"server_url"`
	NetworkKey string `json:"network_key"`
	DeviceId   string `json:"device_id"`
	DeviceName string `json:"device_name"`
	PrivateKey string `json:"private_key"`
	PublicKey  string `json:"public_key"`
	VirtualIP  string `json:"virtual_ip"`
}

// ConfigDir returns ~/.efrc
func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(home, ".efrc"), nil
}

// ConfigPath returns ~/.efrc/config.json
func ConfigPath() (string, error) {
	configDir, err := ConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "config.json"), nil
}

// LoadConfig reads config from file
func LoadConfig() (*Config, error) {
	path, err := ConfigPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config

	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// Save writes config to file
func (c *Config) Save() error {
	path, err := ConfigPath()
	if err != nil {
		return err
	}

	// Create directory if not exists
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}

	// Marshal to pretty JSON
	data, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}

	// Write file
	return os.WriteFile(path, data, 0o600)
}

// InitConfig creates new config with generated keys
func InitConfig(name string) (*Config, error) {
	privateKey, publicKey, err := wg.GenerateKeyPair()
	if err != nil {
		return nil, err
	}

	config := &Config{
		DeviceName: name,
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}

	if err := config.Save(); err != nil {
		return nil, err
	}

	return config, nil
}

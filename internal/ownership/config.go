package ownership

import (
	"encoding/json"
	"fmt"
	"os"
)

// basically this file is doing : mapping files to teams using simple path rules, so we can see who owns what and track where conflicts are coming from
type OwnersConfig struct {
	Owners map[string]string `json:"owners"`
}

// LoadConfig : reads .gitresolve/owners.json
func LoadConfig(repoPath string) (*OwnersConfig, error) {
	configPath := repoPath + "/.gitresolve/owners.json"

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &OwnersConfig{Owners: make(map[string]string)}, nil
		}
		return nil, fmt.Errorf("LoadConfig: %w", err)
	}

	var config OwnersConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("LoadConfig: parsing owners.json: %w", err)
	}

	return &config, nil
}

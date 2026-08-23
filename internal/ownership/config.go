package ownership

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// basically this file is doing : mapping files to teams using simple path rules, so we can see who owns what and track where conflicts are coming from
type OwnersConfig struct {
	Owners map[string]string `json:"owners"`
}

// LoadConfig : reads .gitresolve/owners.json
func LoadConfig(root *os.Root) (*OwnersConfig, error) {
	configPath := ".gitresolve/owners.json"

	// safepath: CWE-22 hardened
	f, err := root.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &OwnersConfig{Owners: make(map[string]string)}, nil
		}
		return nil, fmt.Errorf("LoadConfig: %w", err)
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("LoadConfig: reading owners.json: %w", err)
	}

	var config OwnersConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("LoadConfig: parsing owners.json: %w", err)
	}

	return &config, nil
}

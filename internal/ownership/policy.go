package ownership

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

const (
	PolicyAuto       = "auto"
	PolicyStrict     = "strict"
	PolicyBalanced   = "balanced"
	PolicyAggressive = "aggressive"
)

type PolicyConfig struct {
	DefaultProfile string            `json:"default"`
	PathProfiles   map[string]string `json:"path_profiles"`
	TeamProfiles   map[string]string `json:"team_profiles"`
}

func IsValidPolicyProfile(profile string) bool {
	switch strings.ToLower(strings.TrimSpace(profile)) {
	case PolicyAuto, PolicyStrict, PolicyBalanced, PolicyAggressive:
		return true
	default:
		return false
	}
}

func normalizePolicyProfile(profile string) string {
	p := strings.ToLower(strings.TrimSpace(profile))
	if p == "" {
		return PolicyBalanced
	}
	if !IsValidPolicyProfile(p) {
		return PolicyBalanced
	}
	return p
}

func LoadPolicyConfig(repoPath string) (*PolicyConfig, error) {
	configPath := repoPath + "/.gitresolve/policy.json"

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &PolicyConfig{
				DefaultProfile: PolicyBalanced,
				PathProfiles:   make(map[string]string),
				TeamProfiles:   make(map[string]string),
			}, nil
		}
		return nil, fmt.Errorf("LoadPolicyConfig: %w", err)
	}

	var cfg PolicyConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("LoadPolicyConfig: parsing policy.json: %w", err)
	}

	cfg.DefaultProfile = normalizePolicyProfile(cfg.DefaultProfile)
	if cfg.PathProfiles == nil {
		cfg.PathProfiles = make(map[string]string)
	}
	if cfg.TeamProfiles == nil {
		cfg.TeamProfiles = make(map[string]string)
	}

	for k, v := range cfg.PathProfiles {
		cfg.PathProfiles[k] = normalizePolicyProfile(v)
	}
	for k, v := range cfg.TeamProfiles {
		cfg.TeamProfiles[k] = normalizePolicyProfile(v)
	}

	return &cfg, nil
}

func ResolvePolicyProfile(repoPath, filePath, explicitProfile string) (string, error) {
	explicit := normalizePolicyProfile(explicitProfile)
	if explicit != PolicyAuto {
		return explicit, nil
	}

	cfg, err := LoadPolicyConfig(repoPath)
	if err != nil {
		return "", err
	}

	bestPrefix := ""
	selected := ""
	for prefix, profile := range cfg.PathProfiles {
		if strings.HasPrefix(filePath, prefix) && len(prefix) > len(bestPrefix) {
			bestPrefix = prefix
			selected = profile
		}
	}
	if selected != "" {
		return selected, nil
	}

	owners, err := LoadConfig(repoPath)
	if err != nil {
		return "", err
	}
	team := CheckOwnership(owners, filePath)
	if team != "" {
		if p, ok := cfg.TeamProfiles[team]; ok {
			return p, nil
		}
	}

	return cfg.DefaultProfile, nil
}

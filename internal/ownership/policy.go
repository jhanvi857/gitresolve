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

type PolicyResolution struct {
	RequestedProfile string `json:"requested_profile"`
	ResolvedProfile  string `json:"resolved_profile"`
	Source           string `json:"source"`
	MatchedPath      string `json:"matched_path,omitempty"`
	MatchedTeam      string `json:"matched_team,omitempty"`
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
	resolution, err := ResolvePolicy(repoPath, filePath, explicitProfile)
	if err != nil {
		return "", err
	}
	return resolution.ResolvedProfile, nil
}

func ResolvePolicy(repoPath, filePath, explicitProfile string) (*PolicyResolution, error) {
	requested := normalizePolicyProfile(explicitProfile)
	if requested != PolicyAuto {
		return &PolicyResolution{
			RequestedProfile: requested,
			ResolvedProfile:  requested,
			Source:           "explicit",
		}, nil
	}

	cfg, err := LoadPolicyConfig(repoPath)
	if err != nil {
		return nil, err
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
		return &PolicyResolution{
			RequestedProfile: requested,
			ResolvedProfile:  selected,
			Source:           "path",
			MatchedPath:      bestPrefix,
		}, nil
	}

	owners, err := LoadConfig(repoPath)
	if err != nil {
		return nil, err
	}
	team := CheckOwnership(owners, filePath)
	if team != "" {
		if profile, ok := cfg.TeamProfiles[team]; ok {
			return &PolicyResolution{
				RequestedProfile: requested,
				ResolvedProfile:  profile,
				Source:           "team",
				MatchedTeam:      team,
			}, nil
		}
	}

	return &PolicyResolution{
		RequestedProfile: requested,
		ResolvedProfile:  cfg.DefaultProfile,
		Source:           "default",
	}, nil
}

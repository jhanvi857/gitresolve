package ownership

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
)

const (
	PolicyAuto       = "auto"
	PolicyStrict     = "strict"
	PolicyBalanced   = "balanced"
	PolicyAggressive = "aggressive"
)

var (
	ErrPolicyTooLarge       = errors.New("policy.json exceeds maximum size of 1MB")
	ErrPolicyTooManyEntries = errors.New("policy.json contains too many entries")
)

func ErrPolicyUnknownKey(key string) error {
	return fmt.Errorf("unknown top-level key in policy.json: %s", key)
}

func ErrPolicyInvalidProfile(val string) error {
	return fmt.Errorf("invalid policy profile: %s", val)
}

type PolicyConfig struct {
	DefaultProfile string            `json:"default"`
	PathProfiles   map[string]string `json:"path_profiles"`
	TeamProfiles   map[string]string `json:"team_profiles"`
	SortedPathKeys []string          `json:"-"`
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

func LoadPolicyConfig(root *os.Root) (*PolicyConfig, error) {
	configPath := ".gitresolve/policy.json"

	// safepath: CWE-22 hardened
	stat, err := root.Stat(configPath)
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

	if stat.Size() > 1024*1024 {
		return nil, ErrPolicyTooLarge
	}

	// safepath: CWE-22 hardened
	f, err := root.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var cfg PolicyConfig
	dec := json.NewDecoder(f)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&cfg); err != nil {
		if strings.Contains(err.Error(), "unknown field") {
			parts := strings.Split(err.Error(), "\"")
			if len(parts) >= 2 {
				return nil, ErrPolicyUnknownKey(parts[1])
			}
		}
		return nil, fmt.Errorf("LoadPolicyConfig: parsing policy.json: %w", err)
	}

	if cfg.DefaultProfile != "" && !IsValidPolicyProfile(cfg.DefaultProfile) {
		return nil, ErrPolicyInvalidProfile(cfg.DefaultProfile)
	}
	cfg.DefaultProfile = normalizePolicyProfile(cfg.DefaultProfile)

	if len(cfg.PathProfiles) > 500 {
		return nil, ErrPolicyTooManyEntries
	}
	cfg.SortedPathKeys = make([]string, 0, len(cfg.PathProfiles))
	for k, v := range cfg.PathProfiles {
		if !IsValidPolicyProfile(v) {
			return nil, ErrPolicyInvalidProfile(v)
		}
		cfg.PathProfiles[k] = normalizePolicyProfile(v)
		cfg.SortedPathKeys = append(cfg.SortedPathKeys, k)
	}
	sort.Strings(cfg.SortedPathKeys)

	if len(cfg.TeamProfiles) > 100 {
		return nil, ErrPolicyTooManyEntries
	}
	for k, v := range cfg.TeamProfiles {
		if !IsValidPolicyProfile(v) {
			return nil, ErrPolicyInvalidProfile(v)
		}
		cfg.TeamProfiles[k] = normalizePolicyProfile(v)
	}

	return &cfg, nil
}

func ResolvePolicyProfile(root *os.Root, filePath, explicitProfile string) (string, error) {
	resolution, err := ResolvePolicy(root, filePath, explicitProfile)
	if err != nil {
		return "", err
	}
	return resolution.ResolvedProfile, nil
}

func ResolvePolicy(root *os.Root, filePath, explicitProfile string) (*PolicyResolution, error) {
	requested := normalizePolicyProfile(explicitProfile)
	if requested != PolicyAuto {
		return &PolicyResolution{
			RequestedProfile: requested,
			ResolvedProfile:  requested,
			Source:           "explicit",
		}, nil
	}

	cfg, err := LoadPolicyConfig(root)
	if err != nil {
		return nil, err
	}

	bestPrefix := ""
	selected := ""

	// Optimized prefix matching using binary search
	idx := sort.Search(len(cfg.SortedPathKeys), func(i int) bool {
		return cfg.SortedPathKeys[i] > filePath
	})

	for i := idx - 1; i >= 0; i-- {
		prefix := cfg.SortedPathKeys[i]
		if strings.HasPrefix(filePath, prefix) {
			selected = cfg.PathProfiles[prefix]
			bestPrefix = prefix
			break
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

	owners, err := LoadConfig(root)
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

package analysis

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

type StructuredMergeResult struct {
	Content      string
	Conflicts    []StructuredConflict
	HasConflicts bool
}

type StructuredConflict struct {
	Key        string
	BaseValue  interface{}
	OurValue   interface{}
	TheirValue interface{}
}

func MergeJSON(base, ours, theirs []byte) (StructuredMergeResult, error) {
	var baseMap map[string]interface{}
	var oursMap map[string]interface{}
	var theirsMap map[string]interface{}

	parseSnippet := func(data []byte, dest *map[string]interface{}) error {
		str := strings.TrimSpace(string(data))
		if len(str) == 0 {
			*dest = make(map[string]interface{})
			return nil
		}

		if err := json.Unmarshal(data, dest); err == nil {
			return nil
		}

		str = strings.TrimSuffix(str, ",")

		wrapped := []byte("{ " + str + " }")
		return json.Unmarshal(wrapped, dest)
	}

	if err := parseSnippet(base, &baseMap); err != nil {
		return StructuredMergeResult{}, fmt.Errorf("MergeJSON: parsing base: %w", err)
	}
	if err := parseSnippet(ours, &oursMap); err != nil {
		return StructuredMergeResult{}, fmt.Errorf("MergeJSON: parsing ours: %w", err)
	}
	if err := parseSnippet(theirs, &theirsMap); err != nil {
		return StructuredMergeResult{}, fmt.Errorf("MergeJSON: parsing theirs: %w", err)
	}
	merged, conflicts := mergeMap(baseMap, oursMap, theirsMap)
	output, err := json.MarshalIndent(merged, "", "  ")
	if err != nil {
		return StructuredMergeResult{}, fmt.Errorf("MergeJSON: marshaling result: %w", err)
	}

	outputStr := string(output)
	oursStr := strings.TrimSpace(string(ours))
	if !strings.HasPrefix(oursStr, "{") && !strings.HasPrefix(oursStr, "[") {
		firstBrace := strings.Index(outputStr, "{")
		lastBrace := strings.LastIndex(outputStr, "}")
		if firstBrace != -1 && lastBrace != -1 && firstBrace < lastBrace {
			outputStr = outputStr[firstBrace+1 : lastBrace]
			outputStr = strings.Trim(outputStr, "\n\r")
			if strings.HasSuffix(oursStr, ",") {
				outputStr += ","
			}
		}
	}

	return StructuredMergeResult{
		Content:      outputStr,
		Conflicts:    conflicts,
		HasConflicts: len(conflicts) > 0,
	}, nil
}

func MergeYAML(base, ours, theirs []byte) (StructuredMergeResult, error) {
	baseVal, err := parseYAMLAny(base)
	if err != nil {
		return StructuredMergeResult{}, fmt.Errorf("MergeYAML: parsing base: %w", err)
	}
	oursVal, err := parseYAMLAny(ours)
	if err != nil {
		return StructuredMergeResult{}, fmt.Errorf("MergeYAML: parsing ours: %w", err)
	}
	theirsVal, err := parseYAMLAny(theirs)
	if err != nil {
		return StructuredMergeResult{}, fmt.Errorf("MergeYAML: parsing theirs: %w", err)
	}

	merged, conflicts := mergeAny("<root>", baseVal, oursVal, theirsVal)
	output, err := yaml.Marshal(merged)
	if err != nil {
		return StructuredMergeResult{}, fmt.Errorf("MergeYAML: marshaling result: %w", err)
	}

	return StructuredMergeResult{
		Content:      string(output),
		Conflicts:    conflicts,
		HasConflicts: len(conflicts) > 0,
	}, nil
}

func MergeTOML(base, ours, theirs []byte) (StructuredMergeResult, error) {
	baseMap, err := parseTOMLDocument(base)
	if err != nil {
		return StructuredMergeResult{}, fmt.Errorf("MergeTOML: parsing base: %w", err)
	}
	oursMap, err := parseTOMLDocument(ours)
	if err != nil {
		return StructuredMergeResult{}, fmt.Errorf("MergeTOML: parsing ours: %w", err)
	}
	theirsMap, err := parseTOMLDocument(theirs)
	if err != nil {
		return StructuredMergeResult{}, fmt.Errorf("MergeTOML: parsing theirs: %w", err)
	}

	merged, conflicts := mergeMap(baseMap, oursMap, theirsMap)
	output, err := toml.Marshal(merged)
	if err != nil {
		return StructuredMergeResult{}, fmt.Errorf("MergeTOML: marshaling result: %w", err)
	}

	return StructuredMergeResult{
		Content:      string(output),
		Conflicts:    conflicts,
		HasConflicts: len(conflicts) > 0,
	}, nil
}

func mergeMap(base, ours, theirs map[string]interface{}) (map[string]interface{}, []StructuredConflict) {
	result := make(map[string]interface{})
	var conflicts []StructuredConflict
	allKeysSet := make(map[string]bool)
	for k := range base {
		allKeysSet[k] = true
	}
	for k := range ours {
		allKeysSet[k] = true
	}
	for k := range theirs {
		allKeysSet[k] = true
	}

	// Sort keys for deterministic output (Fix V5: same input -> same output)
	allKeys := make([]string, 0, len(allKeysSet))
	for k := range allKeysSet {
		allKeys = append(allKeys, k)
	}
	sort.Strings(allKeys)

	for _, key := range allKeys {
		baseVal, baseExists := base[key]
		ourVal, ourExists := ours[key]
		theirVal, theirExists := theirs[key]

		if !baseExists && ourExists && theirExists {
			if valuesEqual(ourVal, theirVal) {
				result[key] = ourVal
				continue
			}
			if ourMap, ok := toStringMap(ourVal); ok {
				if theirMap, ok := toStringMap(theirVal); ok {
					nestedMerged, nestedConflicts := mergeMap(make(map[string]interface{}), ourMap, theirMap)
					result[key] = nestedMerged
					conflicts = append(conflicts, nestedConflicts...)
					continue
				}
			}
			if ourArr, ok := toInterfaceSlice(ourVal); ok {
				if theirArr, ok := toInterfaceSlice(theirVal); ok {
					nestedMerged, nestedConflicts := mergeArray(nil, ourArr, theirArr, key)
					result[key] = nestedMerged
					conflicts = append(conflicts, nestedConflicts...)
					continue
				}
			}
			conflicts = append(conflicts, StructuredConflict{
				Key:        key,
				BaseValue:  nil,
				OurValue:   ourVal,
				TheirValue: theirVal,
			})
			continue
		}

		// Case 1: Added in ours only
		if ourExists && !baseExists && !theirExists {
			result[key] = ourVal
			continue
		}

		// Case 2: Added in theirs only
		if theirExists && !baseExists && !ourExists {
			result[key] = theirVal
			continue
		}

		// Case 3: Deleted in one, but modified in other
		if !ourExists && baseExists && theirExists {
			if valuesEqual(baseVal, theirVal) {
				// both sides deleted it (conceptually, or base==theirs and ours deleted)
				continue
			}
			// conflict: ours deleted, theirs modified
			conflicts = append(conflicts, StructuredConflict{
				Key:        key,
				BaseValue:  baseVal,
				OurValue:   nil,
				TheirValue: theirVal,
			})
			continue
		}
		if ourExists && baseExists && !theirExists {
			if valuesEqual(baseVal, ourVal) {
				continue
			}
			// conflict: theirs deleted, ours modified
			conflicts = append(conflicts, StructuredConflict{
				Key:        key,
				BaseValue:  baseVal,
				OurValue:   ourVal,
				TheirValue: nil,
			})
			continue
		}

		// Case 4: Both exist (or both were base)
		if valuesEqual(ourVal, baseVal) && valuesEqual(theirVal, baseVal) {
			result[key] = baseVal
			continue
		}

		if !valuesEqual(ourVal, baseVal) && valuesEqual(theirVal, baseVal) {
			result[key] = ourVal
			continue
		}

		if valuesEqual(ourVal, baseVal) && !valuesEqual(theirVal, baseVal) {
			result[key] = theirVal
			continue
		}

		// Both modified!
		if !valuesEqual(ourVal, baseVal) && !valuesEqual(theirVal, baseVal) {
			if valuesEqual(ourVal, theirVal) {
				result[key] = ourVal
				continue
			}

			// Recursive merge for maps
			ourMap, ourIsMap := ourVal.(map[string]interface{})
			theirMap, theirIsMap := theirVal.(map[string]interface{})
			baseMap, baseIsMap := baseVal.(map[string]interface{})

			if ourIsMap && theirIsMap && baseIsMap {
				nestedMerged, nestedConflicts := mergeMap(baseMap, ourMap, theirMap)
				result[key] = nestedMerged
				conflicts = append(conflicts, nestedConflicts...)
				continue
			}

			// Array merge attempt
			ourArr, ourIsArr := ourVal.([]interface{})
			theirArr, theirIsArr := theirVal.([]interface{})
			baseArr, baseIsArr := baseVal.([]interface{})

			if ourIsArr && theirIsArr && baseIsArr {
				nestedMerged, nestedConflicts := mergeArray(baseArr, ourArr, theirArr, key)
				result[key] = nestedMerged
				conflicts = append(conflicts, nestedConflicts...)
				continue
			}

			// Scalar conflict
			conflicts = append(conflicts, StructuredConflict{
				Key:        key,
				BaseValue:  baseVal,
				OurValue:   ourVal,
				TheirValue: theirVal,
			})
		}
	}

	return result, conflicts
}

func mergeAny(key string, base, ours, theirs interface{}) (interface{}, []StructuredConflict) {
	oursMap, oursIsMap := toStringMap(ours)
	theirsMap, theirsIsMap := toStringMap(theirs)
	if oursIsMap && theirsIsMap {
		baseMap, baseIsMap := toStringMap(base)
		if !baseIsMap {
			baseMap = make(map[string]interface{})
		}
		merged, conflicts := mergeMap(baseMap, oursMap, theirsMap)
		return merged, conflicts
	}

	oursArr, oursIsArr := toInterfaceSlice(ours)
	theirsArr, theirsIsArr := toInterfaceSlice(theirs)
	if oursIsArr && theirsIsArr {
		baseArr, baseIsArr := toInterfaceSlice(base)
		if !baseIsArr {
			baseArr = nil
		}
		merged, conflicts := mergeArray(baseArr, oursArr, theirsArr, key)
		return merged, conflicts
	}

	if valuesEqual(ours, theirs) {
		return ours, nil
	}
	if valuesEqual(ours, base) {
		return theirs, nil
	}
	if valuesEqual(theirs, base) {
		return ours, nil
	}

	return nil, []StructuredConflict{{
		Key:        key,
		BaseValue:  base,
		OurValue:   ours,
		TheirValue: theirs,
	}}
}

func mergeArray(base, ours, theirs []interface{}, key string) ([]interface{}, []StructuredConflict) {
	// If identical, return either
	if valuesEqual(ours, theirs) {
		return ours, nil
	}

	if base == nil {
		seen := make([]interface{}, 0, len(ours)+len(theirs))
		appendUnique := func(arr []interface{}) {
			for _, v := range arr {
				exists := false
				for _, seenV := range seen {
					if valuesEqual(seenV, v) {
						exists = true
						break
					}
				}
				if !exists {
					seen = append(seen, v)
				}
			}
		}
		appendUnique(ours)
		appendUnique(theirs)
		return seen, nil
	}

	// Semantic merge: combine additions and deletions
	contains := func(arr []interface{}, item interface{}) bool {
		for _, a := range arr {
			if valuesEqual(a, item) {
				return true
			}
		}
		return false
	}

	var conflicts []StructuredConflict

	deletedByUs := make(map[int]bool)
	deletedByThem := make(map[int]bool)
	for i, b := range base {
		if !contains(ours, b) {
			deletedByUs[i] = true
		}
		if !contains(theirs, b) {
			deletedByThem[i] = true
		}
	}

	// Detect conflicting deletions: one side deleted, other side kept
	for i, b := range base {
		if deletedByUs[i] && !deletedByThem[i] && !contains(ours, b) {
			// ours deleted this element, theirs kept it — potential conflict
			// For simple values we auto-resolve by respecting deletion,
			// but for complex values (maps) we flag it.
			if _, isMap := b.(map[string]interface{}); isMap {
				conflicts = append(conflicts, StructuredConflict{
					Key:        fmt.Sprintf("%s[%d]", key, i),
					BaseValue:  b,
					OurValue:   nil,
					TheirValue: b,
				})
			}
		}
		if deletedByThem[i] && !deletedByUs[i] && !contains(theirs, b) {
			if _, isMap := b.(map[string]interface{}); isMap {
				conflicts = append(conflicts, StructuredConflict{
					Key:        fmt.Sprintf("%s[%d]", key, i),
					BaseValue:  b,
					OurValue:   b,
					TheirValue: nil,
				})
			}
		}
	}

	var result []interface{}
	for i, b := range base {
		if !deletedByUs[i] && !deletedByThem[i] {
			result = append(result, b)
		}
	}
	for _, o := range ours {
		if !contains(base, o) {
			result = append(result, o)
		}
	}
	for _, t := range theirs {
		if !contains(base, t) && !contains(ours, t) {
			result = append(result, t)
		}
	}

	return result, conflicts
}

func parseYAMLAny(data []byte) (interface{}, error) {
	if len(strings.TrimSpace(string(data))) == 0 {
		return nil, nil
	}
	var v interface{}
	if err := yaml.Unmarshal(data, &v); err != nil {
		return nil, err
	}
	return normalizeYAMLValue(v), nil
}

func normalizeYAMLValue(v interface{}) interface{} {
	switch typed := v.(type) {
	case map[string]interface{}:
		result := make(map[string]interface{}, len(typed))
		for k, val := range typed {
			result[k] = normalizeYAMLValue(val)
		}
		return result
	case map[interface{}]interface{}:
		result := make(map[string]interface{}, len(typed))
		for k, val := range typed {
			result[fmt.Sprint(k)] = normalizeYAMLValue(val)
		}
		return result
	case []interface{}:
		result := make([]interface{}, 0, len(typed))
		for _, val := range typed {
			result = append(result, normalizeYAMLValue(val))
		}
		return result
	default:
		return v
	}
}

func toStringMap(v interface{}) (map[string]interface{}, bool) {
	if v == nil {
		return nil, false
	}
	switch typed := v.(type) {
	case map[string]interface{}:
		return typed, true
	case map[interface{}]interface{}:
		converted := normalizeYAMLValue(typed)
		mapped, ok := converted.(map[string]interface{})
		return mapped, ok
	default:
		return nil, false
	}
}

func toInterfaceSlice(v interface{}) ([]interface{}, bool) {
	if v == nil {
		return nil, false
	}
	typed, ok := v.([]interface{})
	return typed, ok
}

func parseTOMLDocument(data []byte) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" {
		return result, nil
	}
	if err := toml.Unmarshal(data, &result); err == nil {
		return result, nil
	}

	// Fallback for snippet-like TOML chunks that are not full documents.
	parsed := parseSimpleTOMLMap(trimmed)
	if len(parsed) == 0 {
		return nil, fmt.Errorf("unsupported TOML snippet")
	}
	return parsed, nil
}

func parseSimpleTOMLMap(content string) map[string]interface{} {
	result := make(map[string]interface{})
	section := ""
	for _, raw := range strings.Split(content, "\n") {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section = strings.TrimSpace(strings.TrimSuffix(strings.TrimPrefix(line, "["), "]"))
			continue
		}
		if !strings.Contains(line, "=") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if idx := strings.Index(value, "#"); idx >= 0 {
			value = strings.TrimSpace(value[:idx])
		}
		if key == "" || value == "" {
			continue
		}
		fullKey := key
		if section != "" {
			fullKey = section + "." + key
		}
		result[fullKey] = value
	}

	if len(result) > 0 {
		ordered := make([]string, 0, len(result))
		for k := range result {
			ordered = append(ordered, k)
		}
		sort.Strings(ordered)
		normalized := make(map[string]interface{}, len(result))
		for _, k := range ordered {
			normalized[k] = result[k]
		}
		return normalized
	}

	return result
}

func valuesEqual(a, b interface{}) bool {
	// Optimization for nil
	if a == nil || b == nil {
		return a == b
	}
	aJSON, err := json.Marshal(a)
	if err != nil {
		return false
	}
	bJSON, err := json.Marshal(b)
	if err != nil {
		return false
	}
	return string(aJSON) == string(bJSON)
}

func IsStructuredFile(filePath string) bool {
	ext := strings.ToLower(filepath.Ext(filePath))
	return ext == ".json" || ext == ".yaml" || ext == ".yml" || ext == ".toml"
}

func IsCriticalFile(filePath string) bool {
	name := strings.ToLower(filepath.Base(filePath))
	return name == "package.json" ||
		name == "go.mod" ||
		name == "cargo.toml" ||
		name == "composer.json" ||
		name == "podfile" ||
		name == "yarn.lock" ||
		name == "package-lock.json"
}

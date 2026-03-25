package analysis

import (
	"encoding/json"
	"fmt"
	"strings"

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

	// json.Unmarshal converts raw JSON bytes into Go data structures
	// if the JSON is malformed this returns an error immediately
	if err := json.Unmarshal(base, &baseMap); err != nil {
		return StructuredMergeResult{}, fmt.Errorf("MergeJSON: parsing base: %w", err)
	}
	if err := json.Unmarshal(ours, &oursMap); err != nil {
		return StructuredMergeResult{}, fmt.Errorf("MergeJSON: parsing ours: %w", err)
	}
	if err := json.Unmarshal(theirs, &theirsMap); err != nil {
		return StructuredMergeResult{}, fmt.Errorf("MergeJSON: parsing theirs: %w", err)
	}
	merged, conflicts := mergeMap(baseMap, oursMap, theirsMap)
	output, err := json.MarshalIndent(merged, "", "  ")
	if err != nil {
		return StructuredMergeResult{}, fmt.Errorf("MergeJSON: marshaling result: %w", err)
	}

	return StructuredMergeResult{
		Content:      string(output),
		Conflicts:    conflicts,
		HasConflicts: len(conflicts) > 0,
	}, nil
}

func MergeYAML(base, ours, theirs []byte) (StructuredMergeResult, error) {
	var baseMap map[string]interface{}
	var oursMap map[string]interface{}
	var theirsMap map[string]interface{}

	if err := yaml.Unmarshal(base, &baseMap); err != nil {
		return StructuredMergeResult{}, fmt.Errorf("MergeYAML: parsing base: %w", err)
	}
	if err := yaml.Unmarshal(ours, &oursMap); err != nil {
		return StructuredMergeResult{}, fmt.Errorf("MergeYAML: parsing ours: %w", err)
	}
	if err := yaml.Unmarshal(theirs, &theirsMap); err != nil {
		return StructuredMergeResult{}, fmt.Errorf("MergeYAML: parsing theirs: %w", err)
	}

	merged, conflicts := mergeMap(baseMap, oursMap, theirsMap)
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
func mergeMap(base, ours, theirs map[string]interface{}) (map[string]interface{}, []StructuredConflict) {
	result := make(map[string]interface{})
	var conflicts []StructuredConflict
	allKeys := make(map[string]bool)
	for k := range base {
		allKeys[k] = true
	}
	for k := range ours {
		allKeys[k] = true
	}
	for k := range theirs {
		allKeys[k] = true
	}

	for key := range allKeys {
		baseVal, baseExists := base[key]
		ourVal, ourExists := ours[key]
		theirVal, theirExists := theirs[key]
		if ourExists && !baseExists && !theirExists {
			result[key] = ourVal
			continue
		}

		if theirExists && !baseExists && !ourExists {
			result[key] = theirVal
			continue
		}
		if !ourExists && baseExists && theirExists {
			if valuesEqual(baseVal, theirVal) {
				continue
			}
			conflicts = append(conflicts, StructuredConflict{
				Key:        key,
				BaseValue:  baseVal,
				OurValue:   nil,
				TheirValue: theirVal,
			})
			continue
		}

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

		if !valuesEqual(ourVal, baseVal) && !valuesEqual(theirVal, baseVal) {
			if valuesEqual(ourVal, theirVal) {
				result[key] = ourVal
				continue
			}
			ourMap, ourIsMap := ourVal.(map[string]interface{})
			theirMap, theirIsMap := theirVal.(map[string]interface{})
			baseMap, baseIsMap := baseVal.(map[string]interface{})

			if ourIsMap && theirIsMap && baseIsMap {
				nestedMerged, nestedConflicts := mergeMap(baseMap, ourMap, theirMap)
				result[key] = nestedMerged
				conflicts = append(conflicts, nestedConflicts...)
				continue
			}

			conflicts = append(conflicts, StructuredConflict{
				Key:        key,
				BaseValue:  baseVal,
				OurValue:   ourVal,
				TheirValue: theirVal,
			})
			result[key] = ourVal
		}
	}

	return result, conflicts
}

func valuesEqual(a, b interface{}) bool {
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
	return strings.HasSuffix(filePath, ".json") ||
		strings.HasSuffix(filePath, ".yaml") ||
		strings.HasSuffix(filePath, ".yml") ||
		strings.HasSuffix(filePath, ".toml")
}

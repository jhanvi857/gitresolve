package analysis

import (
	"encoding/json"
	"fmt"
	"path/filepath"
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

func MergeTOML(base, ours, theirs []byte) (StructuredMergeResult, error) {
	var baseMap map[string]interface{}
	var oursMap map[string]interface{}
	var theirsMap map[string]interface{}

	if err := toml.Unmarshal(base, &baseMap); err != nil {
		return StructuredMergeResult{}, fmt.Errorf("MergeTOML: parsing base: %w", err)
	}
	if err := toml.Unmarshal(ours, &oursMap); err != nil {
		return StructuredMergeResult{}, fmt.Errorf("MergeTOML: parsing ours: %w", err)
	}
	if err := toml.Unmarshal(theirs, &theirsMap); err != nil {
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

func mergeArray(base, ours, theirs []interface{}, key string) ([]interface{}, []StructuredConflict) {
	// If identical, return either
	if valuesEqual(ours, theirs) {
		return ours, nil
	}

	// If one is unchanged from base, return the other
	if valuesEqual(ours, base) {
		return theirs, nil
	}
	if valuesEqual(theirs, base) {
		return ours, nil
	}

	// Make a copy of ours
	merged := make([]interface{}, len(ours))
	copy(merged, ours)

	contains := func(arr []interface{}, item interface{}) bool {
		for _, a := range arr {
			if valuesEqual(a, item) {
				return true
			}
		}
		return false
	}

	// Append additions seamlessly
	for _, item := range theirs {
		if !contains(base, item) && !contains(ours, item) {
			merged = append(merged, item)
		}
	}

	if len(ours) < len(base) || len(theirs) < len(base) {
		return merged, []StructuredConflict{{Key: key, BaseValue: base, OurValue: ours, TheirValue: theirs}}
	}

	return merged, nil
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

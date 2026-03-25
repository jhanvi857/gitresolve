package ownership

import "strings"

// checking the ownership of file.
func CheckOwnership(config *OwnersConfig, filePath string) string {
	for pattern, team := range config.Owners {
		if strings.HasPrefix(filePath, pattern) {
			return team
		}
	}
	return ""
}

// is there exists a cross ownership for the file like if two people have changed same file or not.
func IsCrossOwnership(config *OwnersConfig, filePath, changedBy string) bool {
	owner := CheckOwnership(config, filePath)
	if owner == "" {
		return false
	}
	return owner != changedBy
}

package safety

// purpose : Never write directly to a file. Write to a temp file first, then rename.
// so if program crashes mid-write the original file is untouched.
import (
	"fmt"
	"os"
)

func getDirOf(filePath string) string {
	for i := len(filePath) - 1; i >= 0; i-- {
		if filePath[i] == '/' || filePath[i] == '\\' {
			return filePath[:i]
		}
	}
	return "."
}

// normally writing with os.WriteFile and if program crashes in middle then data lost, file corruption will be there.
// why need of writing ? : like suppose gitresolve was resolving conflict in some file and program crashed or laptop issue then half one'll have resolved content and other half will be gone

// that's why what to do is that write in temp file :
// 1. create a TEMP file (.gitresolve-tmp-abc123)
// 2. write line 1 into temp
// 3. write line 2 into temp
// 4. write line 3 into temp   -> crash here? src/config.js is 100% untouched
// 5. write line 4 into temp
// 6. os.Rename(temp -> src/config.js)  -> this one operation is instant so no issue with data.
func writeAtomic(targetPath string, content []byte) error {
	// creating temp file in same directory
	tmpFile, err := os.CreateTemp(getDirOf(targetPath), ".gitresolve-tmp-*")
	if err != nil {
		return fmt.Errorf("WriteAtomic : creating temp file : %w", err)
	}
	tmpPath := tmpFile.Name()
	defer func() {
		tmpFile.Close()
		os.Remove(tmpPath)
	}()
	_, err = tmpFile.Write(content)
	if err != nil {
		return fmt.Errorf("WriteAtomic: writing to temp file: %w", err)
	}

	err = tmpFile.Sync()
	if err != nil {
		return fmt.Errorf("WriteAtomic: syncing temp file: %w", err)
	}

	tmpFile.Close()
	err = os.Rename(tmpPath, targetPath)
	if err != nil {
		return fmt.Errorf("WriteAtomic: renaming to target: %w", err)
	}

	return nil
}

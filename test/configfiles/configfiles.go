package configfiles

import "path/filepath"

func GetConfigFilesDir() string {
	return filepath.Join("test", "configfiles")
}

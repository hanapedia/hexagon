package generate

import (
	"fmt"
	"os"
)

func formatManifest(manifest []byte) string {
	return fmt.Sprintf("---\n%s\n", manifest)
}

// createFile create and open file in append
func createFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
}

func (mg ManifestGenerator) getFilePath(name, kind string) string {
	return fmt.Sprintf("%s/%s-%s.yaml", mg.Output, name, kind)
}

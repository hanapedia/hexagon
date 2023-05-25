package generate

import (
	"fmt"
	"os"
)

func formatManifest(manifest []byte) string {
	return fmt.Sprintf("---\n%s\n", manifest)
}

func createFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
}

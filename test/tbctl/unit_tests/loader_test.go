package unit_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hanapedia/the-bench/internal/tbctl/loader"
)

func TestGetYAMLFiles(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "yamlfiles_test")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test directory structure with YAML files
	yamlTestFiles := []string{
		"file1.yaml",
		"file2.yml",
		"subdir1/file3.yaml",
		"subdir1/file4.yml",
		"subdir2/file5.yaml",
		"subdir2/file6.txt",
	}

	for _, file := range yamlTestFiles {
		filePath := filepath.Join(tempDir, file)
		dir := filepath.Dir(filePath)

		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}

		if !strings.HasSuffix(file, ".txt") {
			if err := ioutil.WriteFile(filePath, []byte("test"), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}
		}
	}

	// Run the GetYAMLFiles function
	result, err := loader.GetYAMLFiles(tempDir)
	if err != nil {
		t.Fatalf("GetYAMLFiles failed: %v", err)
	}

	// Validate the result
	if len(result) != len(yamlTestFiles)-1 { // Minus the .txt file
		t.Errorf("Expected %d YAML files, got %d", len(yamlTestFiles)-1, len(result))
	}

	for _, file := range result {
		if !strings.Contains(file, ".yaml") && !strings.Contains(file, ".yml") {
			t.Errorf("Non-YAML file detected: %s", file)
		}
	}
}

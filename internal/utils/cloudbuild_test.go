package util

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/graph/encoding/dot"
)

func init() {
	// Change the current working dir to root dir
	os.Chdir("../../")
}

func TestYamlToDAG(t *testing.T) {
	var testFiles []string
	err := filepath.Walk("./test/fixtures/cloudbuild", func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".yaml" {
			testFiles = append(testFiles, path)
		}
		return nil
	})
	if err != nil {
		assert.Error(t, err)
	}
	for _, file := range testFiles {
		t.Run(fmt.Sprintf("TestYamlToDAG:%s", file), func(t *testing.T) {
			cloudBuild, err := ParseYaml(file)
			assert.Empty(t, err)
			dag := BuildStepsToDAG(cloudBuild.Steps)
			result, err := dot.Marshal(dag, "test", "", "  ")
			assert.Empty(t, err)
			dotFilePath := getDotFilePath(file)
			expected, err := ioutil.ReadFile(dotFilePath)
			assert.Empty(t, err)
			assert.Equal(t, strings.Trim(string(expected), "\n"), strings.Trim(string(result), "\n"))
		})
	}
}

func getDotFilePath(filePath string) string {
	ext := path.Ext(filePath)
	filePath = filePath[0:len(filePath)-len(ext)] + ".dot"
	return strings.ReplaceAll(filePath, "test/fixtures/cloudbuild", "test/fixtures/dot")
}

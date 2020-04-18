package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	graphviz "github.com/goccy/go-graphviz"
	"github.com/stretchr/testify/assert"
)

var supportedFormat = []string{".yaml", ".yml", ".json"}

func init() {
	// Change the current working dir to root dir
	os.Chdir("../../")
}

func TestYamlToDAG(t *testing.T) {
	var testFiles []string
	cloudbuildPath := filepath.Join("./", "test", "fixtures", "cloudbuild")
	err := filepath.Walk(cloudbuildPath, func(path string, info os.FileInfo, err error) error {
		if Contains(supportedFormat, filepath.Ext(path)) {
			testFiles = append(testFiles, path)
		}
		return nil
	})
	if err != nil {
		assert.Error(t, err)
	}
	for _, file := range testFiles {
		t.Run(fmt.Sprintf("TestYamlToDAG:%s", file), func(t *testing.T) {
			g := graphviz.New()
			cloudBuild, err := ParseYaml(file)
			assert.Empty(t, err)
			graph := BuildStepsToDAG(cloudBuild.Steps)
			var buf bytes.Buffer
			err = g.Render(graph, "dot", &buf)
			assert.Empty(t, err)
			dotFilePath := getDotFilePath(file)
			expected, err := ioutil.ReadFile(dotFilePath)
			assert.Empty(t, err)
			assert.Equal(t, strings.ReplaceAll(string(expected), "\r\n", "\n"), buf.String())
		})
	}
}

func getDotFilePath(filePath string) string {
	ext := path.Ext(filePath)
	_, filename := filepath.Split(filePath)
	dotFile := filename[0:len(filename)-len(ext)] + ".dot"
	return filepath.Join("./", "test", "fixtures", "dot", dotFile)
}

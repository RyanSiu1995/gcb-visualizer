package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"testing"

	graphviz "github.com/goccy/go-graphviz"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Change the current working dir to root dir
	os.Chdir("../../")
}

func TestYamlToDAG(t *testing.T) {
	var testFiles []string
	cloudbuildPath := filepath.Join("./", "test", "fixtures", "cloudbuild")
	err := filepath.Walk(cloudbuildPath, func(path string, info os.FileInfo, err error) error {
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
			assert.Equal(t, string(expected), buf.String())
		})
	}
}

func getDotFilePath(filePath string) string {
	ext := path.Ext(filePath)
	filename := path.Base(filePath)
	dotFile := filename[0:len(filename)-len(ext)] + ".dot"
	return filepath.Join("./", "test", "fixtures", "dot", dotFile)
}

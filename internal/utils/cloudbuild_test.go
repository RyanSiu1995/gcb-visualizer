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
	"github.com/goccy/go-graphviz/cgraph"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

var supportedFormat = []string{".yaml", ".yml", ".json"}

type result struct {
	Nodes []string `yaml:"nodes"`
	Edges []struct {
		From string `yaml:"from"`
		To   string `yaml:"to"`
	} `yaml:"edges"`
}

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
			cloudBuild, err := ParseYaml(file)
			assert.Empty(t, err)
			graph := BuildStepsToDAG(cloudBuild.Steps)
			assert.Empty(t, err)

			expectedResult := getExpectedResult(file)
			isExpectedGraph(t, expectedResult, graph)
		})
	}
}

func getExpectedResult(filePath string) *result {
	ext := path.Ext(filePath)
	_, filename := filepath.Split(filePath)
	yamlFile := filename[0:len(filename)-len(ext)] + ".graph.yaml"
	fullPath := filepath.Join("./", "test", "fixtures", "graph", yamlFile)
	resultObj := result{}
	buf, _ := ioutil.ReadFile(fullPath)
	yaml.Unmarshal(buf, &resultObj)
	return &resultObj
}

func isExpectedGraph(t *testing.T, expected *result, actual *cgraph.Graph) {
	assert.Equalf(t, len(expected.Nodes), actual.NumberNodes(), "Should have the same number of nodes.")
	for _, node := range expected.Nodes {
		n, err := actual.Node(node)
		assert.Empty(t, err)
		assert.NotEmptyf(t, n, "Expected Node %s but not found in generated graph", node)
	}
	assert.Equalf(t, len(expected.Edges), actual.NumberEdges(), "Should have the same number of edges.")
	// FIXME This may not be too good to compare the edge. We may need a better way to compare them programmatically
	var buf bytes.Buffer
	g := graphviz.New()
	g.Render(actual, "dot", &buf)
	actualOutput := buf.String()
	for _, edge := range expected.Edges {
		edge.From = sanitizeDotString(edge.From)
		edge.To = sanitizeDotString(edge.To)
		edgeInDot := fmt.Sprintf("%s -> %s", edge.From, edge.To)
		assert.Truef(t, strings.Contains(actualOutput, edgeInDot), "Should contain the edge \"%s\"\n Actual: %s", edgeInDot, actualOutput)
	}
}

func sanitizeDotString(s string) string {
	if strings.Contains(s, " ") {
		return fmt.Sprintf("\"%s\"", s)
	}
	return s
}

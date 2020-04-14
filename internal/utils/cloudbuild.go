package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	yamlUtil "github.com/ghodss/yaml"
	graphviz "github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
	log "github.com/sirupsen/logrus"
	"github.com/skratchdot/open-golang/open"
	cloudbuild "google.golang.org/api/cloudbuild/v1"
)

func init() {
	if _, exist := os.LookupEnv("DEBUG"); exist {
		log.SetLevel(log.DebugLevel)
	}
}

// ParseYaml takes a string of file path and returns the cloud build object
func ParseYaml(filePath string) (*cloudbuild.Build, error) {
	yamlFileInByte, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	jsonFileInByte, err := yamlUtil.YAMLToJSON(yamlFileInByte)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	var build cloudbuild.Build
	if err := json.Unmarshal(jsonFileInByte, &build); err != nil {
		fmt.Println(err.Error())
	}

	return &build, nil
}

// BuildStepsToDAG takes the list of build steps and build a directed acyclic graph
func BuildStepsToDAG(steps []*cloudbuild.BuildStep) *cgraph.Graph {
	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		log.Fatal(err.Error())
	}

	var nodeList []*cgraph.Node
	for idx, step := range steps {
		name := step.Id
		if name == "" {
			name = fmt.Sprintf("Step %d", idx)
		}
		node, err := graph.CreateNode(name)
		if err != nil {
			log.Fatal(err.Error())
		}
		nodeList = append(nodeList, node)
	}
	mapping := buildIDToIdxMap(steps)
	for idx := range steps {
		handleWaitFor(steps, idx, mapping, graph, nodeList)
	}
	return graph
}

// Visualize is a high level API to show the graph
func Visualize(graph *cgraph.Graph) error {
	dir, err := ioutil.TempDir("", "gcb-temp")
	if err != nil {
		return err
	}
	filename := filepath.Join(dir, "temp.png")

	g := graphviz.New()

	if err := g.RenderFilename(graph, graphviz.PNG, filename); err != nil {
		return err
	}

	err = open.Run(filename)
	if err != nil {
		return err
	}
	return nil
}

func handleWaitFor(steps []*cloudbuild.BuildStep, idx int, mapping map[string]int, graph *cgraph.Graph, nodes []*cgraph.Node) {
	waitFor := steps[idx].WaitFor
	if len(waitFor) == 1 {
		if !contains(waitFor, "-") {
			log.Debugf("Node %d is handled with normal waitFor case with length 1...", idx)
			from := mapping[waitFor[0]]
			log.Debugf("Node %d has waitFor %s. Adding the edge from %d to %d...", idx, waitFor[0], from, idx)
			graph.CreateEdge("", nodes[from], nodes[idx])
		} else {
			// If the "-" in the start, it will be nothing to do
			if idx != 0 {
				log.Debugf("Node %d are handled with \"-\" case...", idx)
				// Search for next node without waitFor
				for i := idx; i < len(steps); i++ {
					if len(steps[i].WaitFor) == 0 {
						log.Debugf("The next node with waitFor for immediately started node %d is node %d. Adding the edge from %d to %d...", idx, i, idx, i)
						graph.CreateEdge("", nodes[idx], nodes[i])
						break
					}
				}
			}
		}
	} else if len(waitFor) != 0 {
		log.Debugf("Node %d is handled with normal waitFor case...", idx)
		// Handle all normal cases
		for _, item := range waitFor {
			from := mapping[item]
			log.Debugf("Node %d has waitFor %s. Adding the edge from %d to %d...", idx, item, from, idx)
			graph.CreateEdge("", nodes[from], nodes[idx])
		}
	} else {
		log.Debugf("Node %d is handled with no waitFor case...", idx)
		// Handle all cases without waitFor
		for i := idx - 1; i >= 0; i-- {
			if len(steps[i].WaitFor) != 0 {
				if isLastNode(steps, mapping, idx, i) {
					log.Debugf("Found the last node with waitFor before node %d is node %d. Adding the edge from %d to %d...", idx, i, i, idx)
					graph.CreateEdge("", nodes[i], nodes[idx])
				}
			} else {
				log.Debugf("Last node without waitFor for node %d is node %d. Adding the edge from %d to %d...", idx, i, i, idx)
				graph.CreateEdge("", nodes[i], nodes[idx])
				// If we encounter the last node without WaitFor, all the rest of cases should be routed to this node instead
				break
			}
		}
	}
}

func isLastNode(steps []*cloudbuild.BuildStep, mapping map[string]int, threshold int, idx int) bool {
	id := steps[idx].Id
	for i := idx; i < threshold; i++ {
		if contains(steps[i].WaitFor, id) {
			return false
		}
	}
	return true
}

func buildIDToIdxMap(steps []*cloudbuild.BuildStep) map[string]int {
	var mapping = make(map[string]int)
	for idx, step := range steps {
		if step.Id != "" {
			mapping[step.Id] = idx
		}
	}
	return mapping
}

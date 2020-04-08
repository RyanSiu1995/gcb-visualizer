package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	yamlUtil "github.com/ghodss/yaml"
	"gonum.org/v1/gonum/graph/simple"
	cloudbuild "google.golang.org/api/cloudbuild/v1"
)

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
func BuildStepsToDAG(steps []*cloudbuild.BuildStep) *simple.DirectedGraph {
	dag := simple.NewDirectedGraph()
	mapping := buildIDToIdxMap(steps)
	for idx := range steps {
		// Adding the nodes first
		dag.AddNode(simple.Node(idx))
		handleWaitFor(steps, idx, mapping, dag)
	}

	return dag
}

func handleWaitFor(steps []*cloudbuild.BuildStep, idx int, mapping map[string]int, dag *simple.DirectedGraph) {
	waitFor := steps[idx].WaitFor
	if len(waitFor) == 1 {
		if !contains(waitFor, "-") {
			from := mapping[waitFor[0]]
			dag.SetEdge(simple.Edge{F: simple.Node(from), T: simple.Node(idx)})
		} else {
			// If the "-" in the start, it will be nothing to do
			if idx != 0 {
				// Search for next node without waitFor
				for i := idx; i < len(steps); i++ {
					if len(steps[i].WaitFor) == 0 {
						dag.SetEdge(simple.Edge{F: simple.Node(idx), T: simple.Node(i)})
						break
					}
				}
			}
		}
	} else if len(waitFor) != 0 {
		// Handle all normal cases
		for _, item := range waitFor {
			from := mapping[item]
			dag.SetEdge(simple.Edge{F: simple.Node(from), T: simple.Node(idx)})
		}
	} else {
		// Handle all cases without waitFor
		for i := idx; idx < 0; i-- {
			if len(steps[i].WaitFor) != 0 {
				if isLastNode(steps, mapping, idx, i) {
					dag.SetEdge(simple.Edge{F: simple.Node(i), T: simple.Node(idx)})
				}
			} else {
				dag.SetEdge(simple.Edge{F: simple.Node(i), T: simple.Node(idx)})
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

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	yamlUtil "github.com/ghodss/yaml"
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

package util

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	cloudbuild "google.golang.org/api/cloudbuild/v1"
)

var sample001 = &cloudbuild.Build{
	Steps: []*cloudbuild.BuildStep{
		&cloudbuild.BuildStep{
			Id:   "A",
			Name: "foo",
		},
		&cloudbuild.BuildStep{
			Id:   "B",
			Name: "bar",
			WaitFor: []string{
				"-",
			},
		},
		&cloudbuild.BuildStep{
			Name: "baz",
		},
	},
}

func init() {
	// Change the current working dir to root dir
	os.Chdir("../../")
}

func TestParseYaml(t *testing.T) {
	cloudBuild, err := ParseYaml("./test/fixtures/cloudbuild.sample.001.yaml")

	assert.Empty(t, err)
	assert.EqualValues(t, sample001, cloudBuild)
}

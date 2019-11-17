package helpers

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_YAMLIntoStruct(t *testing.T) {
	assert := assert.New(t)

	// First let's write some yaml to a file
	type yamlTest struct {
		Key  string   `yaml:"key"`
		List []string `yaml:"list"`
	}

	var l []string
	l = append(l, "first item")
	l = append(l, "second item")
	l = append(l, "third item")

	expected := yamlTest{
		Key:  "Value",
		List: l,
	}

	const (
		yamlTestFile = "/tmp/test.yaml"
		yamlFileData = `
key: Value
list:
  - first item
  - second item
  - third item
`
	)
	err := ioutil.WriteFile(yamlTestFile, []byte(yamlFileData), 0644)
	assert.Nil(err)

	// Now let's get it
	var actual yamlTest

	err = YAMLIntoStruct(yamlTestFile, &actual)
	assert.Nil(err)

	assert.Equal(expected, actual)

	// Cleanup
	err = os.RemoveAll(yamlTestFile)
	assert.Nil(err)
}

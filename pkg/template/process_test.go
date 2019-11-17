package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
	yaml "gopkg.in/yaml.v3"
)

func Test_GetExcludes(t *testing.T) {
	assert := assert.New(t)

	expected := []string{"a", "b", "c"}

	testConfig := make(map[string]interface{})
	testConfig[excludedKey] = expected

	y, err := yaml.Marshal(testConfig)
	assert.Nil(err)

	err = yaml.Unmarshal(y, &testConfig)
	assert.Nil(err)

	actual := getExcludes(testConfig)

	assert.Equal(expected, actual)
}

func Test_GetEmptyExcludes(t *testing.T) {
	assert := assert.New(t)

	testConfig := make(map[string]interface{})
	testConfig["SomeKey"] = "SomeValue"

	y, err := yaml.Marshal(testConfig)
	assert.Nil(err)

	err = yaml.Unmarshal(y, &testConfig)
	assert.Nil(err)

	actual := getExcludes(testConfig)

	var expected []string

	assert.Equal(expected, actual)
}

func Test_Exclusions(t *testing.T) {
	assert := assert.New(t)

	fileList := []string{
		"some/path/service1/config.yml",
		"some/path/service1/deployment.yml",
		"some/path/service2/secret.yml",
		"some/path/service2/deployment.yml",
		"some/path/service2/service.yml",
		"another/path/service3/deployment.yml",
	}

	excludeList := []string{
		"service2",
	}

	expectedList := []string{
		"some/path/service1/config.yml",
		"some/path/service1/deployment.yml",
		"another/path/service3/deployment.yml",
	}

	var actual []string
	for _, file := range fileList {
		if isExcluded(file, &excludeList) {
			continue
		}
		actual = append(actual, file)
	}

	assert.Equal(expectedList, actual)

	excludeList = []string{
		"some",
	}

	expectedList = []string{
		"another/path/service3/deployment.yml",
	}

	actual = []string{}
	for _, file := range fileList {
		if isExcluded(file, &excludeList) {
			continue
		}
		actual = append(actual, file)
	}

	assert.Equal(expectedList, actual)
}

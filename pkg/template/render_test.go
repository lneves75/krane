package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RenderOneResource(t *testing.T) {
	assert := assert.New(t)

	testTemplate := []byte("SomeKey: {{ .Value }}")

	testConfig := make(map[string]interface{})
	testConfig["Value"] = "SomeValue"

	actual, err := render(&testTemplate, testConfig)
	assert.Nil(err)

	expected := []byte("SomeKey: SomeValue")

	assert.Len(actual, 1)
	assert.Equal(expected, *actual[0])
}

func Test_RenderMultipleResources(t *testing.T) {
	assert := assert.New(t)

	testTemplate := []byte("FirstKey: {{ .ValueA }}\n---\nSecondKey: {{ .ValueB }}")

	testConfig := make(map[string]interface{})
	testConfig["ValueA"] = "SomeValue"
	testConfig["ValueB"] = "AnotherValue"

	actual, err := render(&testTemplate, testConfig)
	assert.Nil(err)

	expectedA := []byte("FirstKey: SomeValue")
	expectedB := []byte("SecondKey: AnotherValue")

	assert.Len(actual, 2)
	assert.Equal(expectedA, *actual[0])
	assert.Equal(expectedB, *actual[1])
}

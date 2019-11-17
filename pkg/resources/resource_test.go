package resources

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetMetadataDeployment(t *testing.T) {
	assert := assert.New(t)

	body, err := ioutil.ReadFile("../../tests/samples/deployment.yml")
	assert.Nil(err)

	r := &Resource{
		body: &body,
	}

	namespace, err := r.getMetadata()
	assert.Nil(err)

	assert.Equal("Deployment", r.kind)
	assert.Equal("my-app", r.name)
	assert.Equal("", namespace)
}

func Test_GetNewNameConfigMap(t *testing.T) {
	assert := assert.New(t)

	body, err := ioutil.ReadFile("../../tests/samples/configMap.yml")
	assert.Nil(err)

	r := &Resource{
		body: &body,
	}

	namespace, err := r.getMetadata()
	assert.Nil(err)

	assert.Equal("ConfigMap", r.kind)
	assert.Equal("service1", r.name)
	assert.Equal("", namespace)

	newName, err := r.getNewName()
	assert.Nil(err)

	assert.Equal("service1-47fkf8688t", newName)
}

func Test_UpdateNameConfigMap(t *testing.T) {
	assert := assert.New(t)

	body, err := ioutil.ReadFile("../../tests/samples/configMap.yml")
	assert.Nil(err)

	r := &Resource{
		body: &body,
	}

	newName := "new-name-for-config"

	err = r.updateName(newName)
	assert.Nil(err)

	assert.Equal(newName, r.name)
}

func Test_GetNewNameSecret(t *testing.T) {
	assert := assert.New(t)

	body, err := ioutil.ReadFile("../../tests/samples/secret.yml")
	assert.Nil(err)

	r := &Resource{
		body: &body,
	}

	namespace, err := r.getMetadata()
	assert.Nil(err)

	assert.Equal("Secret", r.kind)
	assert.Equal("my-secret", r.name)
	assert.Equal("someNamespace", namespace)

	newName, err := r.getNewName()
	assert.Nil(err)

	assert.Equal("my-secret-h6k78ttkb9", newName)
}

func Test_UpdateNameSecret(t *testing.T) {
	assert := assert.New(t)

	body, err := ioutil.ReadFile("../../tests/samples/secret.yml")
	assert.Nil(err)

	r := &Resource{
		body: &body,
	}

	newName := "new-name-for-secret"

	err = r.updateName(newName)
	assert.Nil(err)

	assert.Equal(newName, r.name)
}

package resources

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_UpdateReferencesInDeployment(t *testing.T) {
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

	newNames := make(map[string]map[string]string)
	configs := make(map[string]string)
	secrets := make(map[string]string)

	configs["my-app-config-1"] = "some-other-config-1"
	configs["my-config-2"] = "some-other-config-2"

	secrets["my-app-secret-1"] = "some-other-secret-1"
	secrets["my-secret-2"] = "some-other-secret-2"

	newNames["ConfigMap"] = configs
	newNames["Secret"] = secrets

	err = r.updateManifestReferences(r.kind, newNames)
	assert.Nil(err)

	expectedBody, err := ioutil.ReadFile("../../tests/samples/updatedDeployment.yml")
	assert.Nil(err)

	assert.Equal(expectedBody, *r.body)
}

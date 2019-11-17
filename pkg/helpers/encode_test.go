package helpers

import (
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	passwordA = "somePassword"
	passwordB = "anotherPassword"
	tokenA    = "2e8da0d7-096e-4ebf-bc09-b58384294df8"
	tokenB    = "f6fb44d2-481d-42a1-acf3-b125fcff8e39"
)

func Test_Encode(t *testing.T) {
	assert := assert.New(t)

	secretsServiceA := make(map[string]interface{})
	secretsServiceA["password"] = passwordA
	secretsServiceA["token"] = tokenA

	secretsServiceB := make(map[string]interface{})
	secretsServiceB["password"] = passwordB
	secretsServiceB["token"] = tokenB

	secrets := make(map[string]interface{})
	secrets["serviceA"] = secretsServiceA
	secrets["serviceB"] = secretsServiceB

	expectedEncodedSecretsServiceA := make(map[string]string)
	expectedEncodedSecretsServiceA["password"] = base64.StdEncoding.EncodeToString([]byte(passwordA))
	expectedEncodedSecretsServiceA["token"] = base64.StdEncoding.EncodeToString([]byte(tokenA))

	expectedEncodedSecretsServiceB := make(map[string]string)
	expectedEncodedSecretsServiceB["password"] = base64.StdEncoding.EncodeToString([]byte(passwordB))
	expectedEncodedSecretsServiceB["token"] = base64.StdEncoding.EncodeToString([]byte(tokenB))

	expectedEncodedSecrets := make(map[string]interface{})
	expectedEncodedSecrets["serviceA"] = expectedEncodedSecretsServiceA
	expectedEncodedSecrets["serviceB"] = expectedEncodedSecretsServiceB

	actual := EncodeSecrets(secrets)

	assert.Equal(expectedEncodedSecrets, actual)

}

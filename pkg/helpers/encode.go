package helpers

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// EncodeSecrets will take a layered structure of maps that ends with a map of strings
// whose values will be base64 encoded
func EncodeSecrets(secrets map[string]interface{}) map[string]interface{} {
	encodedSecrets := make(map[string]interface{})

	for serviceName, serviceMap := range secrets {
		serviceSecrets := serviceMap.(map[string]interface{})

		encodedMap := make(map[string]string)
		for key, value := range serviceSecrets {
			if v, ok := value.(string); ok {
				encodedValue := base64.StdEncoding.EncodeToString([]byte(v))
				encodedMap[key] = encodedValue
			}
		}
		encodedSecrets[serviceName] = encodedMap

	}
	return encodedSecrets
}

// Mostly taken from
// https://github.com/kubernetes/kubernetes/blob/70b1c436576d8fe778e7f8cf8975ce89b8ade9f0/staging/src/k8s.io/kubectl/pkg/util/hash/hash.go#L89
// and
// https://github.com/kubernetes/kubernetes/blob/70b1c436576d8fe778e7f8cf8975ce89b8ade9f0/staging/src/k8s.io/kubectl/pkg/util/hash/hash.go#L112
func Hash(data *[]byte) (string, error) {
	hex := fmt.Sprintf("%x", sha256.Sum256(*data))

	if len(hex) < 10 {
		return "", fmt.Errorf(
			"input length must be at least 10")
	}
	enc := []rune(hex[:10])
	for i := range enc {
		switch enc[i] {
		case '0':
			enc[i] = 'g'
		case '1':
			enc[i] = 'h'
		case '3':
			enc[i] = 'k'
		case 'a':
			enc[i] = 'm'
		case 'e':
			enc[i] = 't'
		}
	}
	return string(enc), nil
}

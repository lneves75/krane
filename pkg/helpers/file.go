package helpers

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

// YAMLIntoStruct reads a YAML file into a struct passed as reference
func YAMLIntoStruct(filePath string, s interface{}) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading config: %s", err)
	}

	if err = yaml.Unmarshal(data, s); err != nil {
		return fmt.Errorf("error parsing %s: %s", filePath, err)
	}
	return nil
}

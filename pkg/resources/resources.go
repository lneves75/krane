package resources

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"

	"gopkg.in/yaml.v3"

	"github.com/lneves75/krane/pkg/helpers"
)

// Resource is a manifest body along with some of its attributes and references to it
// in the case of Secret of ConfigMap
type Resource struct {
	kind     string
	name     string
	fullPath string
	body     *[]byte
}

// Resources will contain all manifests mapped by namespace
// It will also hold a map of names for ConfigMaps and Secrets that will change due to hashing
// Those are also mapped by namespace, then kind then original name
type Resources struct {
	Namespace map[string][]*Resource
	NewNames  map[string]map[string]map[string]string
}

// Initialize is an helper function to initialize all mapping structures
func (rs *Resources) Initialize() {

	ns := make(map[string][]*Resource)
	newNames := make(map[string]map[string]map[string]string)
	k := make(map[string]map[string]string)

	k["ConfigMap"] = nil
	k["Secret"] = nil
	newNames["default"] = k

	rs.Namespace = ns
	rs.NewNames = newNames

}

// Add will receive a path to a resource and its body and add it to the list of resources in that namespace
func (rs *Resources) Add(fullPath string, body *[]byte) error {
	var err error

	r := &Resource{
		fullPath: fullPath,
		body:     body,
	}

	namespace, err := r.getMetadata()
	if err != nil {
		return err
	}

	if namespace == "" {
		namespace = "default"
	}

	if r.kind == "ConfigMap" || r.kind == "Secret" {
		newName, err := r.getNewName()
		if err != nil {
			return err
		}

		new := make(map[string]string)
		new[r.name] = newName

		rs.NewNames[namespace][r.kind] = new

		if err = r.updateName(newName); err != nil {
			return err
		}
	}

	rs.Namespace[namespace] = append(rs.Namespace[namespace], r)

	return nil
}

// UpdateReferences will go through all manifests and update the references to ConfigMaps and Secrets in all
// other manifest kinds
func (rs *Resources) UpdateReferences() error {
	for namespace, resources := range rs.Namespace {
		for _, resource := range resources {
			if err := resource.updateManifestReferences(resource.kind, rs.NewNames[namespace]); err != nil {
				return err
			}
		}
	}

	return nil
}

// WriteToDisk will write all resources to actual files
func (rs *Resources) WriteToDisk(logger *logrus.Logger) error {
	for _, resources := range rs.Namespace {
		for _, resource := range resources {
			logger.Debugf("Writing %s named %s to path %s", resource.kind, resource.name, resource.fullPath)

			filePath := filepath.Dir(resource.fullPath)
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				return fmt.Errorf("failed to create directory for file %s: %s", filePath, err)
			}

			file, err := os.Create(resource.fullPath)
			if err != nil {
				return fmt.Errorf("failed to create file %s: %s", resource.fullPath, err)
			}
			defer file.Close()

			if _, err = file.Write([]byte("# File created with krane. Do not edit\n#\n\n")); err != nil {
				return fmt.Errorf("error writing to resource file %s: %s", resource.fullPath, err)
			}

			if _, err = file.Write(*resource.body); err != nil {
				return fmt.Errorf("error writing to resource file %s: %s", resource.fullPath, err)
			}
		}
	}
	return nil
}

func (r *Resource) getMetadata() (string, error) {
	var commonMeta commonMeta

	err := yaml.Unmarshal(*r.body, &commonMeta)
	if err != nil {
		return "", err
	}

	r.kind = commonMeta.Kind
	r.name = commonMeta.Metadata.Name

	return commonMeta.Metadata.Namespace, nil
}

func (r *Resource) getNewName() (string, error) {
	hash, err := helpers.Hash(r.body)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s-%s", r.name, hash), nil
}

func (r *Resource) updateName(newName string) error {
	type manifest struct {
		Version    string            `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
		Kind       string            `json:"kind,omitempty" yaml:"kind,omitempty"`
		Metadata   metadata          `json:"metadata,omitempty" yaml:"metadata,omitempty"`
		Data       map[string]string `json:"data,omitempty" yaml:"data,omitempty"`
		StringData map[string]string `json:"stringData,omitempty" yaml:"stringData,omitempty"`
		Type       string            `json:"type,omitempty" yaml:"type,omitempty"`
	}

	var m manifest
	err := yaml.Unmarshal(*r.body, &m)
	if err != nil {
		return err
	}

	m.Metadata.Name = newName
	r.name = newName

	*r.body, err = yaml.Marshal(m)

	return err
}

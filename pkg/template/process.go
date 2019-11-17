package template

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/lneves75/krane/pkg/helpers"
	"github.com/lneves75/krane/pkg/resources"
	"github.com/sirupsen/logrus"
)

// ProcessClusters will loop through the clusters in ConfigValues
// and render the manifests for each one
func (v ConfigValues) ProcessClusters(logger *logrus.Logger, templates *Templates, l *helpers.Locations) error {
	for clusterName, clusterConfig := range v {
		logger.Debugf("Processing resources for cluster %s", clusterName)

		excludeList := getExcludes(clusterConfig.(map[string]interface{}))
		resourcesPath := filepath.Join(l.ManifestFolder, clusterName)

		var clusterResources resources.Resources
		clusterResources.Initialize()

		for file, data := range *templates {
			if isExcluded(file, &excludeList) {
				logger.Debugf("%s excluded from cluster %s", file, clusterName)
				continue
			}
			logger.Debugf("Rendering template %s", file)
			rendered, err := render(data, clusterConfig.(map[string]interface{}))
			if err != nil {
				return fmt.Errorf("failed to render template %s, %s", file, err)
			}

			for _, r := range rendered {
				err = clusterResources.Add(strings.Replace(file, l.TemplateFolder, resourcesPath, 1), r)
				if err != nil {
					return fmt.Errorf("failed to create resources from template %s, %s", file, err)
				}
			}
		}

		if err := clusterResources.UpdateReferences(); err != nil {
			return fmt.Errorf("failed to update references between manifests: %s", err)
		}

		if err := clusterResources.WriteToDisk(logger); err != nil {
			return fmt.Errorf("failed to write manifests to disk: %s", err)
		}
	}

	return nil
}

func getExcludes(clusterConfig map[string]interface{}) []string {
	var excludeList []string

	e, ok := clusterConfig[excludedKey].([]interface{})
	if !ok {
		// No exclude list
		return excludeList
	}
	for _, exclude := range e {
		excludeList = append(excludeList, exclude.(string))
	}
	return excludeList
}

func isExcluded(file string, excludeList *[]string) bool {
	for _, exclude := range *excludeList {
		if strings.Contains(file, exclude) {
			return true
		}
	}
	return false
}

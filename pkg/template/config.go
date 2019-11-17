package template

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/imdario/mergo"
	"github.com/sirupsen/logrus"

	"github.com/lneves75/krane/pkg/helpers"
)

type config struct {
	Fleets  []fleets               `yaml:"fleets"`
	Globals map[string]interface{} `yaml:"globals"`
}

type fleets struct {
	Name     string     `yaml:"name"`
	Clusters []clusters `yaml:"clusters"`
}

type clusters struct {
	Name string `yaml:"name"`
}

// ConfigValues is a map with the cluster name as a key and will contain all necessary configuration values
// to render the templates for that cluster
type ConfigValues map[string]interface{}

const excludedKey = "Exclude"

// ConfigLoad will load config files from the fleet, secrets and each of the clusters config (if any defined)
// and merge them together for the templating process
func (v *ConfigValues) ConfigLoad(logger *logrus.Logger, l *helpers.Locations) error {
	var err error
	config := &config{}
	values := make(map[string]interface{})

	if err = helpers.YAMLIntoStruct(l.ConfigFile, config); err != nil {
		logger.Fatal(err)
	}

	// For each environemnt load it's config, loop through its clusters and load any per-cluster configs
	for _, env := range config.Fleets {
		logger.Debugf("Loading config for fleet %s", env.Name)

		envConfig := make(map[string]interface{})
		if err = helpers.YAMLIntoStruct(filepath.Join(l.ConfigFolder, fmt.Sprintf("%s.yml", env.Name)), &envConfig); err != nil {
			return err
		}
		// Add globals to fleet config
		envConfig["Globals"] = config.Globals

		for _, cluster := range env.Clusters {
			e := envConfig
			logger.Debugf("Loading config for cluster %s", cluster.Name)

			// Load any cluster specific config here
			clusterConfig := make(map[string]interface{})
			if err = helpers.YAMLIntoStruct(filepath.Join(l.ClustersConfigFolder, fmt.Sprintf("%s.yml", cluster.Name)), &clusterConfig); err != nil {
				if strings.Contains(err.Error(), "no such file or directory") {
					logger.Infof("No specific configuration found for cluster %s", cluster.Name)
				} else {
					return err
				}
			}

			// Merge fleet config with the cluster config with the latter taking precendence
			if err = mergo.Merge(&clusterConfig, e); err != nil {
				return fmt.Errorf("error merging configs for cluster %s: %s", cluster, err)
			}

			if clusterConfig["Secrets"] != nil {
				logger.Debugf("Encoding secrets for cluster %s", cluster.Name)
				secrets := clusterConfig["Secrets"].(map[string]interface{})
				clusterConfig["Secrets"] = helpers.EncodeSecrets(secrets)
			}

			values[cluster.Name] = clusterConfig
			*v = values
		}
	}

	return nil
}

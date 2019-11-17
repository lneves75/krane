package helpers

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

const (
	defaultTestsFolder          = "tests"
	defaultConfigFolder         = "config"
	defaultClustersConfigFolder = "clusters"
	defaultConfigFile           = "config.yml"
	defaultTemplateFolder       = "templates"
	defaultManifestFolder       = "manifests"
)

// Locations is a collection of source locations for the
// various elements of the resource templating process
type Locations struct {
	PWD                  string // PWD will hold the current working directory
	ConfigFolder         string // Path where we can find the main configuration file
	ClustersConfigFolder string // Path where we can find specific configuration for a given cluster
	ConfigFile           string // Name of the main configuration file holidng the values
	TemplateFolder       string // Path where the template files are located
	ManifestFolder       string // Path where we should write the rendered manifests
}

// ValidateLocations will check if certain required folders exist
func (l *Locations) ValidateLocations() error {
	folders := []string{l.ConfigFolder, l.TemplateFolder}

	for _, folder := range folders {
		f := filepath.Join(folder)
		stat, err := os.Stat(f)
		if os.IsNotExist(err) || !stat.IsDir() {
			return fmt.Errorf("couldn't find folder %s: %v", f, err)
		}
	}
	return nil
}

// RemoveResourcesFolder wipes clean the folder where the rendered manifests are
func (l *Locations) RemoveResourcesFolder() error {
	if err := os.RemoveAll(l.ManifestFolder); err != nil {
		return fmt.Errorf("error removing %s's contents: %s", l.ManifestFolder, err)
	}
	return nil
}

// GetRootDir will fetch the working directory for krane from environment variable
// If not defined use the present working directory as default
func GetRootDir(logger *logrus.Logger) string {
	var workingDirectory string

	pwd, err := os.Getwd()
	if err != nil {
		logger.Fatal(err)
	}

	if rootDir, ok := os.LookupEnv("KRANE_ROOT"); ok {
		workingDirectory = rootDir
		logger.Debugf("KRANE_ROOT = %s", rootDir)
	} else {
		logger.Debugf("PWD = %s", pwd)
		workingDirectory = pwd
	}

	return workingDirectory
}

// LoadLocations will return the folder paths where the process will take place
func LoadLocations(logger *logrus.Logger, workingDirectory string) (*Locations, error) {
	var l Locations

	l.PWD = workingDirectory

	if _, ok := os.LookupEnv("CI"); ok {
		logger.Debug("Environment variable CI set, using test folders")
		l.ConfigFolder = filepath.Join(l.PWD, defaultTestsFolder, defaultConfigFolder)
		l.ClustersConfigFolder = filepath.Join(l.ConfigFolder, defaultClustersConfigFolder)
		l.ConfigFile = filepath.Join(l.ConfigFolder, defaultConfigFile)
		l.TemplateFolder = filepath.Join(l.PWD, defaultTestsFolder, defaultTemplateFolder)
		l.ManifestFolder = filepath.Join(l.PWD, defaultTestsFolder, defaultManifestFolder)
	} else {
		l.ConfigFolder = filepath.Join(l.PWD, defaultConfigFolder)
		l.ClustersConfigFolder = filepath.Join(l.ConfigFolder, defaultClustersConfigFolder)
		l.ConfigFile = filepath.Join(l.ConfigFolder, defaultConfigFile)
		l.TemplateFolder = filepath.Join(l.PWD, defaultTemplateFolder)
		l.ManifestFolder = filepath.Join(l.PWD, defaultManifestFolder)
	}
	return &l, nil
}

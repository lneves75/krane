package helpers

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lneves75/krane/pkg/logging"

	"github.com/stretchr/testify/assert"
)

const (
	kraneRoot         = "/tmp"
	defaultRootFolder = "/some/default/folder"
)

func Test_GetRootDir(t *testing.T) {
	assert := assert.New(t)

	err := os.Unsetenv("KRANE_ROOT")
	assert.Nil(err)

	rootFolder := GetRootDir(logging.TestLogger())
	assert.Equal(defaultRootFolder, defaultRootFolder)

	err = os.Setenv("KRANE_ROOT", kraneRoot)
	assert.Nil(err)

	rootFolder = GetRootDir(logging.TestLogger())
	assert.Equal(kraneRoot, rootFolder)
}

func Test_LoadLocationsCI(t *testing.T) {
	assert := assert.New(t)

	err := os.Setenv("CI", "true")
	assert.Nil(err)

	var expectedLocations Locations

	rootFolder := GetRootDir(logging.TestLogger())
	expectedLocations.PWD = rootFolder

	expectedLocations.ConfigFolder = filepath.Join(expectedLocations.PWD, defaultTestsFolder, defaultConfigFolder)
	expectedLocations.ClustersConfigFolder = filepath.Join(expectedLocations.ConfigFolder, defaultClustersConfigFolder)
	expectedLocations.ConfigFile = filepath.Join(expectedLocations.ConfigFolder, defaultConfigFile)
	expectedLocations.TemplateFolder = filepath.Join(expectedLocations.PWD, defaultTestsFolder, defaultTemplateFolder)
	expectedLocations.ManifestFolder = filepath.Join(expectedLocations.PWD, defaultTestsFolder, defaultManifestFolder)

	l, err := LoadLocations(logging.TestLogger(), rootFolder)
	assert.Nil(err)

	assert.Equal(&expectedLocations, l)
}

func Test_LoadLocationsNoCI(t *testing.T) {
	assert := assert.New(t)

	err := os.Unsetenv("CI")
	assert.Nil(err)

	var expectedLocations Locations

	rootFolder := GetRootDir(logging.TestLogger())
	expectedLocations.PWD = rootFolder

	expectedLocations.ConfigFolder = filepath.Join(expectedLocations.PWD, defaultConfigFolder)
	expectedLocations.ClustersConfigFolder = filepath.Join(expectedLocations.ConfigFolder, defaultClustersConfigFolder)
	expectedLocations.ConfigFile = filepath.Join(expectedLocations.ConfigFolder, defaultConfigFile)
	expectedLocations.TemplateFolder = filepath.Join(expectedLocations.PWD, defaultTemplateFolder)
	expectedLocations.ManifestFolder = filepath.Join(expectedLocations.PWD, defaultManifestFolder)

	l, err := LoadLocations(logging.TestLogger(), rootFolder)
	assert.Nil(err)

	assert.Equal(&expectedLocations, l)
}

func Test_ValidateFolders(t *testing.T) {
	assert := assert.New(t)

	rootFolder := GetRootDir(logging.TestLogger())
	l, err := LoadLocations(logging.TestLogger(), rootFolder)
	assert.Nil(err)

	// First clean up to make sure we get an doesn't exist error
	err = os.RemoveAll(l.ConfigFolder)
	assert.Nil(err)
	err = os.RemoveAll(l.TemplateFolder)
	assert.Nil(err)

	err = l.ValidateLocations()
	assert.NotNil(err)

	err = os.MkdirAll(l.ConfigFolder, 0755)
	assert.Nil(err)
	err = os.MkdirAll(l.TemplateFolder, 0755)
	assert.Nil(err)

	assert.Nil(l.ValidateLocations())
}

func Test_RemoveResourcesFolder(t *testing.T) {
	assert := assert.New(t)

	rootFolder := GetRootDir(logging.TestLogger())
	l, err := LoadLocations(logging.TestLogger(), rootFolder)
	assert.Nil(err)

	// Check if folder exists and if not create it
	_, err = os.Stat(l.ManifestFolder)
	if os.IsNotExist(err) {
		err = os.MkdirAll(l.ManifestFolder, 0755)
	}
	assert.Nil(err)

	err = l.RemoveResourcesFolder()
	assert.Nil(err)

	// Check if folder was removed
	_, err = os.Stat(l.ManifestFolder)
	assert.True(os.IsNotExist(err))
}

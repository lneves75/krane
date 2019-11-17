package main

import (
	"github.com/lneves75/krane/pkg/helpers"
	"github.com/lneves75/krane/pkg/logging"
	"github.com/lneves75/krane/pkg/template"
)

func main() {
	logger := logging.NewLogger()

	// Get the current working directory or use environment variable KRANE_ROOT
	workingDirectory := helpers.GetRootDir(logger)
	logger.Infof("Running krane in %s", workingDirectory)

	logger.Info("Loading location sources for templates, configurations and final manifests")
	locations, err := helpers.LoadLocations(logger, workingDirectory)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Infof("Confirming that required folders exist")
	if err = locations.ValidateLocations(); err != nil {
		logger.Fatal(err)
	}

	logger.Infof("Removing previous manifests at %s", locations.ManifestFolder)
	if err = locations.RemoveResourcesFolder(); err != nil {
		logger.Fatal(err)
	}

	logger.Debug("Loading templates file")
	var templates template.Templates
	if err = templates.ReadAll(logger, locations.TemplateFolder); err != nil {
		logger.Fatal(err)
	}

	logger.Debug("Loading configuration values for templating")
	var values template.ConfigValues
	if err = values.ConfigLoad(logger, locations); err != nil {
		logger.Fatal(err)
	}

	logger.Debug("Rendering templates for each cluster")
	if err = values.ProcessClusters(logger, &templates, locations); err != nil {
		logger.Fatal(err)
	}

	logger.Info("Done")
}

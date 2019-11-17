package template

import (
	"fmt"
	"io/ioutil"
	"regexp"

	"github.com/sirupsen/logrus"

	"github.com/karrick/godirwalk"
)

// Templates is a map with a path to a template file as key
// and the template content as value
type Templates map[string]*[]byte

// ReadAll will take a path to the template files location
// and use godirwalk function to traverse that path and load all files
func (f *Templates) ReadAll(logger *logrus.Logger, templateLocation string) error {
	templates := make(map[string]*[]byte)

	err := godirwalk.Walk(templateLocation, &godirwalk.Options{
		Callback: func(filePath string, de *godirwalk.Dirent) error {
			logger.Debugf("Current file: %s", filePath)

			// We're only interested in yml files for the templates
			var fileNamePattern = regexp.MustCompile("\\.yml$")
			if !fileNamePattern.MatchString(filePath) {
				return nil
			}

			data, err := ioutil.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("error reading template data from %s: %s", filePath, err)
			}

			templates[filePath] = &data

			return nil
		},
		Unsorted: true,
	})

	*f = templates

	return err
}

package template

import (
	"bytes"
	"fmt"
	txtTpl "text/template"

	"github.com/Masterminds/sprig"
)

// render will take the template and config data and will return a list
// of addresses to the rendered resources
func render(template *[]byte, clusterConfig map[string]interface{}) ([]*[]byte, error) {
	tpl := txtTpl.New("manifest").Funcs(sprig.TxtFuncMap())
	tplParser := txtTpl.Must(
		tpl.Parse(string(*template)),
	)

	renderedTpl := new(bytes.Buffer)
	err := tplParser.Execute(renderedTpl, clusterConfig)

	return splitResources(renderedTpl), err
}

func splitResources(renderedTpl *bytes.Buffer) []*[]byte {
	var resourceList []*[]byte

	resources := bytes.Split(renderedTpl.Bytes(), []byte(fmt.Sprintf("%s%s%s", "\n", "---", "\n")))

	for _, resource := range resources {
		if len(resource) > 0 {
			r := resource
			resourceList = append(resourceList, &r)
		}
	}
	return resourceList
}

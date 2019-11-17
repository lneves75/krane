package resources

import (
	"gopkg.in/yaml.v3"
)

func (r *Resource) updateManifestReferences(kind string, newNames map[string]map[string]string) error {
	var d map[string]interface{}

	if err := yaml.Unmarshal(*r.body, &d); err != nil {
		return err
	}

	var spec map[string]interface{}

	// Pass on the right spec to the update function
	switch kind {
	case "Deployment", "DaemonSet":
		// For deployment and daemonset is spec.template.spec
		spec = d["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})
	case "CronJob":
		// For cronjob is spec.jobTemplate.spec.template.spec
		jobTemplate := d["spec"].(map[string]interface{})["jobTemplate"].(map[string]interface{})
		spec = jobTemplate["spec"].(map[string]interface{})["template"].(map[string]interface{})["spec"].(map[string]interface{})
	default:
		return nil
	}

	if err := updateSpec(spec, newNames); err != nil {
		return err
	}

	body, err := yaml.Marshal(d)
	if err != nil {
		return err
	}

	*r.body = body

	return nil
}

func updateSpec(spec map[string]interface{}, newNames map[string]map[string]string) error {
	if _, ok := spec["initContainers"]; ok {
		if err := updateContainers(spec["initContainers"].([]interface{}), newNames); err != nil {
			return err
		}
	}
	if _, ok := spec["containers"]; ok {
		if err := updateContainers(spec["containers"].([]interface{}), newNames); err != nil {
			return err
		}
	}
	if _, ok := spec["volumes"]; ok {
		if err := updateVolumes(spec["volumes"].([]interface{}), newNames); err != nil {
			return err
		}
	}
	return nil
}

func updateContainers(containers []interface{}, newNames map[string]map[string]string) error {
	for _, c := range containers {
		container := c.(map[string]interface{})
		if _, ok := container["env"]; ok {
			if err := updateEnvVars(container["env"].([]interface{}), newNames); err != nil {
				return err
			}
		}
		if _, ok := container["envFrom"]; ok {
			if err := updateEnvFromVars(container["envFrom"].([]interface{}), newNames); err != nil {
				return err
			}
		}
	}
	return nil
}

func updateEnvVars(envVars []interface{}, newNames map[string]map[string]string) error {
	for _, e := range envVars {
		env := e.(map[string]interface{})
		if valueFrom, ok := env["valueFrom"].(map[string]interface{}); ok {
			if configMapRef, ok := valueFrom["configMapKeyRef"].(map[string]interface{}); ok {
				configMapRef["name"] = newNames["ConfigMap"][configMapRef["name"].(string)]
			}
			if secretKeyRef, ok := valueFrom["secretKeyRef"].(map[string]interface{}); ok {
				secretKeyRef["name"] = newNames["Secret"][secretKeyRef["name"].(string)]
			}
		}
	}
	return nil
}

func updateEnvFromVars(envVars []interface{}, newNames map[string]map[string]string) error {
	for _, e := range envVars {
		env := e.(map[string]interface{})
		if configMapRef, ok := env["configMapRef"].(map[string]interface{}); ok {
			configMapRef["name"] = newNames["ConfigMap"][configMapRef["name"].(string)]
		}
		if secretRef, ok := env["secretRef"].(map[string]interface{}); ok {
			secretRef["name"] = newNames["Secret"][secretRef["name"].(string)]
		}
	}
	return nil
}

func updateVolumes(volumes []interface{}, newNames map[string]map[string]string) error {
	for _, v := range volumes {
		volume := v.(map[string]interface{})
		if configMap, ok := volume["configMap"].(map[string]interface{}); ok {
			configMap["name"] = newNames["ConfigMap"][configMap["name"].(string)]
		}
		if secret, ok := volume["secret"].(map[string]interface{}); ok {
			secret["secretName"] = newNames["Secret"][secret["secretName"].(string)]
		}
	}
	return nil
}

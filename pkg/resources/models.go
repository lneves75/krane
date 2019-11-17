package resources

type commonMeta struct {
	Version  string   `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	Kind     string   `json:"kind,omitempty" yaml:"kind,omitempty"`
	Metadata metadata `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

type metadata struct {
	Name                       string            `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace                  string            `json:"namespace,omitempty" yaml:"namespace,omitempty"`
	Labels                     map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`
	Annotations                map[string]string `json:"annotations,omitempty" yaml:"annotations,omitempty"`
	ClusterName                string            `json:"clusterName,omitempty" yaml:"clusterName,omitempty"`
	CreationTimestamp          string            `json:"creationTimestamp,omitempty" yaml:"creationTimestamp,omitempty"`
	DeletionGracePeriodSeconds interface{}       `json:"deletionGracePeriodSeconds,omitempty" yaml:"deletionGracePeriodSeconds,omitempty"`
	DeletionTimestamp          string            `json:"deletionTimestamp,omitempty" yaml:"deletionTimestamp,omitempty"`
	Finalizers                 []string          `json:"finalizers,omitempty" yaml:"finalizers,omitempty"`
	GenerateName               string            `json:"generateName,omitempty" yaml:"generateName,omitempty"`
	Generation                 interface{}       `json:"generation,omitempty" yaml:"generation,omitempty"`
	Initializers               interface{}       `json:"initializers,omitempty" yaml:"initializers,omitempty"`
	OwnerReferences            interface{}       `json:"ownerReferences,omitempty" yaml:"ownerReferences,omitempty"`
	ResourceVersion            string            `json:"resourceVersion,omitempty" yaml:"resourceVersion,omitempty"`
	SelfLink                   string            `json:"selfLink,omitempty" yaml:"selfLink,omitempty"`
	UID                        string            `json:"uid,omitempty" yaml:"uid,omitempty"`
}

type kindDeployment struct {
	Spec struct {
		Template struct {
			Spec commonSpec `json:"spec,omitempty" yaml:"spec,omitempty"`
		} `json:"template,omitempty" yaml:"template,omitempty"`
	} `json:"spec,omitempty" yaml:"spec,omitempty"`
}

type commonSpec struct {
	Containers     []Container `json:"containers,omitempty" yaml:"containers,omitempty"`
	InitContainers []Container `json:"initContainers,omitempty" yaml:"initContainers,omitempty"`
	Volumes        []Volume    `json:"volume,omitempty" yaml:"volume,omitempty"`
}

type Container struct {
	Env     []Env     `json:"env,omitempty" yaml:"env,omitempty"`
	EnvFrom []EnvFrom `json:"envFrom,omitempty" yaml:"envFrom,omitempty"`
}

type Env struct {
	ValueFrom struct {
		SecretRef       Ref `json:"secretKeyRef,omitempty" yaml:"secretKeyRef,omitempty"`
		ConfigMapKeyRef Ref `json:"configMapKeyRef,omitempty" yaml:"configMapKeyRef,omitempty"`
	} `json:"valueFrom,omitempty" yaml:"valueFrom,omitempty"`
}

type EnvFrom struct {
	SecretRef       Ref `json:"secretKeyRef,omitempty" yaml:"secretKeyRef,omitempty"`
	ConfigMapKeyRef Ref `json:"configMapKeyRef,omitempty" yaml:"configMapKeyRef,omitempty"`
}

type Ref struct {
	Name string `json:"name,omitempty" yaml:"name,omitempty"`
}

type Volume struct {
	ConfigMap struct {
		Name string `json:"name,omitempty" yaml:"name,omitempty"`
	} `json:"configMap,omitempty" yaml:"configMap,omitempty"`
	Secret struct {
		Name string `json:"secretName,omitempty" yaml:"secretName,omitempty"`
	} `json:"secret,omitempty" yaml:"secret,omitempty"`
}

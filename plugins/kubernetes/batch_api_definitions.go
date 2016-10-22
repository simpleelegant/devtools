package kubernetes

// SELinuxOptions ...
type SELinuxOptions struct {
	User  string `json:"user"`
	Role  string `json:"role"`
	Type  string `json:"type"`
	Level string `json:"level"`
}

// ObjectFieldSelector ...
type ObjectFieldSelector struct {
	APIVersion string `json:"apiVersion"`
	FieldPath  string `json:"fieldPath"`
}

// VolumeMount ...
type VolumeMount struct {
	Name      string `json:"name"`
	ReadOnly  bool   `json:"readOnly"`
	MountPath string `json:"mountPath"`
}

// NFSVolumeSource ...
type NFSVolumeSource struct {
	Server   string `json:"server"`
	Path     string `json:"path"`
	ReadOnly bool   `json:"readOnly"`
}

// LabelSelector ...
type LabelSelector struct {
	MatchLabels      map[string]string          `json:"matchLabels,omitempty"`
	MatchExpressions []LabelSelectorRequirement `json:"matchExpressions,omitempty"`
}

// CephFSVolumeSource ...
type CephFSVolumeSource struct {
	Monitors  []string              `json:"monitors"`
	Path      string                `json:"path"`
	User      string                `json:"user"`
	SecretRef *LocalObjectReference `json:"secretRef"`
	ReadOnly  bool                  `json:"readOnly"`
}

// HTTPHeader ...
type HTTPHeader struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// FCVolumeSource ...
type FCVolumeSource struct {
	TargetWWNs []string `json:"targetWWNs"`
	Lun        int32    `json:"lun"`
	FsType     string   `json:"fsType"`
	ReadOnly   bool     `json:"readOnly"`
}

// DownwardAPIVolumeSource ...
type DownwardAPIVolumeSource struct {
	Items []DownwardAPIVolumeFile `json:"items"`
}

// StatusCause ...
type StatusCause struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
	Field   string `json:"field"`
}

// GCEPersistentDiskVolumeSource ...
type GCEPersistentDiskVolumeSource struct {
	PdName    string `json:"pdName"`
	FsType    string `json:"fsType"`
	Partition int32  `json:"partition"`
	ReadOnly  bool   `json:"readOnly"`
}

// TCPSocketAction ...
type TCPSocketAction struct {
	Port string `json:"port"`
}

// ConfigMapVolumeSource ...
type ConfigMapVolumeSource struct {
	Name  string      `json:"name"`
	Items []KeyToPath `json:"items"`
}

// KeyToPath ...
type KeyToPath struct {
	Key  string `json:"key"`
	Path string `json:"path"`
}

// Status often return by BadRequest response
type Status struct {
	Kind       string         `json:"kind"`
	APIVersion string         `json:"apiVersion"`
	Metadata   *ListMeta      `json:"metadata"`
	Status     string         `json:"status"`
	Message    string         `json:"message"`
	Reason     string         `json:"reason"`
	Details    *StatusDetails `json:"details"`
	Code       int32          `json:"code"`
}

// StatusDetails ...
type StatusDetails struct {
	Name              string        `json:"name"`
	Group             string        `json:"group"`
	Kind              string        `json:"kind"`
	Causes            []StatusCause `json:"causes"`
	RetryAfterSeconds int32         `json:"retryAfterSeconds"`
}

// ListMeta ...
type ListMeta struct {
	SelfLink        string `json:"selfLink"`
	ResourceVersion string `json:"resourceVersion"`
}

// GitRepoVolumeSource ...
type GitRepoVolumeSource struct {
	Repository string `json:"repository"`
	Revision   string `json:"revision"`
	Directory  int32  `json:"directory"`
}

// HTTPGetAction ...
type HTTPGetAction struct {
	Path        string       `json:"path"`
	Port        string       `json:"port"`
	Host        string       `json:"host"`
	Scheme      string       `json:"scheme"`
	HTTPHeaders []HTTPHeader `json:"httpHeaders"`
}

// Capabilities ...
type Capabilities struct {
	Add  []Capability `json:"add"`
	Drop []Capability `json:"drop"`
}

// LocalObjectReference ...
type LocalObjectReference struct {
	Name string `json:"name"`
}

// Container ...
type Container struct {
	Name                   string               `json:"name" description:"Required. Name of the container specified as a DNS_LABEL"`
	Image                  string               `json:"image,omitempty" description:"Optional. Docker image name."`
	Command                []string             `json:"command,omitempty" description:"Optional. Entrypoint array."`
	Args                   []string             `json:"args,omitempty" description:"Optional. Arguments to the entrypoint."`
	WorkingDir             string               `json:"workingDir,omitempty" description:"Optional. Container's working directory."`
	Ports                  []ContainerPort      `json:"ports,omitempty" description:"Optional. List of ports to expose from the container."`
	Env                    []EnvVar             `json:"env,omitempty" description:"Optional. List of environment variables to set in the container."`
	Resources              ResourceRequirements `json:"resources,omitempty" description:"Optional. Compute Resources required by this container."`
	VolumeMounts           []VolumeMount        `json:"volumeMounts,omitempty" description:"Optional. Pod volumes to mount into the container's filesyste."`
	LivenessProbe          *Probe               `json:"livenessProbe,omitempty" description:"Optional. Periodic probe of container liveness."`
	ReadinessProbe         *Probe               `json:"readinessProbe,omitempty" description:"Optional. Periodic probe of container service readiness."`
	Lifecycle              *Lifecycle           `json:"lifecycle,omitempty" description:"Optional. Actions that the management system should take in response to container lifecycle events."`
	TerminationMessagePath string               `json:"terminationMessagePath,omitempty" description:"Optional. Path at which the file to which the container's termination message will be written is mounted into the container’s filesystem."`
	ImagePullPolicy        string               `json:"imagePullPolicy,omitempty" description:"Optional. Image pull policy."`
	SecurityContext        *SecurityContext     `json:"securityContext,omitempty" description:"Optional. Security options the pod should run with."`
	Stdin                  bool                 `json:"stdin,omitempty" description:"Optional. Whether this container should allocate a buffer for stdin in the container runtime."`
	StdinOnce              bool                 `json:"stdinOnce,omitempty" description:"Optional. Whether the container runtime should close the stdin channel after it has been opened by a single attach."`
	TTY                    bool                 `json:"tty,omitempty" description:"Optional. Whether this container should allocate a TTY for itself, also requires stdin to be true."`
}

// PodSecurityContext ...
type PodSecurityContext struct {
	SELinuxOptions     *SELinuxOptions `json:"seLinuxOptions,omitempty"`
	RunAsUser          int64           `json:"runAsUser,omitempty"`
	RunAsNonRoot       bool            `json:"runAsNonRoot,omitempty"`
	SupplementalGroups []int32         `json:"supplementalGroups,omitempty"`
	FsGroup            int64           `json:"fsGroup,omitempty"`
}

// ExecAction ...
type ExecAction struct {
	Command []string `json:"command"`
}

// JobStatus ...
type JobStatus struct {
	Conditions     []JobCondition `json:"conditions,omitempty"`
	StartTime      string         `json:"startTime,omitempty"`
	CompletionTime string         `json:"completionTime,omitempty"`
	Active         int32          `json:"active,omitempty"`
	Succeeded      int32          `json:"succeeded,omitempty"`
	Failed         int32          `json:"failed,omitempty"`
}

// ObjectMeta ...
type ObjectMeta struct {
	Name                       string            `json:"name,omitempty" description:"Optional. Name must be unique within a namespace."`
	GenerateName               string            `json:"generateName,omitempty" description:"Optional. GenerateName is an optional prefix, used by the server, to generate a unique name ONLY IF the Name field has not been provided."`
	Namespace                  string            `json:"namespace,omitempty" description:"Optional. namespace, but 'default' is the canonical representation."`
	SelfLink                   string            `json:"selfLink,omitempty" description:"Optional. SelfLink is a URL representing this object."`
	UID                        string            `json:"uid,omitempty" description:"Optional. UID is the unique in time and space value for this object."`
	ResourceVersion            string            `json:"resourceVersion,omitempty" description:"Optional. An opaque value that represents the internal version of this object that can be used by clients to determine when objects have changed."`
	Generation                 int64             `json:"generation,omitempty" description:"Optional. A sequence number representing a specific generation of the desired state."`
	CreationTimestamp          string            `json:"creationTimestamp,omitempty" description:"Optional. CreationTimestamp is a timestamp representing the server time when this object was created."`
	DeletionTimestamp          string            `json:"deletionTimestamp,omitempty" description:"Optional. This field is set by the server when a graceful deletion is requested by the user."`
	DeletionGracePeriodSeconds int64             `json:"deletionGracePeriodSeconds,omitempty" description:"Optional. Number of seconds allowed for this object to gracefully terminate before it will be removed from the system."`
	Labels                     map[string]string `json:"labels,omitempty" description:"Optional. Map of string keys and values that can be used to organize and categorize (scope and select) objects."`
	Annotations                map[string]string `json:"annotations,omitempty" description:"Optional. Annotations is an unstructured key value map stored with a resource that may be set by external tools to store and retrieve arbitrary metadata."`
}

// HostPathVolumeSource ...
type HostPathVolumeSource struct {
	Path string `json:"path"`
}

// AzureFileVolumeSource ...
type AzureFileVolumeSource struct {
	SecretName string `json:"secretName"`
	ShareName  string `json:"ShareName"`
	ReadOnly   bool   `json:"readOnly"`
}

// ISCSIVolumeSource ...
type ISCSIVolumeSource struct {
	TargetPortal   string `json:"targetPortal"`
	Iqn            string `json:"iqn"`
	Lun            int32  `json:"lun"`
	IscsiInterface string `json:"iscsiInterface"`
	FsType         string `json:"fsType"`
	ReadOnly       bool   `json:"readOnly"`
}

// WatchEvent ...
type WatchEvent struct {
	Type   string `json:"type"`
	Object string `json:"object"`
}

// EmptyDirVolumeSource ...
type EmptyDirVolumeSource struct {
	Medium string `json:"medium"`
}

// Job represents the configuration of a single job.
type Job struct {
	Kind       string      `json:"kind,omitempty" description:"Optional. Kind is a string value representing the REST resource this object represents."`
	APIVersion string      `json:"apiVersion,omitempty" description:"Optional. APIVersion defines the versioned schema of this representation of an object."`
	Metadata   *ObjectMeta `json:"metadata,omitempty" description:"Optional. Standard object’s metadata."`
	Spec       *JobSpec    `json:"spec,omitempty" description:"Optional. Spec is a structure defining the expected behavior of a job."`
	Status     *JobStatus  `json:"status,omitempty" description:"Optional. Status is a structure describing current status of a job."`
}

// JobSpec ...
type JobSpec struct {
	Parallelism           int32            `json:"parallelism,omitempty" description:"Optional. Parallelism specifies the maximum desired number of pods the job should run at any given time."`
	Completions           int32            `json:"completions,omitempty" description:"Optional. Completions specifies the desired number of successfully finished pods the job should be run with."`
	ActiveDeadlineSeconds int64            `json:"activeDeadlineSeconds,omitempty" description:"Optional. Optional duration in seconds relative to the startTime that the job may be active before the system tries to terminate it."`
	Selector              *LabelSelector   `json:"selector,omitempty" description:"Optional. Selector is a label query over pods that should match the pod count."`
	ManualSelector        bool             `json:"ManualSelector,omitempty" description:"Optional. ManualSelector controls generation of pod labels and pod selectors."`
	Template              *PodTemplateSpec `json:"template" description:"Required. Template is the object that describes the pod that will be created when executing a job."`
}

// PodTemplateSpec ...
type PodTemplateSpec struct {
	Metadata *ObjectMeta `json:"metadata,omitempty" description:"Optional. Standard object’s metadata."`
	Spec     *PodSpec    `json:"spec,omitempty" description:"Optional. Specification of the desired behavior of the pod."`
}

// JobCondition ...
type JobCondition struct {
	Type               string `json:"type"`
	Status             string `json:"status"`
	LastProbeTime      string `json:"lastProbeTime"`
	LastTransitionTime string `json:"lastTransitionTime"`
	Reason             string `json:"reason"`
	Message            string `json:"message"`
}

// LabelSelectorRequirement ...
type LabelSelectorRequirement struct {
	Key      string   `json:"key"`
	Operator string   `json:"operator"`
	Values   []string `json:"values"`
}

// PodSpec ...
type PodSpec struct {
	Volumes                       []Volume               `json:"volumes,omitempty" description:"Optional. List of volumes that can be mounted by containers belonging to the pod."`
	Containers                    []Container            `json:"containers,omitempty" description:"Optional. List of containers belonging to the pod."`
	RestartPolicy                 string                 `json:"restartPolicy,omitempty" description:"Optional. Restart policy for all containers within the pod."`
	TerminationGracePeriodSeconds int64                  `json:"terminationGracePeriodSeconds,omitempty" description:"Optional. Optional duration in seconds the pod needs to terminate gracefully."`
	ActiveDeadlineSeconds         int64                  `json:"activeDeadlineSeconds,omitempty" description:"Optional. Optional duration in seconds the pod may be active on the node relative to StartTime before the system will actively try to mark it failed and kill associated containers."`
	DNSPolicy                     string                 `json:"dnsPolicy,omitempty" description:"Optional. Set DNS policy for containers within the pod."`
	NodeSelector                  map[string]string      `json:"nodeSelector,omitempty" description:"Optional. NodeSelector is a selector which must be true for the pod to fit on a node."`
	ServiceAccountName            string                 `json:"serviceAccountName,omitempty" description:"Optional. ServiceAccountName is the name of the ServiceAccount to use to run this pod."`
	ServiceAccount                string                 `json:"serviceAccount,omitempty" description:"Optional. ServiceAccount is a depreciated alias for ServiceAccountName."`
	NodeName                      string                 `json:"nodeName,omitempty" description:"Optional. NodeName is a request to schedule this pod onto a specific node."`
	HostNetwork                   bool                   `json:"hostNetwork,omitempty" description:"Optional. Host networking requested for this pod."`
	HostPID                       bool                   `json:"hostPID,omitempty" description:"Optional. Use the host’s pid namespace."`
	HostIPC                       bool                   `json:"hostIPC,omitempty" description:"Optional. Use the host’s ipc namespace."`
	SecurityContext               *PodSecurityContext    `json:"securityContext,omitempty" description:"Optional. SecurityContext holds pod-level security attributes and common container settings."`
	ImagePullSecrets              []LocalObjectReference `json:"imagePullSecrets,omitempty" description:"Optional. ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec."`
}

// Volume ...
type Volume struct {
	Name                  string                             `json:"name"`
	HostPath              *HostPathVolumeSource              `json:"hostPath"`
	EmptyDir              *EmptyDirVolumeSource              `json:"emptyDir"`
	GcePersistentDisk     *GCEPersistentDiskVolumeSource     `json:"gcePersistentDisk"`
	AwsElasticBlockStore  *AWSElasticBlockStoreVolumeSource  `json:"awsElasticBlockStore"`
	GitRepo               *GitRepoVolumeSource               `json:"gitRepo"`
	Secret                *SecretVolumeSource                `json:"secret"`
	Nfs                   *NFSVolumeSource                   `json:"nfs"`
	Iscsi                 *ISCSIVolumeSource                 `json:"iscsi"`
	Glusterfs             *GlusterfsVolumeSource             `json:"glusterfs"`
	PersistentVolumeClaim *PersistentVolumeClaimVolumeSource `json:"persistentVolumeClaim"`
	Rbd                   *RBDVolumeSource                   `json:"rbd"`
	FlexVolume            *FlexVolumeSource                  `json:"flexVolume"`
	Cinder                *CinderVolumeSource                `json:"cinder"`
	Cephfs                *CephFSVolumeSource                `json:"cephfs"`
	Flocker               *FlockerVolumeSource               `json:"flocker"`
	DownwardAPI           *DownwardAPIVolumeSource           `json:"downwardAPI"`
	Fc                    *FCVolumeSource                    `json:"fc"`
	AzureFile             *AzureFileVolumeSource             `json:"azureFile"`
}

// AWSElasticBlockStoreVolumeSource ...
type AWSElasticBlockStoreVolumeSource struct {
	VolumeID  string `json:"volumeID"`
	FsType    string `json:"fsType"`
	Partition int32  `json:"partition"`
	ReadOnly  bool   `json:"readOnly"`
}

// SecretVolumeSource ...
type SecretVolumeSource struct {
	SecretName string `json:"secretName"`
}

// GlusterfsVolumeSource ...
type GlusterfsVolumeSource struct {
	Endpoints string `json:"endpoints"`
	Path      string `json:"path"`
	ReadOnly  bool   `json:"readOnly"`
}

// PersistentVolumeClaimVolumeSource ...
type PersistentVolumeClaimVolumeSource struct {
	ClaimName string `json:"claimName"`
	ReadOnly  bool   `json:"readOnly"`
}

// RBDVolumeSource ...
type RBDVolumeSource struct {
	Monitors  []string              `json:"monitors"`
	Image     string                `json:"image"`
	FsType    string                `json:"fsType"`
	Pool      string                `json:"pool"`
	User      string                `json:"user"`
	Keyring   string                `json:"keyring"`
	SecretRef *LocalObjectReference `json:"secretRef"`
	ReadOnly  bool                  `json:"readOnly"`
}

// FlexVolumeSource ...
type FlexVolumeSource struct {
	Driver    string                `json:"driver"`
	FsType    string                `json:"fsType"`
	SecretRef *LocalObjectReference `json:"secretRef"`
	ReadOnly  bool                  `json:"readOnly"`
	Options   map[string]string     `json:"options"`
}

// CinderVolumeSource ...
type CinderVolumeSource struct {
	VolumeID string `json:"volumeID"`
	FsType   string `json:"fsType"`
	ReadOnly bool   `json:"readOnly"`
}

// FlockerVolumeSource ...
type FlockerVolumeSource struct {
	DatasetName string `json:"datasetName"`
}

// DownwardAPIVolumeFile ...
type DownwardAPIVolumeFile struct {
	Path     string               `json:"path"`
	FieldRef *ObjectFieldSelector `json:"fieldRef"`
}

// ContainerPort ...
type ContainerPort struct {
	Name          string `json:"name"`
	HostPort      int32  `json:"hostPort"`
	ContainerPort int32  `json:"containerPort"`
	Protocol      string `json:"protocol"`
	HostIP        string `json:"hostIP"`
}

// EnvVar ...
type EnvVar struct {
	Name      string        `json:"name" description:"Required. Name of the environment variable."`
	Value     string        `json:"value,omitempty" description:"Optional. Value of the environment variable."`
	ValueFrom *EnvVarSource `json:"valueFrom,omitempty" description:"Optional. Source for the environment variable’s value."`
}

// ResourceRequirements ...
type ResourceRequirements struct {
	Limits   map[string]string `json:"limits,omitempty"`
	Requests map[string]string `json:"requests,omitempty"`
}

// EnvVarSource ...
type EnvVarSource struct {
	FieldRef        *ObjectFieldSelector  `json:"fieldRef"`
	ConfigMapKeyRef *ConfigMapKeySelector `json:"configMapKeyRef"`
	SecretKeyRef    *SecretKeySelector    `json:"secretKeyRef"`
}

// ConfigMapKeySelector ...
type ConfigMapKeySelector struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

// SecretKeySelector ...
type SecretKeySelector struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

// Probe ...
type Probe struct {
	Exec                *ExecAction      `json:"exec"`
	HTTPGet             *HTTPGetAction   `json:"httpGet"`
	TCPSocket           *TCPSocketAction `json:"tcpSocket"`
	InitialDelaySeconds int32            `json:"initialDelaySeconds"`
	TimeoutSeconds      int32            `json:"timeoutSeconds"`
	PeriodSeconds       int32            `json:"periodSeconds"`
	SuccessThreshold    int32            `json:"successThreshold"`
	FailureThreshold    int32            `json:"failureThreshold"`
}

// Lifecycle ...
type Lifecycle struct {
	PostStart *Handler `json:"postStart"`
	PreStop   *Handler `json:"preStop"`
}

// Handler ...
type Handler struct {
	Exec      *ExecAction      `json:"exec"`
	HTTPGet   *HTTPGetAction   `json:"httpGet"`
	TCPSocket *TCPSocketAction `json:"tcpSocket"`
}

// SecurityContext ...
type SecurityContext struct {
	Capabilities           *Capabilities   `json:"capabilities,omitempty" description:"Optional. The capabilities to add/drop when running containers."`
	Privileged             bool            `json:"privileged,omitempty" description:"Optional. Run container in privileged mode."`
	SELinuxOptions         *SELinuxOptions `json:"seLinuxOptions,omitempty" description:"Optional. The SELinux context to be applied to the container."`
	RunAsUser              int64           `json:"runAsUser,omitempty" description:"Optional. The UID to run the entrypoint of the container process."`
	RunAsNonRoot           bool            `json:"runAsNonRoot,omitempty" description:"Optional. Indicates that the container must run as a non-root user."`
	ReadOnlyRootFilesystem bool            `json:"readOnlyRootFilesystem,omitempty" description:"Optional. Whether this container has a read-only root filesystem."`
}

// Capability ...
type Capability struct {
}

// JobList ...
type JobList struct {
	Kind       string    `json:"kind,omitempty"`
	APIVersion string    `json:"apiVersion,omitempty"`
	Metadata   *ListMeta `json:"metadata,omitempty"`
	Items      []*Job    `json:"items"`
}

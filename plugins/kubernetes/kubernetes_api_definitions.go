package kubernetes

// PodList is a list of Pods.
type PodList struct {
	APIVersion string            `json:"apiversion,omitempty"`
	Kind       string            `json:"kind,omitempty"`
	Metadata   *ListMeta         `json:"metadata,omitempty"`
	Items      []PodTemplateSpec `json:"items,omitempty"`
}

//Pod is a collection of containers that can run on a host.
type Pod struct {
	Kind       string      `json:"kind,omitempty"`
	APIVersion string      `json:"apiVersion,omitempty"`
	Metadata   *ObjectMeta `json:"metadata,omitempty"`
	Spec       *PodSpec    `json:"spec,omitempty"`
	Status     *PodStatus  `json:"status,omitempty"`
}

// PodStatus ...
type PodStatus struct {
	Phase             string            `json:"phase,omitempty"`
	Conditions        []PodCondition    `json:"conditions,omitempty"`
	Message           string            `json:"message,omitempty"`
	Reason            string            `json:"reason,omitempty"`
	HostIP            string            `json:"hostIP,omitempty"`
	PodIP             string            `json:"podIP,omitempty"`
	StartTime         string            `json:"startTime,omitempty"`
	ContainerStatuses []ContainerStatus `json:"containerStatuses,omitempty"`
}

// PodCondition ...
type PodCondition struct {
	Type               string `json:"type"`
	Status             string `json:"status"`
	LastProbeTime      string `json:"lastProbeTime,omitempty"`
	LastTransitionTime string `json:"lastTransitionTime,omitempty"`
	Reason             string `json:"reason,omitempty"`
	Message            string `json:"message,omitempty"`
}

// ContainerStatus contains details for the current status of this container.
type ContainerStatus struct {
	Name                 string          `json:"name"`
	State                *ContainerState `json:"state,omitempty"`
	LastTerminationState *ContainerState `json:"lastState,omitempty"`
	Ready                bool            `json:"ready"`
	RestartCount         int             `json:"restartCount"`
	Image                string          `json:"image"`
	ImageID              string          `json:"imageID"`
	ContainerID          string          `json:"containerID,omitempty"`
}

// ContainerState holds a possible state of container.
// Only one of its members may be specified.
// If none of them is specified, the default one is ContainerStateWaiting.
type ContainerState struct {
	Waiting    *ContainerStateWaiting    `json:"waiting,omitempty"`
	Running    *ContainerStateRunning    `json:"running,omitempty"`
	Terminated *ContainerStateTerminated `json:"terminated,omitempty"`
}

// ContainerStateWaiting ...
type ContainerStateWaiting struct {
	Reason  string `json:"reason,omitempty"`
	Message string `json:"message,omitempty"`
}

// ContainerStateRunning ...
type ContainerStateRunning struct {
	StartedAt string `json:"startedAt,omitempty"`
}

// ContainerStateTerminated ...
type ContainerStateTerminated struct {
	ExitCode    int    `json:"exitCode"`
	Signal      int    `json:"signal,omitempty"`
	Reason      string `json:"reason,omitempty"`
	Message     string `json:"message,omitempty"`
	StartedAt   string `json:"startedAt,omitempty"`
	FinishedAt  string `json:"finishedAt,omitempty"`
	ContainerID string `json:"containerID,omitempty"`
}

// DeleteOptions
type DeleteOptions struct {
	Kind               string `json:"kind,omitempty"`
	APIVersion         string `json:"apiVersion,omitempty"`
	GracePeriodSeconds int64  `json:"gracePeriodSeconds"`
}

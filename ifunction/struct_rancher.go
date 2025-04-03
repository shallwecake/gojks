package ifunction

import (
	"time"
)

// Cluster 表示 Rancher API 返回的集群数据（简化的结构体）
type Cluster struct {
	// 集群id
	ID string `json:"id"`
	// 命名空间
	Name string `json:"name"`
}

type ClusterResponse struct {
	Data []Cluster `json:"data"`
}

// PodList 表示整个 Pod 集合的响应
type PodList struct {
	Type         string            `json:"type"`
	Links        map[string]string `json:"links"`
	CreateTypes  map[string]string `json:"createTypes"`
	Actions      map[string]string `json:"actions"`
	ResourceType string            `json:"resourceType"`
	Revision     string            `json:"revision"`
	Count        int               `json:"count"`
	Data         []Pod             `json:"data"`
}

// Pod 表示单个 Pod 的结构
type Pod struct {
	ID         string            `json:"id"`
	Type       string            `json:"type"`
	Links      map[string]string `json:"links"`
	APIVersion string            `json:"apiVersion"`
	Kind       string            `json:"kind"`
	Metadata   Metadata          `json:"metadata"`
	Spec       Spec              `json:"spec"`
	Status     Status            `json:"status"`
}

// Metadata 表示 Pod 的元数据
type Metadata struct {
	Annotations       map[string]string `json:"annotations"`
	CreationTimestamp time.Time         `json:"creationTimestamp"`
	Fields            []string          `json:"fields"`
	GenerateName      string            `json:"generateName"`
	Labels            map[string]string `json:"labels"`
	ManagedFields     []ManagedField    `json:"managedFields"`
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	OwnerReferences   []OwnerReference  `json:"ownerReferences"`
	Relationships     []Relationship    `json:"relationships"`
	ResourceVersion   string            `json:"resourceVersion"`
	State             State             `json:"state"`
	UID               string            `json:"uid"`
}

// ManagedField 表示元数据中的管理字段
type ManagedField struct {
	APIVersion string    `json:"apiVersion"`
	FieldsType string    `json:"fieldsType"`
	FieldsV1   FieldsV1  `json:"fieldsV1"`
	Manager    string    `json:"manager"`
	Operation  string    `json:"operation"`
	Time       time.Time `json:"time"`
	// 如果有 subresource，可以添加 Subresource string `json:"subresource,omitempty"`
}

// FieldsV1 表示 ManagedField 中的字段定义（简化为 interface{}，可根据需要扩展）
type FieldsV1 map[string]interface{}

// OwnerReference 表示 Pod 的拥有者引用
type OwnerReference struct {
	APIVersion         string `json:"apiVersion"`
	BlockOwnerDeletion bool   `json:"blockOwnerDeletion"`
	Controller         bool   `json:"controller"`
	Kind               string `json:"kind"`
	Name               string `json:"name"`
	UID                string `json:"uid"`
}

// Relationship 表示 Pod 的资源关系
type Relationship struct {
	ToID    string `json:"toId"`
	ToType  string `json:"toType"`
	Rel     string `json:"rel"`
	State   string `json:"state"`
	Message string `json:"message"`
}

// State 表示 Pod 的状态
type State struct {
	Error         bool   `json:"error"`
	Message       string `json:"message"`
	Name          string `json:"name"`
	Transitioning bool   `json:"transitioning"`
}

// Spec 表示 Pod 的规格
type Spec struct {
	Containers                    []Container       `json:"containers"`
	DNSPolicy                     string            `json:"dnsPolicy"`
	EnableServiceLinks            bool              `json:"enableServiceLinks"`
	ImagePullSecrets              []ImagePullSecret `json:"imagePullSecrets"`
	NodeName                      string            `json:"nodeName"`
	PreemptionPolicy              string            `json:"preemptionPolicy"`
	Priority                      int               `json:"priority"`
	RestartPolicy                 string            `json:"restartPolicy"`
	SchedulerName                 string            `json:"schedulerName"`
	SecurityContext               interface{}       `json:"securityContext"` // 可根据需要定义具体结构
	ServiceAccount                string            `json:"serviceAccount"`
	ServiceAccountName            string            `json:"serviceAccountName"`
	TerminationGracePeriodSeconds int               `json:"terminationGracePeriodSeconds"`
	Tolerations                   []Toleration      `json:"tolerations"`
	Volumes                       []Volume          `json:"volumes"`
}

// Container 表示 Pod 中的容器
type Container struct {
	Args                     []string      `json:"args"`
	Command                  []string      `json:"command"`
	Image                    string        `json:"image"`
	ImagePullPolicy          string        `json:"imagePullPolicy"`
	Name                     string        `json:"name"`
	Resources                interface{}   `json:"resources"` // 可根据需要定义具体结构
	TerminationMessagePath   string        `json:"terminationMessagePath"`
	TerminationMessagePolicy string        `json:"terminationMessagePolicy"`
	VolumeMounts             []VolumeMount `json:"volumeMounts"`
}

// VolumeMount 表示容器的挂载卷
type VolumeMount struct {
	MountPath string `json:"mountPath"`
	Name      string `json:"name"`
	ReadOnly  bool   `json:"readOnly,omitempty"`
}

// ImagePullSecret 表示镜像拉取密钥
type ImagePullSecret struct {
	Name string `json:"name"`
}

// Toleration 表示 Pod 的容忍规则
type Toleration struct {
	Effect            string `json:"effect"`
	Key               string `json:"key"`
	Operator          string `json:"operator"`
	Value             string `json:"value,omitempty"`
	TolerationSeconds int    `json:"tolerationSeconds,omitempty"`
}

// Volume 表示 Pod 的卷
type Volume struct {
	Name      string           `json:"name"`
	ConfigMap *ConfigMapVolume `json:"configMap,omitempty"`
	Projected *ProjectedVolume `json:"projected,omitempty"`
}

// ConfigMapVolume 表示 ConfigMap 类型的卷
type ConfigMapVolume struct {
	DefaultMode int    `json:"defaultMode"`
	Name        string `json:"name"`
}

// ProjectedVolume 表示 Projected 类型的卷（简化为 interface{}，可扩展）
type ProjectedVolume struct {
	DefaultMode int           `json:"defaultMode"`
	Sources     []interface{} `json:"sources"` // 可根据实际数据定义具体结构
}

// Status 表示 Pod 的状态
type Status struct {
	Conditions        []Condition       `json:"conditions"`
	ContainerStatuses []ContainerStatus `json:"containerStatuses"`
	HostIP            string            `json:"hostIP"`
	Phase             string            `json:"phase"`
	PodIP             string            `json:"podIP"`
	PodIPs            []PodIP           `json:"podIPs"`
	QOSClass          string            `json:"qosClass"`
	StartTime         string            `json:"startTime"` // 可改为 time.Time，根据格式调整
}

// Condition 表示 Pod 的状态条件
type Condition struct {
	Error              bool       `json:"error"`
	LastProbeTime      *time.Time `json:"lastProbeTime,omitempty"`
	LastTransitionTime time.Time  `json:"lastTransitionTime"`
	LastUpdateTime     time.Time  `json:"lastUpdateTime"`
	Status             string     `json:"status"`
	Transitioning      bool       `json:"transitioning"`
	Type               string     `json:"type"`
}

// ContainerStatus 表示容器状态
type ContainerStatus struct {
	ContainerID  string         `json:"containerID"`
	Image        string         `json:"image"`
	ImageID      string         `json:"imageID"`
	LastState    interface{}    `json:"lastState"` // 可根据需要定义具体结构
	Name         string         `json:"name"`
	Ready        bool           `json:"ready"`
	RestartCount int            `json:"restartCount"`
	Started      bool           `json:"started"`
	State        ContainerState `json:"state"`
}

// ContainerState 表示容器当前状态
type ContainerState struct {
	Running struct {
		StartedAt string `json:"startedAt"`
	} `json:"running"`
}

// PodIP 表示 Pod 的 IP 地址
type PodIP struct {
	IP string `json:"ip"`
}

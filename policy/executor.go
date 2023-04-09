package policy

import (
	"github.com/capitalonline/cds-eci-profile/resource"
	"github.com/capitalonline/cds-eci-profile/utils"
	v1 "k8s.io/api/core/v1"
)

type PatchInfo struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

type Executor interface {
	OnPodCreating(selector *resource.Selector, pod *v1.Pod) ([]PatchInfo, error)
	OnPodUnscheduled(selector *resource.Selector, pod *v1.Pod) (*utils.PatchOption, error)
}

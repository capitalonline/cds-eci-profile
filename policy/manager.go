package policy

import (
	"github.com/capitalonline/cds-eci-profile/resource"
	"github.com/capitalonline/cds-eci-profile/utils"
	v1 "k8s.io/api/core/v1"
)

const (
	ExecutorNameVirtualNodeOnly = "VirtualNodeOnly"
)

type Manager struct {
	executors map[string]Executor
}

func NewManager() *Manager {
	return &Manager{
		executors: map[string]Executor{
			ExecutorNameVirtualNodeOnly: NewVirtualNodeOnlyExecutor(),
		},
	}
}

func (m *Manager) OnPodCreating(selector *resource.Selector, pod *v1.Pod) ([]PatchInfo, error) {
	executor := m.findExecutor()
	return executor.OnPodCreating(selector, pod)
}

func (m *Manager) OnPodUnscheduled(selector *resource.Selector, pod *v1.Pod) (*utils.PatchOption, error) {
	executor := m.findExecutor()
	return executor.OnPodUnscheduled(selector, pod)
}

func (m *Manager) findExecutor() Executor {
	executorName := ExecutorNameVirtualNodeOnly
	return m.executors[executorName]
}

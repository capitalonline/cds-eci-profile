package policy

import (
	"encoding/json"
	"github.com/capitalonline/cds-eci-profile/resource"
	v1 "k8s.io/api/core/v1"
	"os"
)

const (
	vnodeNodeSelectorKey = "virtual-kubelet.io/provider"
	vnodeNodeSelectorVal = "cds-provider"
)

type VKTaint struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Effect string `json:"effect"`
}

func init() {
	vkTaintStr := os.Getenv("VK_TAINTS")
	if vkTaintStr == "" {
		virtualNodeTolerationList = append(virtualNodeTolerationList, v1.Toleration{
			Key:      vnodeNodeSelectorKey,
			Value:    vnodeNodeSelectorVal,
			Operator: v1.TolerationOpEqual,
			Effect:   v1.TaintEffectNoSchedule,
		})
	} else {
		var l []VKTaint
		err := json.Unmarshal([]byte(vkTaintStr), &l)
		if err != nil {
			virtualNodeTolerationList = append(virtualNodeTolerationList, v1.Toleration{
				Key:      vnodeNodeSelectorKey,
				Value:    vnodeNodeSelectorVal,
				Operator: v1.TolerationOpEqual,
				Effect:   v1.TaintEffectNoSchedule,
			})
		} else {
			for _, t := range l {
				var effect v1.TaintEffect
				switch t.Effect {
				case "NoSchedule":
					effect = v1.TaintEffectNoSchedule
				case "NoExecute":
					effect = v1.TaintEffectNoExecute
				case "PreferNoSchedule":
					effect = v1.TaintEffectPreferNoSchedule
				default:
					continue
				}
				virtualNodeTolerationList = append(virtualNodeTolerationList, v1.Toleration{
					Key:      t.Key,
					Value:    t.Value,
					Operator: v1.TolerationOpEqual,
					Effect:   effect,
				})
			}
		}
	}
}

var (
	virtualNodeTolerationList []v1.Toleration
	//virtualNodeToleration     = v1.Toleration{
	//	Key:      vnodeNodeSelectorKey,
	//	Value:    vnodeNodeSelectorVal,
	//	Operator: v1.TolerationOpEqual,
	//	Effect:   v1.TaintEffectNoSchedule,
	//}
)

func addVirtualNodeToleration(pod *v1.Pod) PatchInfo {
	tolerations := pod.Spec.Tolerations
	tolerations = append(tolerations, virtualNodeTolerationList...)
	return PatchInfo{
		Op:    "add",
		Path:  "/spec/tolerations",
		Value: tolerations,
	}
}

func addVirtualNodeSelector() PatchInfo {
	return PatchInfo{
		Op:   "replace",
		Path: "/spec/nodeSelector",
		Value: map[string]string{
			"type": "virtual-kubelet",
		},
	}
}

func addAnnotations(selector *resource.Selector, pod *v1.Pod) PatchInfo {
	annotations := pod.Annotations
	if annotations == nil {
		annotations = map[string]string{
			"effect-match-selector": selector.Name,
		}
	}
	for key, value := range selector.Spec.Effect.Annotations {
		annotations[key] = value
	}
	return PatchInfo{
		Op:    "add",
		Path:  "/metadata/annotations",
		Value: annotations,
	}
}

func addLabels(selector *resource.Selector, pod *v1.Pod) PatchInfo {
	labels := pod.Labels
	if labels == nil {
		labels = map[string]string{}
	}

	for key, value := range selector.Spec.Effect.Labels {
		labels[key] = value
	}

	return PatchInfo{
		Op:    "add",
		Path:  "/metadata/labels",
		Value: labels,
	}
}

func existVirtualTolerations(tolerations []v1.Toleration) bool {
	for _, toleration := range tolerations {
		if toleration.Key == vnodeNodeSelectorKey &&
			toleration.Value == vnodeNodeSelectorVal &&
			toleration.Operator == v1.TolerationOpEqual &&
			toleration.Effect == v1.TaintEffectNoSchedule {
			return true
		}
	}

	return false
}

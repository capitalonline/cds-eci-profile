package resource

import (
	"encoding/json"
	"github.com/capitalonline/cds-eci-profile/utils"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	SELECTORS         = "selectors"
	ObjectSelector    = "objectSelector"
	NamespaceSelector = "namespaceSelector"
	MatchLabels       = "matchLabels"
	Effect            = "effect"
	Annotations       = "annotations"
	Labels            = "labels"
)

type Selector struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SelectorSpec `json:"spec"`
}

type SelectorSpec struct {
	NamespaceLabels *metav1.LabelSelector `json:"namespaceLabels,omitempty"`
	ObjectLabels    *metav1.LabelSelector `json:"objectLabels,omitempty"`
	Effect          *SideEffect           `json:"effect,omitempty"`
}

type SideEffect struct {
	Annotations map[string]string `json:"annotations,omitempty"` // 需要追加的annotation
	Labels      map[string]string `json:"labels,omitempty"`      // 需要追加的label
}

func MakeSelectors(configMap map[string]string) (selectors []*Selector) {
	v, have := configMap[SELECTORS]
	if have {
		var list []map[string]interface{}
		_ = json.Unmarshal([]byte(v), &list)
		for _, selector := range list {
			objectSelector := new(metav1.LabelSelector)
			namespaceSelector := new(metav1.LabelSelector)
			effect := new(SideEffect)
			if data, ok := selector[ObjectSelector]; ok {
				matchLabels := utils.Map(data)
				if matchLabels != nil {
					objectSelector.MatchLabels = utils.MapStrStr(matchLabels[MatchLabels])
				}
			}
			if data, ok := selector[NamespaceSelector]; ok {
				matchLabels := utils.Map(data)
				if matchLabels != nil {
					namespaceSelector.MatchLabels = utils.MapStrStr(matchLabels[MatchLabels])
				}
			}
			if data, ok := selector[Effect]; ok {
				effectContent := utils.Map(data)
				if effectContent != nil {
					effect.Annotations = utils.MapStrStr(effectContent[Annotations])
					effect.Labels = utils.MapStrStr(effectContent[Labels])
				}
			}
			selectors = append(selectors, &Selector{
				ObjectMeta: metav1.ObjectMeta{Name: utils.String(selector["name"])},
				Spec: SelectorSpec{
					NamespaceLabels: namespaceSelector,
					ObjectLabels:    objectSelector,
					Effect:          effect,
				},
			})
		}
	}
	return
}

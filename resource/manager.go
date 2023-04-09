package resource

import (
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	listercorev1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

type Manager struct {
	coreV1InformerFactory informers.SharedInformerFactory
	podInformer           cache.SharedIndexInformer
	podLister             listercorev1.PodLister
	nsInformer            cache.SharedIndexInformer
	nsLister              listercorev1.NamespaceLister
	selectorInformer      cache.SharedIndexInformer
	selectorLister        listercorev1.ConfigMapLister
}

func NewManager(k8sClient *kubernetes.Clientset) *Manager {
	coreV1InformerFactory := informers.NewSharedInformerFactory(k8sClient, 30*time.Second)
	return &Manager{
		coreV1InformerFactory: coreV1InformerFactory,
		podInformer:           coreV1InformerFactory.Core().V1().Pods().Informer(),
		podLister:             coreV1InformerFactory.Core().V1().Pods().Lister(),
		nsInformer:            coreV1InformerFactory.Core().V1().Namespaces().Informer(),
		nsLister:              coreV1InformerFactory.Core().V1().Namespaces().Lister(),
		selectorInformer:      coreV1InformerFactory.Core().V1().ConfigMaps().Informer(),
		selectorLister:        coreV1InformerFactory.Core().V1().ConfigMaps().Lister(),
	}
}

func (m *Manager) Run(stopChan <-chan struct{}) {
	go m.coreV1InformerFactory.Start(stopChan)
}

func (m *Manager) HasSynced() bool {
	return m.podInformer.HasSynced() &&
		m.nsInformer.HasSynced() &&
		m.selectorInformer.HasSynced()
}

func (m *Manager) AddPodEventHandler(handler cache.ResourceEventHandler) {
	_, _ = m.podInformer.AddEventHandler(handler)
}

func (m *Manager) GetNamespace(name string) (*v1.Namespace, error) {
	return m.nsLister.Get(name)
}

func (m *Manager) ListSelectors(namespace, name string) ([]*Selector, error) {
	configMap, err := m.GetSelector(namespace, name)
	if err != nil {
		return nil, err
	}
	return MakeSelectors(configMap.Data), nil
}

func (m *Manager) GetSelector(namespace, name string) (*v1.ConfigMap, error) {
	return m.selectorLister.ConfigMaps(namespace).Get(name)
}

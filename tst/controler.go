package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

// NamespaceController следит через Kubernetes API за изменениями
// в пространствах имен и создает RoleBinding для конкретного namespace.
type NamespaceController struct {
	namespaceInformer cache.SharedIndexInformer
	kclient           *kubernetes.Clientset
}

// NewNamespaceController создает новый NewNamespaceController
func NewNamespaceController(kclient *kubernetes.Clientset) *NamespaceController {
	namespaceWatcher := &NamespaceController{}

	// Создаем информер для слежения за Namespaces
	namespaceInformer := cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				return kclient.Core().Namespaces().List(options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				return kclient.Core().Namespaces().Watch(options)
			},
		},
		&v1.Namespace{},
		3*time.Minute,
		cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
	)

	namespaceInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: namespaceWatcher.createRoleBinding,
	})

	namespaceWatcher.kclient = kclient
	namespaceWatcher.namespaceInformer = namespaceInformer

	return namespaceWatcher
}

func (c *NamespaceController) createRoleBinding(obj interface{}) {
	namespaceObj := obj.(*v1.Namespace)
	namespaceName := namespaceObj.Name

	roleBinding := &v1beta1.RoleBinding{
		TypeMeta: metav1.TypeMeta{
			Kind:       "RoleBinding",
			APIVersion: "rbac.authorization.k8s.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("ad-kubernetes-%s", namespaceName),
			Namespace: namespaceName,
		},
		Subjects: []v1beta1.Subject{
			v1beta1.Subject{
				Kind: "Group",
				Name: fmt.Sprintf("ad-kubernetes-%s", namespaceName),
			},
		},
		RoleRef: v1beta1.RoleRef{
			APIGroup: "rbac.authorization.k8s.io",
			Kind:     "ClusterRole",
			Name:     "edit",
		},
	}

	_, err := c.kclient.Rbac().RoleBindings(namespaceName).Create(roleBinding)

	if err != nil {
		log.Println(fmt.Sprintf("Failed to create Role Binding: %s", err.Error()))
	} else {
		log.Println(fmt.Sprintf("Created AD RoleBinding for Namespace: %s", roleBinding.Name))
	}
}

// Run запускает процесс ожидания изменений в пространствах имён
// и действия в соответствии с этими изменениями.
func (c *NamespaceController) Run(stopCh <-chan struct{}, wg *sync.WaitGroup) {
	// Когда эта функция завершена, пометим как выполненную
	defer wg.Done()

	// Инкрементируем wait group, т.к. собираемся вызвать goroutine
	wg.Add(1)

	// Вызываем goroutine
	go c.namespaceInformer.Run(stopCh)

	// Ожидаем получения стоп-сигнала
	<-stopCh
}

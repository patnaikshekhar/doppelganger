package k8s

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

type K8sClient struct {
	client *kubernetes.Clientset
}

func NewClient() (*K8sClient, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	config, err := clientcmd.BuildConfigFromFlags("", path.Join(home, ".kube/config"))
	if err != nil {
		return nil, err
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &K8sClient{
		client: client,
	}, nil

}

func (c *K8sClient) NewInformerForServices() error {
	factory := informers.NewSharedInformerFactory(c.client, 5*time.Second)
	svcInformer := factory.Core().V1().Services().Informer()
	svcInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			s := obj.(*corev1.Service)
			fmt.Printf("Informer event: Service ADDED %s/%s\n", s.GetNamespace(), s.GetName())
		},
		UpdateFunc: func(old, new interface{}) {
			s := old.(*corev1.Service)
			fmt.Printf("Informer event: Service UPDATED %s/%s\n", s.GetNamespace(), s.GetName())
		},
		DeleteFunc: func(obj interface{}) {
			s := obj.(*corev1.Service)
			fmt.Printf("Informer event: Service DELETED %s/%s\n", s.GetNamespace(), s.GetName())
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	factory.Start(ctx.Done())

	for informerType, ok := range factory.WaitForCacheSync(ctx.Done()) {
		if !ok {
			return fmt.Errorf("could not sync informer %s", informerType)
		}
	}

	return nil

}

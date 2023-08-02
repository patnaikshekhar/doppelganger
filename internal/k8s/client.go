package k8s

import (
	"context"
	"doppelganger/internal/services"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/utils/strings/slices"
)

type K8sClient struct {
	client *kubernetes.Clientset
	config *rest.Config
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
		config: config,
	}, nil

}

func (c *K8sClient) NewInformerForServices(
	ctx context.Context,
	all bool, namespaces []string,
	minlocalPort *uint32,
	fwdServices *services.ForwardedServices,
	events chan services.ForwardedService) error {
	factory := informers.NewSharedInformerFactory(c.client, 5*time.Second)
	svcInformer := factory.Core().V1().Services().Informer()
	svcInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			s := obj.(*corev1.Service)
			if all || slices.Contains(namespaces, s.Namespace) {
				fmt.Printf("Informer event: Service ADDED: %s/%s\n", s.GetNamespace(), s.GetName())
				servicePort := []int{}
				for _, p := range s.Spec.Ports {
					log.Printf("Here 1 port is %d", p.Port)
					servicePort = append(servicePort, int(p.Port))
				}
				err := c.PortForward(ctx, s.Name, s.Namespace, servicePort, minlocalPort, s.Spec.Selector, fwdServices, events)
				if err != nil {
					log.Printf("Error in port forward %s", err)
				}
			}
		},
		DeleteFunc: func(obj interface{}) {
			s := obj.(*corev1.Service)
			if all || slices.Contains(namespaces, s.Namespace) {
				fmt.Printf("Informer event: Service DELETED: %s/%s\n", s.GetNamespace(), s.GetName())
			}
		},
	})

	factory.Start(ctx.Done())

	for informerType, ok := range factory.WaitForCacheSync(ctx.Done()) {
		if !ok {
			return fmt.Errorf("could not sync informer %s", informerType)
		}
	}

	return nil

}

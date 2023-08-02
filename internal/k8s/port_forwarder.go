package k8s

import (
	"context"
	"doppelganger/internal/services"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync/atomic"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

func (c *K8sClient) PortForward(
	ctx context.Context,
	svcName,
	namespace string,
	ports []int,
	localPort *uint32,
	svcLabels map[string]string,
	fwdServices *services.ForwardedServices,
	events chan services.ForwardedService) error {

	// Get Pods with the same labels as the service
	pods, err := c.client.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labels.SelectorFromSet(svcLabels).String(),
	})
	if err != nil {
		return err
	}

	if len(pods.Items) == 0 {
		return fmt.Errorf("the service %s does not have any pods", svcName)
	}

	pod := pods.Items[0]

	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward",
		namespace, pod.Name)
	hostIP := strings.TrimLeft(c.config.Host, "https://")

	transport, upgrader, err := spdy.RoundTripperFor(c.config)
	if err != nil {
		return err
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, http.MethodPost, &url.URL{Scheme: "https", Path: path, Host: hostIP})

	for _, servicePort := range ports {
		readyCh := make(chan struct{})
		lp := atomic.AddUint32(localPort, 1)
		fw, err := portforward.New(dialer, []string{fmt.Sprintf("%d:%d", lp, servicePort)}, ctx.Done(), readyCh, os.Stdout, os.Stderr)
		if err != nil {
			return err
		}

		// a go-routine helps prevent blocking calls
		go fw.ForwardPorts()

		go func(readyCh chan struct{}, svcName string, namespace string, lp uint32, servicePort uint32) {
			<-readyCh
			fmt.Printf("Service %s in namespace %s forwarded to port %d\n", svcName, namespace, lp)
			fsvc := services.ForwardedService{
				Name:        svcName,
				Namespace:   namespace,
				LocalPort:   lp,
				ServicePort: servicePort,
			}

			fwdServices.Lock()
			defer fwdServices.Unlock()

			fwdServices.Services = append(fwdServices.Services, fsvc)

			events <- fsvc

		}(readyCh, svcName, namespace, lp, uint32(servicePort))
	}

	return nil
}

/*
127.0.0.1 nginx.dg-test
127.0.0.1 service.default
*/

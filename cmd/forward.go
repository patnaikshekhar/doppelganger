package cmd

import (
	"context"
	"doppelganger/internal/dns"
	"doppelganger/internal/http_server"
	"doppelganger/internal/k8s"
	"doppelganger/internal/proxy"
	"doppelganger/internal/services"
	"log"

	"github.com/spf13/cobra"
)

var (
	all          bool
	minLocalPort uint32
	namespaces   []string
)

var forwardCmd = &cobra.Command{
	Use:     "forward",
	Aliases: []string{"fwd"},
	Short:   "Setup port forward for k8s services",
	Long:    `Setup port forward for k8s services`,
	RunE: func(cmd *cobra.Command, args []string) error {
		waitCtx, cancel := context.WithCancel(context.Background())
		defer cancel()

		creationEventChannel := make(chan services.ForwardedService)

		msp := proxy.NewMultiPortProxy(creationEventChannel)
		go msp.Start()

		client, err := k8s.NewClient()
		if err != nil {
			return err
		}

		fwdServices := services.ForwardedServices{}
		// minPort := atomic.AddInt32()
		err = client.NewInformerForServices(waitCtx, all, namespaces, &minLocalPort, &fwdServices, creationEventChannel)
		if err != nil {
			return err
		}

		dnsProvider := dns.NewHosts()

		go func() {
			for event := range creationEventChannel {
				err := dnsProvider.Add(event)
				if err != nil {
					log.Printf("Error occured %s", err)
				}
			}
		}()

		server := http_server.New(&fwdServices)
		err = server.Start()
		if err != nil {
			return err
		}

		// Handling k8s informer and keeping it alive till a signal to kill is made available
		<-waitCtx.Done()
		return nil
	},
}

func init() {
	forwardCmd.Flags().BoolVar(&all, "all", false, "Defaults to true, set to false if you would like to select services to forward")
	forwardCmd.Flags().Uint32Var(&minLocalPort, "min-port", 30000, "")
	forwardCmd.Flags().StringSliceVar(&namespaces, "namespaces", []string{}, "Namespaces to filter by")
}

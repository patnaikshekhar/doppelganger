package cmd

import (
	"doppelganger/internal/k8s"

	"github.com/spf13/cobra"
)

var all bool

var forwardCmd = &cobra.Command{
	Use:     "forward",
	Aliases: []string{"fwd"},
	Short:   "Setup port forward for k8s services",
	Long:    `Setup port forward for k8s services`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := k8s.NewClient()
		if err != nil {
			return err
		}

		err = client.NewInformerForServices()
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	forwardCmd.Flags().BoolVar(&all, "all", true, "Defaults to true, set to false if you would like to select services to forward")
}

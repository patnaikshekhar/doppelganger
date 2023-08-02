package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "doppleganger",
	Short: "Doppleganger is a CLI to make services running in k8s clusters accessible locally",
	Long: `Doppleganger is a CLI to get services running in 
			your Kubernetes clusters seamlessly closer 
			to you on your local machine`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(forwardCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(killCmd)
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "List port forwarded k8s services",
	Long:  `List port forwarded k8s services`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Listig port forwarded services")
	},
}

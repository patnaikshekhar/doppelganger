package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var killCmd = &cobra.Command{
	Use:   "kill",
	Short: "Ends all active processes handling port forwards",
	Long:  `Ends all active processes handling port forwards`,
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := http.Post("http://localhost:1234/kill", "text/plain", nil)
		if err != nil {
			return err
		}

		if resp.StatusCode > 399 {
			return fmt.Errorf("encountered an error while killing processes %d", resp.StatusCode)
		}

		return nil
	},
}

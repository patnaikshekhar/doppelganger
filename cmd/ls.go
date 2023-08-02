package cmd

import (
	"doppelganger/internal/services"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "List port forwarded k8s services",
	Long:  `List port forwarded k8s services`,
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := http.Get("http://localhost:1234/list")
		if err != nil {
			return err
		}

		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		currentListOfForwardedServices := services.ForwardedServices{}
		err = json.Unmarshal(respBytes, &currentListOfForwardedServices)
		if err != nil {
			return err
		}

		headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
		columnFmt := color.New(color.FgYellow).SprintfFunc()

		tbl := table.New("Service Name", "Namespace", "Local Port", "URL")
		tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

		for _, svc := range currentListOfForwardedServices.Services {
			url := fmt.Sprintf("%s.%s:%d", svc.Name, svc.Namespace, svc.LocalPort)
			tbl.AddRow(svc.Name, svc.Namespace, svc.LocalPort, url)
		}

		tbl.Print()

		return nil
	},
}

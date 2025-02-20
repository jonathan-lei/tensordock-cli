package commands

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

var (
	serversCmd = &cobra.Command{
		Use:   "servers",
		Short: "Manage servers",
	}
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List servers",
		RunE:  serverList,
	}
	infoCmd = &cobra.Command{
		Use:   "info [flags] server_id",
		Short: "Get server info",
		Args:  cobra.ExactArgs(1),
		RunE:  serverInfo,
	}
	startCmd = &cobra.Command{
		Use:     "start [flags] server_id",
		Short:   "Start a server",
		Args:    cobra.ExactArgs(1),
		RunE:    startServer,
		PostRun: logAction("success"),
	}
	stopCmd = &cobra.Command{
		Use:     "stop [flags] server_id",
		Short:   "Stop a server",
		Args:    cobra.ExactArgs(1),
		RunE:    stopServer,
		PostRun: logAction("success"),
	}
	deleteCmd = &cobra.Command{
		Use:     "delete [flags] server_id",
		Short:   "Delete a server",
		Args:    cobra.ExactArgs(1),
		RunE:    deleteServer,
		PostRun: logAction("success"),
	}
	deployCmd = &cobra.Command{
		Use:     "deploy [flags] name admin_user admin_pass",
		Short:   "Deploy a server",
		Args:    cobra.ExactArgs(3),
		RunE:    deployServer,
		PostRun: logAction("success"),
	}
	manageCmd = &cobra.Command{
		Use:   "manage [flags] server_id",
		Short: "Open server management panel in a browser",
		Args:  cobra.ExactArgs(1),
		RunE:  manageServer,
	}
)

func init() {
	serversCmd.AddCommand(listCmd)

	serversCmd.AddCommand(infoCmd)

	serversCmd.AddCommand(stopCmd)

	serversCmd.AddCommand(startCmd)

	serversCmd.AddCommand(deleteCmd)

	serversCmd.AddCommand(deployCmd)
	deployCmd.Flags().String("gpuModel", "A40", "The GPU model that you would like to provision")
	deployCmd.Flags().String("location", "na-us-las-1", "Location")
	deployCmd.Flags().String("instanceType", "gpu", "Either \"gpu\" or \"cpu\"")
	deployCmd.Flags().Int("gpuCount", 1, "The number of GPUs of the model you specified earlier")
	deployCmd.Flags().Int("vcpus", 1, "Number of vCPUs that you would like")
	deployCmd.Flags().Int("storage", 20, "Number of GB of networked storage")
	deployCmd.Flags().String("storageClass", "st1", "io1 or st1, depending on storage class desired")
	deployCmd.Flags().Int("ram", 2, "Number of GB of RAM to be deployed.")
	deployCmd.Flags().String("os", "Ubuntu 18.04 LTS", "Operating system")

	serversCmd.AddCommand(manageCmd)

	rootCmd.AddCommand(serversCmd)
}

func serverList(cmd *cobra.Command, args []string) error {
	res, err := client.ListServers()
	if err != nil {
		return err
	}

	if !res.Success {
		return errors.New("endpoint returned error")
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Id", "Name", "Location", "Status"})
	for _, elem := range res.Servers {
		t.AppendRow(table.Row{elem.Id, elem.Name, elem.Location, elem.Status})
	}
	t.Render()

	return nil
}

func serverInfo(cmd *cobra.Command, args []string) error {
	server := args[0]
	res, err := client.GetServer(server)
	if err != nil {
		return err
	}

	if !res.Success {
		return errors.New("endpoint returned error")
	}

	props := []map[string]string{
		{"name": "ID", "value": res.Server.Id},
		{"name": "Name", "value": res.Server.Name},
		{"name": "Location", "value": res.Server.Location},
		{"name": "IP", "value": res.Server.Ip},
		{"name": "Charged Cost", "value": fmt.Sprintf("%v", res.Server.Cost.Charged)},
		{"name": "Hour-On Cost", "value": fmt.Sprintf("%v", res.Server.Cost.HourOn)},
		{"name": "Minutes-On Cost", "value": fmt.Sprintf("%v", res.Server.Cost.MinutesOn)},
		{"name": "Hour-Off Cost", "value": fmt.Sprintf("%v", res.Server.Cost.HourOff)},
		{"name": "Minutes-Off Cost", "value": fmt.Sprintf("%v", res.Server.Cost.MinutesOff)},
		{"name": "CPU Model", "value": res.Server.CPUModel},
		{"name": "GPU Count", "value": strconv.Itoa(res.Server.GPUCount)},
		{"name": "GPU Model", "value": res.Server.GPUModel},
		{"name": "RAM", "value": fmt.Sprintf("%vGB", res.Server.Ram)},
		{"name": "Status", "value": res.Server.Status},
		{"name": "Storage", "value": fmt.Sprintf("%vGB", res.Server.Storage)},
		{"name": "Storage Class", "value": res.Server.StorageClass},
		{"name": "Type", "value": res.Server.Type},
		{"name": "vCPUs", "value": strconv.Itoa(res.Server.VCPUs)},
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Property", "Value"})
	for _, elem := range props {
		t.AppendRow(table.Row{elem["name"], elem["value"]})
	}
	t.Render()

	return nil
}

func startServer(cmd *cobra.Command, args []string) error {
	server := args[0]
	res, err := client.StartServer(server)
	if err != nil {
		return err
	}

	if !res.Success {
		return err
	}

	return nil
}

func stopServer(cmd *cobra.Command, args []string) error {
	server := args[0]
	res, err := client.StopServer(server)
	if err != nil {
		return err
	}

	if !res.Success {
		return err
	}

	return nil
}

func deleteServer(cmd *cobra.Command, args []string) error {
	server := args[0]
	res, err := client.DeleteServer(server)
	if err != nil {
		return err
	}

	if !res.Success {
		return err
	}

	return nil
}

func deployServer(cmd *cobra.Command, args []string) error {
	flags := cmd.Flags()

	instanceType, err := flags.GetString("instanceType")
	if err != nil {
		return err
	}

	gpuModel, err := flags.GetString("gpuModel")
	if err != nil {
		return err
	}

	gpuCount, err := flags.GetInt("gpuCount")
	if err != nil {
		return err
	}

	vcpus, err := flags.GetInt("vcpus")
	if err != nil {
		return err
	}

	ram, err := flags.GetInt("ram")
	if err != nil {
		return err
	}

	storage, err := flags.GetInt("storage")
	if err != nil {
		return err
	}

	storageClass, err := flags.GetString("storageClass")
	if err != nil {
		return err
	}

	os, err := flags.GetString("os")
	if err != nil {
		return err
	}

	location, err := flags.GetString("location")
	if err != nil {
		return err
	}

	name := args[0]
	adminUser := args[1]
	adminPass := args[2]

	res, err := client.DeployServer(
		adminUser,
		adminPass,
		instanceType,
		gpuModel,
		gpuCount,
		vcpus,
		ram,
		storage,
		storageClass,
		os,
		location,
		name,
	)

	if err != nil {
		return err
	}

	if !res.Success {
		return err
	}

	fmt.Println(res.Server.Id)

	return nil
}

func manageServer(cmd *cobra.Command, args []string) error {
	server := args[0]
	res, err := client.GetServer(server)
	if err != nil {
		return err
	}

	if !res.Success {
		return err
	}

	err = browser.OpenURL(res.Server.Links["dashboard"]["href"])
	if err != nil {
		return err
	}

	return nil
}

func logAction(message string) func(*cobra.Command, []string) {
	return func(c *cobra.Command, s []string) { log.Println(message) }
}

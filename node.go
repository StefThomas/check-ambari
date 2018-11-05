package main

import (
	"fmt"
	"github.com/disaster37/go-ambari-rest/client"
	"github.com/disaster37/go-nagios"
	"gopkg.in/urfave/cli.v1"
	"strings"
)

// Perform a node check
func checkNode(c *cli.Context) error {

	monitoringData := nagiosPlugin.NewMonitoring()

	// Check global parameters
	err := manageGlobalParameters()
	if err != nil {
		return cli.NewExitError(fmt.Sprintf("%v", err), nagiosPlugin.STATUS_UNKNOWN)
	}

	// Check current parameters
	if c.String("cluster-name") == "" {
		return cli.NewExitError("You must set --cluster-name parameter", nagiosPlugin.STATUS_UNKNOWN)
	}
	if c.String("node-name") == "" {
		return cli.NewExitError("You must set --node-name parameter", nagiosPlugin.STATUS_UNKNOWN)
	}

	params := &OptionnalComputeState{}
	if c.String("include-alerts") != "" {
		params.IncludeAlerts = strings.Split(c.String("include-alerts"), ",")
	}
	if c.String("exclude-alerts") != "" {
		params.ExcludeAlerts = strings.Split(c.String("exclude-alerts"), ",")
	}

	// Get Ambari connection
	ambariClient := client.New(ambariUrl+"/api/v1", ambariLogin, ambariPassword)
	ambariClient.DisableVerifySSL()

	// Check node alertes
	alerts, err := ambariClient.AlertsInHost(c.String("cluster-name"), c.String("node-name"))
	if err != nil {
		monitoringData.AddMessage("Somethink wrong when try to check node alerts on %s: %v", c.String("node-name"), err)
		monitoringData.SetStatus(nagiosPlugin.STATUS_UNKNOWN)
		monitoringData.ToSdtOut()
	}

	// Remove all UNKNOWN alert before to compute them
	filterAlerts := make([]client.Alert, 0, 1)
	for _, alert := range alerts {
		if alert.AlertInfo.State != "UNKNOWN" {
			filterAlerts = append(filterAlerts, alert)
		}
	}

	monitoringData, err = computeState(filterAlerts, monitoringData, params)
	if err != nil {
		return err
	}

	monitoringData.ToSdtOut()
	return nil

}

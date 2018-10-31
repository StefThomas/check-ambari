package main

import (
	"fmt"
	"github.com/disaster37/go-ambari-rest/client"
	"github.com/disaster37/go-nagios"
	"gopkg.in/urfave/cli.v1"
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

	if len(alerts) > 0 {
		monitoringData.AddPerfdata("nbAlert", len(alerts), "")
	} else {
		monitoringData.AddMessage("All works fine !")
	}

	for _, alert := range alerts {
		if alert.AlertInfo.ServiceName != "" && alert.AlertInfo.ComponentName != "" {
			monitoringData.AddMessage("%s/%s - %s", alert.AlertInfo.ServiceName, alert.AlertInfo.ComponentName, alert.AlertInfo.Text)
		} else if alert.AlertInfo.ServiceName != "" {
			monitoringData.AddMessage("%s - %s", alert.AlertInfo.ServiceName, alert.AlertInfo.Text)
		} else {
			monitoringData.AddMessage(alert.AlertInfo.Text)
		}

		monitoringData.SetStatusAsString(alert.AlertInfo.State)
	}

	monitoringData.ToSdtOut()
	return nil

}

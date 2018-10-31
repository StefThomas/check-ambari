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

	nbAlert := 0
	for _, alert := range alerts {

		// Skip Unkown alert
		if alert.AlertInfo.State != "UNKNOWN" {

			if nbAlert == 0 {
				monitoringData.AddMessage("There are some problems !")
			}
			nbAlert++

			if alert.AlertInfo.ServiceName != "" && alert.AlertInfo.ComponentName != "" {
				monitoringData.AddMessage("%s - %s/%s - %s", alert.AlertInfo.State, alert.AlertInfo.ServiceName, alert.AlertInfo.ComponentName, alert.AlertInfo.Label)
			} else if alert.AlertInfo.ServiceName != "" {
				monitoringData.AddMessage("%s - %s - %s", alert.AlertInfo.State, alert.AlertInfo.ServiceName, alert.AlertInfo.Label)
			} else {
				monitoringData.AddMessage("%s - %s", alert.AlertInfo.State, alert.AlertInfo.Label)
			}

			monitoringData.SetStatusAsString(alert.AlertInfo.State)
		}
	}

	monitoringData.AddPerfdata("nbAlert", nbAlert, "")
	if nbAlert == 0 {
		monitoringData.AddMessage("All works fine !")
	}

	monitoringData.ToSdtOut()
	return nil

}

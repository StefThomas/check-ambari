package main

import (
	"fmt"
	"github.com/disaster37/go-ambari-rest/client"
	"github.com/disaster37/go-nagios"
	"gopkg.in/urfave/cli.v1"
	"strings"
)

// Perform a service check
func checkService(c *cli.Context) error {

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
	if c.String("service-name") == "" {
		return cli.NewExitError("You must set --service-name parameter", nagiosPlugin.STATUS_UNKNOWN)
	}

	// Get Ambari connection
	ambariClient := client.New(ambariUrl+"/api/v1", ambariLogin, ambariPassword)
	ambariClient.DisableVerifySSL()

	// Check service alertes
	alerts, err := ambariClient.AlertsInService(c.String("cluster-name"), strings.ToUpper(c.String("service-name")))
	if err != nil {
		monitoringData.AddMessage("Somethink wrong when try to check service alerts on %s: %v", c.String("service-name"), err)
		monitoringData.SetStatus(nagiosPlugin.STATUS_UNKNOWN)
		monitoringData.ToSdtOut()
	}

	nbAlert := 0
	for _, alert := range alerts {

		// Skip Unkown alert
		if alert.AlertInfo.State != "UNKNOWN" {

			// Keep only service alerte
			if c.Bool("exclude-node-alerts") && alert.AlertInfo.Scope == "SERVICE" {
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

			} else if c.Bool("exclude-node-alerts") == false {

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
	}

	monitoringData.AddPerfdata("nbAlert", nbAlert, "")
	if nbAlert == 0 {
		monitoringData.AddMessage("All works fine !")
	}

	monitoringData.ToSdtOut()
	return nil

}

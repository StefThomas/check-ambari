package main

import (
	"github.com/disaster37/go-ambari-rest/client"
	"github.com/disaster37/go-nagios"
)

func computeState(alerts []client.Alert, monitoringData *nagiosPlugin.Monitoring, scopes []string) *nagiosPlugin.Monitoring {

	//Keep only alert in desired scope
	var filterAlerts []client.Alert
	if len(scopes) == 0 {
		filterAlerts = alerts
	} else {

		filterAlerts = make([]client.Alert, 0, 0)
		for _, alert := range alerts {
			for _, scope := range scopes {
				if alert.AlertInfo.Scope == scope {
					filterAlerts = append(filterAlerts, alert)
					break
				}
			}
		}
	}

	// Compute the state
	nbAlert := len(filterAlerts)
	monitoringData.AddPerfdata("nbAlert", nbAlert, "")
	if nbAlert == 0 {
		monitoringData.AddMessage("All works fine !")
	} else {
		monitoringData.AddMessage("There are some problems !")
		for _, alert := range filterAlerts {

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

	return monitoringData

}

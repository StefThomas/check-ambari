package main

import (
	"encoding/json"
	"github.com/disaster37/go-ambari-rest/client"
	"github.com/disaster37/go-nagios"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type OptionnalComputeState struct {
	Scopes        []string
	IncludeAlerts []string
	ExcludeAlerts []string
}

func (o *OptionnalComputeState) String() string {
	json, _ := json.Marshal(o)
	return string(json)
}

func computeState(alerts []client.Alert, monitoringData *nagiosPlugin.Monitoring, params *OptionnalComputeState) (*nagiosPlugin.Monitoring, error) {

	log.Debugf("OptionnalComputeState: %s", params)

	if len(params.ExcludeAlerts) > 0 && len(params.IncludeAlerts) > 0 {
		return nil, errors.New("You need to use ExcludeAlerts or IncludeAlerts!")
	}

	//Keep only alert in desired scope
	var filterAlerts []client.Alert

	filterAlerts = make([]client.Alert, 0, 0)
	for _, alert := range alerts {
		log.Debugf("Start to filter alert %s", alert.AlertInfo.Label)
		if len(params.Scopes) > 0 {
			for _, scope := range params.Scopes {
				if alert.AlertInfo.Scope == scope {
					// Check to exclude Alert
					if len(params.ExcludeAlerts) > 0 {
						isExclude := false
						for _, excludeAlert := range params.ExcludeAlerts {
							if excludeAlert == alert.AlertInfo.Label {
								isExclude = true
								log.Debugf("Alert %s is exclude", alert.AlertInfo.Label)
								break
							}
						}
						if isExclude == false {
							filterAlerts = append(filterAlerts, alert)
							log.Debugf("Alert %s is not exclude", alert.AlertInfo.Label)

						}
					} else if len(params.IncludeAlerts) > 0 {
						// Check to include alerts
						for _, includeAlert := range params.IncludeAlerts {
							if includeAlert == alert.AlertInfo.Label {
								filterAlerts = append(filterAlerts, alert)
								log.Debugf("Alert %s is include", alert.AlertInfo.Label)
								break
							}
						}
					} else {
						filterAlerts = append(filterAlerts, alert)
						log.Debugf("Alert %s is not exclude/include", alert.AlertInfo.Label)
					}
				} else {
					log.Debugf("Alert %s in not the require scope", alert.AlertInfo.Label)
				}
			}
		} else {
			if len(params.ExcludeAlerts) > 0 {
				isExclude := false
				for _, excludeAlert := range params.ExcludeAlerts {
					if excludeAlert == alert.AlertInfo.Label {
						isExclude = true
						log.Debugf("Alert %s is exclude", alert.AlertInfo.Label)
						break
					}
				}
				if isExclude == false {
					filterAlerts = append(filterAlerts, alert)
					log.Debugf("Alert %s is not exclude", alert.AlertInfo.Label)

				}
			} else if len(params.IncludeAlerts) > 0 {
				// Check to include alerts
				for _, includeAlert := range params.IncludeAlerts {
					if includeAlert == alert.AlertInfo.Label {
						filterAlerts = append(filterAlerts, alert)
						log.Debugf("Alert %s is include", alert.AlertInfo.Label)
						break
					}
				}
			} else {
				filterAlerts = append(filterAlerts, alert)
				log.Debugf("Alert %s is not include/exclude", alert.AlertInfo.Label)
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
			log.Debugf("Process alert %s", alert.AlertInfo.Label)

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

	return monitoringData, nil

}

package main

import (
	"github.com/disaster37/go-ambari-rest/client"
	"github.com/disaster37/go-nagios"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestComputeState(t *testing.T) {

	// when no alerts
	monitoringData := nagiosPlugin.NewMonitoring()
	monitoringData = computeState(make([]client.Alert, 0, 0), monitoringData, []string{})
	assert.Equal(t, 0, monitoringData.Status())
	assert.Equal(t, 1, len(monitoringData.Perfdatas()))
	assert.Equal(t, 1, len(monitoringData.Messages()))

	// When one alerte
	monitoringData = nagiosPlugin.NewMonitoring()
	alert1 := client.Alert{
		AlertInfo: &client.AlertInfo{
			State:         "WARNING",
			Label:         "label",
			ServiceName:   "service",
			ComponentName: "component",
			Scope:         "SERVICE",
		},
	}
	listAlert := make([]client.Alert, 0, 2)
	listAlert = append(listAlert, alert1)
	monitoringData = computeState(listAlert, monitoringData, []string{})
	assert.Equal(t, 1, monitoringData.Status())
	assert.Equal(t, 1, len(monitoringData.Perfdatas()))
	assert.Equal(t, 2, len(monitoringData.Messages()))

	// When 2 alert with diffrent state
	monitoringData = nagiosPlugin.NewMonitoring()
	alert2 := client.Alert{
		AlertInfo: &client.AlertInfo{
			State:         "CRITICAL",
			Label:         "label2",
			ServiceName:   "service2",
			ComponentName: "component2",
			Scope:         "HOST",
		},
	}
	listAlert = append(listAlert, alert2)
	monitoringData = computeState(listAlert, monitoringData, []string{})
	assert.Equal(t, 2, monitoringData.Status())
	assert.Equal(t, 1, len(monitoringData.Perfdatas()))
	assert.Equal(t, 3, len(monitoringData.Messages()))

	// When include only specific scope
	monitoringData = nagiosPlugin.NewMonitoring()
	monitoringData = computeState(listAlert, monitoringData, []string{"SERVICE"})
	assert.Equal(t, 1, monitoringData.Status())
	assert.Equal(t, 1, len(monitoringData.Perfdatas()))
	assert.Equal(t, 2, len(monitoringData.Messages()))
}

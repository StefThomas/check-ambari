package main

import (
	"github.com/disaster37/go-ambari-rest/client"
	"github.com/disaster37/go-nagios"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestComputeState(t *testing.T) {

	logrus.SetLevel(logrus.DebugLevel)

	// when no alerts
	params := &OptionnalComputeState{}
	monitoringData := nagiosPlugin.NewMonitoring()
	params.Scopes = nil
	params.IncludeAlerts = nil
	params.ExcludeAlerts = nil
	monitoringData, err := computeState(make([]client.Alert, 0, 0), monitoringData, params)
	assert.NoError(t, err)
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
	params.Scopes = nil
	params.IncludeAlerts = nil
	params.ExcludeAlerts = nil
	monitoringData, err = computeState(listAlert, monitoringData, params)
	assert.NoError(t, err)
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
	params.Scopes = nil
	params.IncludeAlerts = nil
	params.ExcludeAlerts = nil
	monitoringData, err = computeState(listAlert, monitoringData, params)
	assert.NoError(t, err)
	assert.Equal(t, 2, monitoringData.Status())
	assert.Equal(t, 1, len(monitoringData.Perfdatas()))
	assert.Equal(t, 3, len(monitoringData.Messages()))

	// When include only specific scope
	params.Scopes = []string{"SERVICE"}
	params.IncludeAlerts = nil
	params.ExcludeAlerts = nil
	monitoringData = nagiosPlugin.NewMonitoring()
	monitoringData, err = computeState(listAlert, monitoringData, params)
	assert.NoError(t, err)
	assert.Equal(t, 1, monitoringData.Status())
	assert.Equal(t, 1, len(monitoringData.Perfdatas()))
	assert.Equal(t, 2, len(monitoringData.Messages()))

	// When we use exclude list
	params.Scopes = nil
	params.IncludeAlerts = nil
	params.ExcludeAlerts = []string{"label2"}
	monitoringData = nagiosPlugin.NewMonitoring()
	monitoringData, err = computeState(listAlert, monitoringData, params)
	assert.NoError(t, err)
	assert.Equal(t, 1, monitoringData.Status())
	assert.Equal(t, 1, len(monitoringData.Perfdatas()))
	assert.Equal(t, 2, len(monitoringData.Messages()))

	// When we use include list
	params.Scopes = nil
	params.ExcludeAlerts = nil
	params.IncludeAlerts = []string{"label"}
	monitoringData = nagiosPlugin.NewMonitoring()
	monitoringData, err = computeState(listAlert, monitoringData, params)
	assert.NoError(t, err)
	assert.Equal(t, 1, monitoringData.Status())
	assert.Equal(t, 1, len(monitoringData.Perfdatas()))
	assert.Equal(t, 2, len(monitoringData.Messages()))

	// When we use exclude list with scope
	params.Scopes = []string{"SERVICE"}
	params.IncludeAlerts = nil
	params.ExcludeAlerts = []string{"label2"}
	monitoringData = nagiosPlugin.NewMonitoring()
	monitoringData, err = computeState(listAlert, monitoringData, params)
	assert.NoError(t, err)
	assert.Equal(t, 1, monitoringData.Status())
	assert.Equal(t, 1, len(monitoringData.Perfdatas()))
	assert.Equal(t, 2, len(monitoringData.Messages()))

	// When we use include list with scope
	params.Scopes = []string{"SERVICE"}
	params.ExcludeAlerts = nil
	params.IncludeAlerts = []string{"label"}
	monitoringData = nagiosPlugin.NewMonitoring()
	monitoringData, err = computeState(listAlert, monitoringData, params)
	assert.NoError(t, err)
	assert.Equal(t, 1, monitoringData.Status())
	assert.Equal(t, 1, len(monitoringData.Perfdatas()))
	assert.Equal(t, 2, len(monitoringData.Messages()))

	// When we use include and exlude list in same time
	params.Scopes = nil
	params.ExcludeAlerts = []string{"label"}
	params.IncludeAlerts = []string{"label"}
	monitoringData = nagiosPlugin.NewMonitoring()
	monitoringData, err = computeState(listAlert, monitoringData, params)
	assert.Error(t, err)
}

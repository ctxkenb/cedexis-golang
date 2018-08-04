package main

import (
	"strconv"

	"github.com/ctxkenb/cedexis-golang/cedexis"
)

var alerts *map[int]*cedexis.Alert

func getAlerts() ([]*cedexis.Alert, error) {
	if alerts == nil {
		newAlerts, err := cClient.GetAlerts()
		if err != nil {
			return nil, err
		}

		alerts = &map[int]*cedexis.Alert{}
		for _, v := range newAlerts {
			(*alerts)[*v.ID] = v
		}
	}

	m := make([]*cedexis.Alert, 0, len(*alerts))
	for _, val := range *alerts {
		m = append(m, val)
	}
	return m, nil
}

func getAlert(name string) (*cedexis.Alert, error) {
	allAlerts, err := getAlerts()
	if err != nil {
		return nil, err
	}

	for _, a := range allAlerts {
		if *(a.Name) == name {
			return a, nil
		}
	}

	return nil, nil
}

func alertsToTable(alerts []*cedexis.Alert) *Table {
	t := Table{
		Columns: []string{"Name", "Enabled", "Platform"},
		Rows:    make([][]string, len(alerts)),
	}

	for i, a := range alerts {
		var pName string
		p, err := cClient.GetPrivatePlatform(*a.Platform)
		if err != nil || p.Name == nil {
			pName = "<err: lookup>"
		} else {
			pName = *p.Name
		}

		t.Rows[i] = []string{*a.Name, strconv.FormatBool(*a.Enabled), pName}
	}

	return &t
}

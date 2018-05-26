package main

import (
	"strconv"

	"github.com/ctxkenb/cedexis-golang/cedexis"
)

func getAlerts() ([]*cedexis.Alert, error) {
	return cClient.GetAlerts()
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

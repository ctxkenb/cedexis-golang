package main

import (
	"fmt"
	"strconv"

	"github.com/ctxkenb/cedexis-golang/cedexis"
)

var apps *map[int]*cedexis.Application

func getApps() ([]*cedexis.Application, error) {
	if apps == nil {
		newApps, err := cClient.GetApplications()
		if err != nil {
			return nil, err
		}

		apps = &map[int]*cedexis.Application{}
		for _, v := range newApps {
			(*apps)[*v.ID] = v
		}
	}

	m := make([]*cedexis.Application, 0, len(*apps))
	for _, val := range *apps {
		m = append(m, val)
	}
	return m, nil
}

func getApp(name string) (*cedexis.Application, error) {
	fmt.Printf("getApp: %s\n", name)
	allApps, err := getApps()
	if err != nil {
		return nil, err
	}

	for _, a := range allApps {
		if *(a.Name) == name {
			return a, nil
		}
	}

	return nil, nil
}

func appsToTable(apps []*cedexis.Application) *Table {
	t := Table{
		Columns: []string{"Name", "Enabled", "CNAME"},
		Rows:    make([][]string, len(apps)),
	}

	for i, a := range apps {
		t.Rows[i] = []string{*a.Name, strconv.FormatBool(*a.Enabled), *a.Cname}
	}

	return &t
}

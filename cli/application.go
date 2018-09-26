package main

import (
	"regexp"
	"strconv"

	"github.com/ctxkenb/cedexis-golang/cedexis"
)

var apps *map[int]*cedexis.Application

func createApp(name string, description string, t string, fallbackName string, availabilityThreshold int, targetPlatform int, targetCname string, useSonar bool) error {
	apps = nil

	newApp := cedexis.NewApplication(name, description, t, fallbackName, availabilityThreshold, []cedexis.ApplicationPlatform{
		{
			ID:           &targetPlatform,
			Cname:        &targetCname,
			SonarEnabled: &useSonar,
		},
	})

	_, err := cClient.CreateApplication(newApp)
	return err
}

func filterApps(apps []*cedexis.Application, filter string) ([]*cedexis.Application, error) {
	if filter == "" {
		return apps, nil
	}

	re, err := regexp.Compile(filter)
	if err != nil {
		return nil, err
	}

	result := make([]*cedexis.Application, 0, len(apps))
	for _, p := range apps {
		if re.MatchString(*p.Name) {
			result = append(result, p)
		}
	}
	return result, nil
}

func deleteApp(name string) error {
	app, err := getApp(name)
	if err != nil {
		return err
	}

	apps = nil

	return cClient.DeleteApplication(*app.ID)
}

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

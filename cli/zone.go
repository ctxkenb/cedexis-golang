package main

import (
	"regexp"
	"strconv"

	"github.com/ctxkenb/cedexis-golang/cedexis"
)

var zones *map[int]*cedexis.Zone

func createZone(name string, description string, tags []string, content *string) error {
	zones = nil
	_, err := cClient.CreateZone(name, description, tags, content)
	return err
}

func filterZones(zones []*cedexis.Zone, filter string) ([]*cedexis.Zone, error) {
	if filter == "" {
		return zones, nil
	}

	re, err := regexp.Compile(filter)
	if err != nil {
		return nil, err
	}

	result := make([]*cedexis.Zone, 0, len(zones))
	for _, p := range zones {
		if re.MatchString(*p.DomainName) {
			result = append(result, p)
		}
	}
	return result, nil
}

func getZones() ([]*cedexis.Zone, error) {
	if zones == nil {
		newZones, err := cClient.GetZones()
		if err != nil {
			return nil, err
		}

		zones = &map[int]*cedexis.Zone{}
		for _, v := range newZones {
			(*zones)[*v.ID] = v
		}
	}

	m := make([]*cedexis.Zone, 0, len(*zones))
	for _, val := range *zones {
		m = append(m, val)
	}
	return m, nil
}

func getZone(name string) (*cedexis.Zone, error) {
	allZones, err := getZones()
	if err != nil {
		return nil, err
	}

	for _, z := range allZones {
		if *(z.DomainName) == name {
			return z, nil
		}
	}

	return nil, nil
}

func deleteZone(name string) error {
	z, err := getZone(name)
	if err != nil {
		return err
	}

	zones = nil
	return cClient.DeleteZone(*z.ID)
}

func zonesToTable(zones []*cedexis.Zone) *Table {
	t := Table{
		Columns: []string{"Domain", "Records", "Description"},
		Rows:    make([][]string, len(zones)),
	}

	for i, a := range zones {
		t.Rows[i] = []string{*a.DomainName, strconv.FormatInt(int64(len(a.Records)), 10), *a.Description}
	}

	return &t
}

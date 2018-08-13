package main

import (
	"strconv"

	"github.com/ctxkenb/cedexis-golang/cedexis"
)

var zones *map[int]cedexis.Zone

func getZones() ([]cedexis.Zone, error) {
	if zones == nil {
		newZones, err := cClient.GetZones()
		if err != nil {
			return nil, err
		}

		zones = &map[int]cedexis.Zone{}
		for _, v := range newZones {
			(*zones)[*v.ID] = v
		}
	}

	m := make([]cedexis.Zone, 0, len(*zones))
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
			return &z, nil
		}
	}

	return nil, nil
}

func zonesToTable(zones []cedexis.Zone) *Table {
	t := Table{
		Columns: []string{"Domain", "Records", "Description"},
		Rows:    make([][]string, len(zones)),
	}

	for i, a := range zones {
		t.Rows[i] = []string{*a.DomainName, strconv.FormatInt(int64(len(a.Records)), 10), *a.Description}
	}

	return &t
}

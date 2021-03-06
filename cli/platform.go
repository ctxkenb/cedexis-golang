package main

import (
	"fmt"
	"regexp"

	"github.com/ctxkenb/cedexis-golang/cedexis"
)

func resetPlatformCache() {
	platforms = map[cedexis.PlatformType][]*cedexis.PlatformInfo{}
}

var platforms = map[cedexis.PlatformType][]*cedexis.PlatformInfo{}

func getPlatforms(t cedexis.PlatformType, category *cedexis.PlatformCategory) []*cedexis.PlatformInfo {
	if platforms[t] == nil {
		var err error
		platforms[t], err = cClient.GetPlatforms(t)
		if err != nil {
			return []*cedexis.PlatformInfo{}
		}
	}

	if category == nil {
		return platforms[t]
	}

	result := make([]*cedexis.PlatformInfo, 0, len(platforms[t]))
	for _, p := range platforms[t] {
		if *p.Category.ID == *category {
			result = append(result, p)
		}
	}

	return result
}

func filterPlatforms(platforms []*cedexis.PlatformInfo, filter string) ([]*cedexis.PlatformInfo, error) {
	if filter == "" {
		return platforms, nil
	}

	re, err := regexp.Compile(filter)
	if err != nil {
		return nil, err
	}

	result := make([]*cedexis.PlatformInfo, 0, len(platforms))
	for _, p := range platforms {
		if re.MatchString(*p.Name) {
			result = append(result, p)
		}
	}
	return result, nil
}

func platformsToTable(platforms []*cedexis.PlatformInfo) *Table {
	t := Table{
		Columns: []string{"Name", "ID", "Category", "Alias"},
		Rows:    make([][]string, len(platforms)),
	}

	for i, p := range platforms {
		archeType := ""
		if p.AliasedPlatform != nil && p.AliasedPlatform.Name != nil {
			archeType = *p.AliasedPlatform.Name
		}
		t.Rows[i] = []string{*p.Name, fmt.Sprintf("%d", *p.ID), *p.Category.Name, archeType}
	}

	return &t
}

func getPlatformID(name string, t cedexis.PlatformType, c *cedexis.PlatformCategory) (int, error) {
	p, err := getPlatform(name, t, c)
	if err != nil {
		return 0, err
	}

	return *p.ID, nil
}

func getPlatform(name string, t cedexis.PlatformType, c *cedexis.PlatformCategory) (*cedexis.PlatformInfo, error) {
	availPlatforms := getPlatforms(t, c)

	for _, p := range availPlatforms {
		if *p.Name == name {
			return p, nil
		}
	}

	return nil, fmt.Errorf("platform '%v' not found", name)
}

func createPlatform(name string, shortName string, description string, archeType int, tags []string, sonar *cedexis.SonarConfig) error {
	p := cedexis.NewPublicCloudPrivatePlatform(shortName, name, description, archeType, tags)
	p.SonarConfig = sonar

	var err error
	p, err = cClient.CreatePrivatePlatform(p)
	if err != nil {
		return err
	}

	resetPlatformCache()
	return nil
}

func deletePlatform(name string, t cedexis.PlatformType) error {
	platformID, err := getPlatformID(name, t, nil)
	if err != nil {
		return err
	}

	err = cClient.DeletePrivatePlatform(platformID)
	if err != nil {
		return err
	}

	resetPlatformCache()
	return nil
}

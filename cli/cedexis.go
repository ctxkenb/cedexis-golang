package main

import (
	"github.com/ctxkenb/cedexis-golang/cedexis"
)

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
		if *p.Category.ID == int(*category) {
			result = append(result, p)
		}
	}

	return result
}

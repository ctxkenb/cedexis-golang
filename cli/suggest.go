package main

import (
	"github.com/ctxkenb/cedexis-golang/cedexis"
	"github.com/ctxkenb/cedexis-golang/cli/parser"
)

func suggestAllPlatforms(s string) []parser.Suggestion {
	return suggestPlatforms(s, cedexis.PlatformsTypeAll, nil)
}

func suggestCommunityPlatforms(s string) []parser.Suggestion {
	return suggestPlatforms(s, cedexis.PlatformsTypeCommunity, nil)
}

func suggestCloudPlatforms(s string) []parser.Suggestion {
	cat := cedexis.PlatformCategoryCloudComputing
	return suggestPlatforms(s, cedexis.PlatformsTypeCommunity, &cat)
}

func suggestPrivatePlatforms(s string) []parser.Suggestion {
	return suggestPlatforms(s, cedexis.PlatformsTypePrivate, nil)
}

func suggestPlatforms(s string, t cedexis.PlatformType, category *cedexis.PlatformCategory) []parser.Suggestion {
	platforms := getPlatforms(t, category)

	result := make([]parser.Suggestion, 0, len(platforms))
	for _, p := range platforms {
		result = append(result, parser.Suggestion{Text: "\"" + *p.Name + "\""})
	}

	return parser.FilterContains(result, s, true)
}

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

func suggestSonarMarket(s string) []parser.Suggestion {
	result := []parser.Suggestion{
		{Text: string(cedexis.MarketGlobal), Description: "Global"},
		{Text: string(cedexis.MarketAsia), Description: "Asia Pacific"},
		{Text: string(cedexis.MarketEurope), Description: "Europe"},
		{Text: string(cedexis.MarketNorthAmerica), Description: "North America"},
	}

	return parser.FilterHasPrefix(result, s, true)
}

func suggestAlerts(s string) []parser.Suggestion {
	alerts, err := getAlerts()
	if err != nil {
		return nil
	}

	result := make([]parser.Suggestion, 0, len(alerts))
	for _, a := range alerts {
		result = append(result, parser.Suggestion{Text: "\"" + *a.Name + "\""})
	}

	return parser.FilterContains(result, s, true)
}

func suggestAlertTypes(s string) []parser.Suggestion {
	result := []parser.Suggestion{
		{Text: cedexis.AlertTypeRadar.String(), Description: "Radar"},
		{Text: cedexis.AlertTypeSonar.String(), Description: "Sonar"},
	}

	return parser.FilterHasPrefix(result, s, true)
}

func suggestAlertChange(s string) []parser.Suggestion {
	result := []parser.Suggestion{
		{Text: cedexis.AlertChangeAny.String(), Description: "Any"},
		{Text: cedexis.AlertChangeToUp.String(), Description: "down-->up"},
		{Text: cedexis.AlertChangeToDown.String(), Description: "up-->down"},
	}

	return parser.FilterHasPrefix(result, s, true)
}

func suggestAlertTiming(s string) []parser.Suggestion {
	result := []parser.Suggestion{
		{Text: cedexis.AlertTimingBoth.String(), Description: "Immediate and summary"},
		{Text: cedexis.AlertTimingImmediate.String(), Description: "Immediate only"},
		{Text: cedexis.AlertTimingSummary.String(), Description: "Summary only"},
	}

	return parser.FilterHasPrefix(result, s, true)
}

func suggestAlertInterval(s string) []parser.Suggestion {
	result := []parser.Suggestion{
		{Text: "5", Description: "5 minutes"},
		{Text: "10", Description: "10 minutes"},
		{Text: "15", Description: "15 minutes"},
		{Text: "30", Description: "30 minutes"},
		{Text: "60", Description: "1 hour"},
	}

	return parser.FilterHasPrefix(result, s, true)
}

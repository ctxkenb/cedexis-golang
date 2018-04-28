package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/ctxkenb/cedexis-golang/cedexis"
	"github.com/ctxkenb/cedexis-golang/cli/parser"
	"golang.org/x/crypto/ssh/terminal"
)

func completer(d prompt.Document) []prompt.Suggest {
	pSuggestions := cliParser.Suggest(d.TextBeforeCursor())

	result := make([]prompt.Suggest, len(pSuggestions))
	for i, s := range pSuggestions {
		result[i] = prompt.Suggest{Text: s.Text, Description: s.Description}
	}
	return result
}

func executor(cmd string) {
	command, err := cliParser.ParseCommand(cmd)

	if err != nil {
		fmt.Println(err)
		return
	}

	if command.Code == int(CmdExit) {
		fmt.Println("Bye!")
		os.Exit(0)
	}

	if (command.Code & 0xff) == int(CmdFragList) {
		var t *Table

		if ((command.Code >> 8) & 0xff) == int(CmdFragPlatform) {
			var platforms []*cedexis.PlatformInfo
			switch CommandCode(command.Code) {
			case CmdListCommunityPlatforms:
				platforms = getPlatforms(cedexis.PlatformsTypeCommunity, nil)
				break
			case CmdListPrivatePlatforms:
				platforms = getPlatforms(cedexis.PlatformsTypePrivate, nil)
				break
			}

			if filter, ok := command.Args[argFilter]; ok {
				platforms, err = filterPlatforms(platforms, filter)
			}

			t = platformsToTable(platforms)
		}

		w, _, err := terminal.GetSize(int(os.Stdout.Fd()))
		if err != nil || w == 0 {
			w = 80
		}

		t.Print(os.Stdout, w)
	} else if command.Code == int(CmdCreateCloudPlatform) {
		shortName := command.Args[argShortName]
		if shortName == "" {
			shortName = strings.Replace(command.Args[argName], " ", "_", -1)
		}

		cat := cedexis.PlatformCategoryCloudComputing
		platformID, err := getPlatformID(command.Args[argRegion], cedexis.PlatformsTypeCommunity, &cat)
		if err != nil {
			fmt.Println(err)
			return
		}

		tags := strings.Split(command.Args[argTags], ",")

		p := cedexis.NewPublicCloudPrivatePlatform(shortName, command.Args[argName],
			command.Args[argDescription], platformID, tags)

		p.SonarConfig, err = parseSonarConfig(command.Args)
		if err != nil {
			fmt.Println(err)
			return
		}

		if p.SonarConfig.Enabled != nil && *p.SonarConfig.Enabled == true && p.SonarConfig.URL == nil {
			fmt.Println("Error: Sonar requires URL to be enabled")
			return
		}

		p, err = cClient.CreatePrivatePlatform(p)
		if err != nil {
			fmt.Println(err)
			return
		}

		resetCache()
	} else if command.Code == int(CmdDeletePlatform) {
		platformID, err := getPlatformID(command.Args[argName], cedexis.PlatformsTypePrivate, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		err = cClient.DeletePrivatePlatform(platformID)
		if err != nil {
			fmt.Println(err)
		}

		resetCache()
	}

}

func parseSonarConfig(vars map[string]string) (*cedexis.SonarConfig, error) {
	sonarEnabled, err := parseBool(vars[argSonarEnabled])
	if err != nil {
		return nil, err
	}

	sonarURL, err := parseURL(vars[argSonarURL])
	if err != nil {
		return nil, err
	}

	sonarPollInterval, err := parseInt(vars[argSonarPollInterval])
	if err != nil {
		return nil, err
	}

	sonarTimeout, err := parseInt(vars[argSonarPollInterval])
	if err != nil {
		return nil, err
	}

	sonarIgnoreSSLErrors, err := parseBool(vars[argSonarIgnoreSSLErrors])
	if err != nil {
		return nil, err
	}

	sonarMaintenanceMode, err := parseBool(vars[argSonarMaintenanceMode])
	if err != nil {
		return nil, err
	}

	sonarMethod := stringOrNil(vars[argSonarMethod])
	if sonarMethod == nil {
		val := "GET"
		sonarMethod = &val
	}

	sonarHost := stringOrNil(vars[argSonarHost])
	sonarMarket := stringOrNil(vars[argSonarMarket])
	sonarRequestContentType := stringOrNil(vars[argSonarRequestContentType])
	sonarResponseBodyMatch := stringOrNil(vars[argSonarResponseBodyMatch])
	sonarResponseMatchType := stringOrNil(vars[argSonarResponseMatchType])

	return &cedexis.SonarConfig{
		Enabled:             sonarEnabled,
		URL:                 sonarURL,
		PollIntervalSeconds: sonarPollInterval,
		Timeout:             sonarTimeout,
		Method:              sonarMethod,
		IgnoreSSLErrors:     sonarIgnoreSSLErrors,
		MaintenanceMode:     sonarMaintenanceMode,
		Host:                sonarHost,
		Market:              sonarMarket,
		RequestContentType:  sonarRequestContentType,
		ResponseBodyMatch:   sonarResponseBodyMatch,
		ResponseMatchType:   sonarResponseMatchType,
	}, nil
}

func stringOrNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func parseBool(s string) (*bool, error) {
	if s == "" {
		return nil, nil
	}
	b, err := strconv.ParseBool(s)
	return &b, err
}

func parseURL(s string) (*string, error) {
	if s == "" {
		return nil, nil
	}
	_, err := url.ParseRequestURI(s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func parseInt(s string) (*int, error) {
	if s == "" {
		return nil, nil
	}
	i64, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return nil, err
	}

	i := int(i64)

	return &i, nil
}

func resetCache() {
	platforms = map[cedexis.PlatformType][]*cedexis.PlatformInfo{}
}

var cClient *cedexis.Client
var platforms = map[cedexis.PlatformType][]*cedexis.PlatformInfo{}

var cliParser = parser.New(commandSpec)

func main() {
	ctx := context.Background()

	cClient = cedexis.NewClient(ctx, os.Getenv("CEDEXIS_KEY_NAME"), os.Getenv("CEDEXIS_KEY_SECRET"))
	err := cClient.Ping()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("### Cedexis interactive shell ###")
	p := prompt.New(
		executor,
		completer,
		prompt.OptionTitle("cedexis-cli: interactive shell for cedexis"),
		prompt.OptionPrefix("> "),
	)
	p.Run()
}

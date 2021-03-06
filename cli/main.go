package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	if command.Handler != nil {
		command.Handler(command)
		return
	}

	fmt.Printf("No handler for command %s\n", CommandCode(command.Code).String())
}

func handleCreatePlatform(command *parser.Command) {
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

	sonar, err := parseSonarConfig(command.Args)
	if err != nil {
		fmt.Println(err)
		return
	}

	if sonar.Enabled != nil && *sonar.Enabled == true && sonar.URL == nil {
		fmt.Println("Error: Sonar requires URL to be enabled")
		return
	}

	err = createPlatform(command.Args[argName], shortName, command.Args[argDescription], platformID, tags, sonar)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func handleCreateAlert(command *parser.Command) {
	alertType, err := cedexis.ParseAlertType(command.Args[argType])
	if err != nil {
		fmt.Println(err)
		return
	}

	platformID, err := getPlatformID(command.Args[argPlatform], cedexis.PlatformsTypePrivate, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	change, err := cedexis.ParseAlertChange(command.Args[argChange])
	if err != nil {
		fmt.Println(err)
		return
	}

	timing, err := cedexis.ParseAlertTiming(command.Args[argTiming])
	if err != nil {
		fmt.Println(err)
		return
	}

	emails := strings.Split(command.Args[argEmails], ",")

	intervalMins, err := parseInt(command.Args[argInterval])
	if err != nil {
		fmt.Println(err)
		return
	}

	interval := 5 * 60
	if intervalMins != nil {
		interval = (*intervalMins) * 60
	}

	err = createAlert(command.Args[argName], alertType, platformID, change, timing, emails, interval)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func handleCreateApplication(command *parser.Command) {
	appType := command.Args[argType]

	fallbackCname := command.Args[argFallbackCname]

	availabilityThreshold, err := parseInt(command.Args[argAvailabilityThreshold])
	if err != nil {
		fmt.Println(err)
		return
	}

	targetPlatformID, err := getPlatformID(command.Args[argTargetPlatform], cedexis.PlatformsTypeAll, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	targetCname := command.Args[argTargetCname]

	sonarEnabled, err := parseBool(command.Args[argSonarEnabled])
	if err != nil {
		fmt.Println(err)
		return
	}
	if sonarEnabled == nil {
		enabled := true
		sonarEnabled = &enabled
	}

	err = createApp(command.Args[argName], command.Args[argDescription], appType, fallbackCname, *availabilityThreshold, targetPlatformID, targetCname, *sonarEnabled)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func handleCreateZone(command *parser.Command) {
	var err error

	name := command.Args[argName]
	description := command.Args[argDescription]
	tags := strings.Split(command.Args[argTags], ",")

	var zoneFile *string
	fileName := command.Args[argZoneFile]
	if len(fileName) > 0 {
		zoneData, err := ioutil.ReadFile(fileName)
		if err != nil {
			fmt.Println(err)
			return
		}

		zoneStr := string(zoneData)
		zoneFile = &zoneStr
	}

	err = createZone(name, description, tags, zoneFile)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func handleList(command *parser.Command) {
	var t *Table
	var err error

	objType := CommandFrag((command.Code >> 8) & 0xff)
	switch objType {
	case CmdFragPlatform:
		var platforms []*cedexis.PlatformInfo
		switch CommandCode(command.Code) {
		case CmdListCommunityPlatforms:
			platforms = getPlatforms(cedexis.PlatformsTypeCommunity, nil)
			break
		case CmdListPrivatePlatforms:
			platforms = getPlatforms(cedexis.PlatformsTypePrivate, nil)
			break
		default:
			fmt.Println("Unrecognized command: " + commandCodeNames[CommandCode(command.Code)])
			return
		}

		if filter, ok := command.Args[argFilter]; ok {
			platforms, err = filterPlatforms(platforms, filter)
		}

		t = platformsToTable(platforms)
	case CmdFragAlert:
		alerts, err := getAlerts()
		if err != nil {
			fmt.Println(err)
			return
		}

		if filter, ok := command.Args[argFilter]; ok {
			alerts, err = filterAlerts(alerts, filter)
		}

		t = alertsToTable(alerts)
	case CmdFragApp:
		apps, err := getApps()
		if err != nil {
			fmt.Println(err)
			return
		}

		if filter, ok := command.Args[argFilter]; ok {
			apps, err = filterApps(apps, filter)
		}

		t = appsToTable(apps)
	case CmdFragZone:
		zones, err := getZones()
		if err != nil {
			fmt.Println(err)
			return
		}

		if filter, ok := command.Args[argFilter]; ok {
			zones, err = filterZones(zones, filter)
		}

		t = zonesToTable(zones)
	default:
		fmt.Println("Unrecognized command: " + commandCodeNames[CommandCode(command.Code)])
		return
	}

	w, _, err := terminal.GetSize(int(os.Stdout.Fd()))
	if err != nil || w == 0 {
		w = 80
	}

	t.Print(os.Stdout, w)
}

func handleShow(command *parser.Command) {
	var err error
	var obj interface{}
	objType := CommandFrag((command.Code >> 8) & 0xff)
	switch objType {
	case CmdFragPlatform:
		pID, err := getPlatformID(command.Args[argName], cedexis.PlatformsTypeAll, nil)
		if err != nil {
			fmt.Println(err)
			return
		}

		obj, err = cClient.GetPrivatePlatform(pID)
		if err != nil {
			fmt.Println(err)
			return
		}
	case CmdFragAlert:
		obj, err = getAlert(command.Args[argName])
		if err != nil {
			fmt.Println(err)
			return
		}
	case CmdFragApp:
		obj, err = getApp(command.Args[argName])
		if err != nil {
			fmt.Println(err)
			return
		}
	case CmdFragZone:
		obj, err = getZone(command.Args[argName])
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(data))
}

func handleDeletePlatform(command *parser.Command) {
	var err error
	if command.Args[argName] != "" {
		err = deletePlatform(command.Args[argName], cedexis.PlatformsTypePrivate)
	} else if command.Args[argFilter] != "" {
		platforms := getPlatforms(cedexis.PlatformsTypePrivate, nil)
		platforms, err = filterPlatforms(platforms, command.Args[argFilter])
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, p := range platforms {
			err = deletePlatform(*p.Name, cedexis.PlatformsTypePrivate)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

	} else {
		err = fmt.Errorf("No platform specified by name or by filter")
	}

	if err != nil {
		fmt.Println(err)
	}
}

func handleDeleteAlert(command *parser.Command) {
	var err error
	if command.Args[argName] != "" {
		err = deleteAlert(command.Args[argName])
	} else if command.Args[argFilter] != "" {
		alerts, err := getAlerts()
		if err != nil {
			fmt.Println(err)
			return
		}

		alerts, err = filterAlerts(alerts, command.Args[argFilter])
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, p := range alerts {
			err = deleteAlert(*p.Name)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

	} else {
		err = fmt.Errorf("No alert specified by name or by filter")
	}

	if err != nil {
		fmt.Println(err)
	}
}

func handleDeleteApplication(command *parser.Command) {
	var err error
	if command.Args[argName] != "" {
		err = deleteApp(command.Args[argName])
	} else if command.Args[argFilter] != "" {
		apps, err := getApps()
		if err != nil {
			fmt.Println(err)
			return
		}

		apps, err = filterApps(apps, command.Args[argFilter])
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, p := range apps {
			err = deleteApp(*p.Name)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

	} else {
		err = fmt.Errorf("No app specified by name or by filter")
	}

	if err != nil {
		fmt.Println(err)
	}
}

func handleDeleteZone(command *parser.Command) {
	var err error
	if command.Args[argName] != "" {
		err = deleteZone(command.Args[argName])
	} else if command.Args[argFilter] != "" {
		zones, err := getZones()
		if err != nil {
			fmt.Println(err)
			return
		}

		zones, err = filterZones(zones, command.Args[argFilter])
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, p := range zones {
			err = deleteZone(*p.DomainName)
			if err != nil {
				fmt.Println(err)
				return
			}
		}

	} else {
		err = fmt.Errorf("No app specified by name or by filter")
	}

	if err != nil {
		fmt.Println(err)
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
	sonarMarket := (*cedexis.Market)(stringOrNil(vars[argSonarMarket]))
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

var cClient *cedexis.Client

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

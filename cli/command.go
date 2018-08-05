package main

import "github.com/ctxkenb/cedexis-golang/cli/parser"

//
// Defines commands of the format:
//   1 or more command fragments to make a command (e.g. create plaform cloud)
//   1 or more positional arguments (e.g. <name>)
//   1 or more optional arguments (e.g. -force <value>)
//

// CommandCode represents a fully specified command (command, sub-command, sub-sub-command)
type CommandCode int

// CommandFrag represents one element of a command
type CommandFrag int

const (
	// CmdFragCreate represents the "create" command
	CmdFragCreate CommandFrag = iota

	// CmdFragList represents the "list" command
	CmdFragList

	// CmdFragShow represents the "show" command
	CmdFragShow

	// CmdFragDelete represents the "delete" command
	CmdFragDelete

	// CmdFragExit represents the "exit" command
	CmdFragExit

	// CmdFragPlatform represents the "xxx platform xxx" sub-command
	CmdFragPlatform

	// CmdFragAlert represents the "xxx alert xxx" sub-command
	CmdFragAlert

	// CmdFragCloud represents the "xxx xxx cloud" sub-command
	CmdFragCloud

	// CmdFragPublic represents the "xxx xxx public" sub-command
	CmdFragPublic

	// CmdFragPrivate represents the "xxx xxx private" sub-command
	CmdFragPrivate

	// CmdFragApp represents the "xxx application" sub-command
	CmdFragApp
)

const (
	// CmdNone represents the default (none) command
	CmdNone CommandCode = 0

	// CmdCreateCloudPlatform represents command "create platform cloud"
	CmdCreateCloudPlatform CommandCode = CommandCode(int(CmdFragCreate | (CmdFragPlatform << 8) | (CmdFragCloud << 16)))

	// CmdDeletePlatform represents command "delete platform <name>"
	CmdDeletePlatform CommandCode = CommandCode(int(CmdFragDelete | (CmdFragPlatform << 8)))

	// CmdListCommunityPlatforms represents command "list platform public"
	CmdListCommunityPlatforms CommandCode = CommandCode(int(CmdFragList | (CmdFragPlatform << 8) | (CmdFragPublic << 16)))

	// CmdListPrivatePlatforms represents command "list platform private"
	CmdListPrivatePlatforms CommandCode = CommandCode(int(CmdFragList | (CmdFragPlatform << 8) | (CmdFragPrivate << 16)))

	// CmdShowPlatform represents command "show platform"
	CmdShowPlatform CommandCode = CommandCode(int(CmdFragShow | (CmdFragPlatform << 8)))

	// CmdCreateAlert represents commdn "create alert"
	CmdCreateAlert CommandCode = CommandCode(int(CmdFragCreate | (CmdFragAlert << 8)))

	// CmdDeleteAlert represents commdn "delete alert"
	CmdDeleteAlert CommandCode = CommandCode(int(CmdFragDelete | (CmdFragAlert << 8)))

	// CmdListAlerts represents command "list alert"
	CmdListAlerts CommandCode = CommandCode(int(CmdFragList | (CmdFragAlert << 8)))

	// CmdShowAlert reprsents command "show alert"
	CmdShowAlert CommandCode = CommandCode(int(CmdFragShow | (CmdFragAlert << 8)))

	// CmdListApplications represents command "list applications"
	CmdListApplications CommandCode = CommandCode(int(CmdFragList | (CmdFragApp << 8)))

	// CmdShowApplication represents command "show application"
	CmdShowApplication CommandCode = CommandCode(int(CmdFragShow | (CmdFragApp << 8)))

	// CmdExit represents "exit" command
	CmdExit CommandCode = CommandCode(int(CmdFragExit))
)

var commandCodeNames = map[CommandCode]string{
	CmdNone:                   "CmdNone",
	CmdCreateCloudPlatform:    "CmdCreateCloudPlatform",
	CmdDeletePlatform:         "CmdDeletePlatform",
	CmdListCommunityPlatforms: "CmdListPublicPlatforms",
	CmdListPrivatePlatforms:   "CmdListPrivatePlatforms",
	CmdShowPlatform:           "CmdShowPlatform",
	CmdCreateAlert:            "CmdCreateAlert",
	CmdDeleteAlert:            "CmdDeleteAlert",
	CmdListAlerts:             "CmdListAlerts",
	CmdShowAlert:              "CmdShowAlert",
	CmdListApplications:       "CmdListApplications",
	CmdShowApplication:        "CmdShowApplication",
	CmdExit:                   "CmdExit",
}

func (c CommandCode) String() string {
	return commandCodeNames[c]
}

const (
	argName                    string = "name"
	argShortName               string = "shortName"
	argDescription             string = "description"
	argFilter                  string = "filter"
	argRegion                  string = "region"
	argTags                    string = "tags"
	argSonarEnabled            string = "sonarEnabled"
	argSonarURL                string = "sonarURL"
	argSonarPollInterval       string = "sonarPollInterval"
	argSonarTimeout            string = "sonarTimeout"
	argSonarMethod             string = "sonarMethod"
	argSonarIgnoreSSLErrors    string = "sonarIgnoreSSLErrors"
	argSonarMaintenanceMode    string = "sonarMaintenanceMode"
	argSonarHost               string = "sonarHost"
	argSonarMarket             string = "sonarMarket"
	argSonarRequestContentType string = "sonarRequestContentType"
	argSonarResponseBodyMatch  string = "sonarResponseBodyMatch"
	argSonarResponseMatchType  string = "sonarResponseMatchType"
	argType                    string = "type"
	argPlatform                string = "platform"
	argChange                  string = "change"
	argEmails                  string = "emails"
	argTiming                  string = "timing"
	argInterval                string = "interval"
)

var commandSpec = map[string]parser.CommandFrag{
	"create": {Desc: "Creates platforms, alerts, etc",
		Args: map[string]parser.NamedArg{
			argShortName:   {Desc: "Set the shortname"},
			argDescription: {Desc: "Set the description"},
		},
		Sub: map[string]parser.CommandFrag{
			"platform": {Desc: "Create a new platform", Sub: map[string]parser.CommandFrag{
				"cloud": {Desc: "New public cloud platform",
					Code:    int(CmdCreateCloudPlatform),
					Handler: handleCreatePlatform,
					PosArgs: []parser.PosArg{{Name: argName, Desc: "Name of platform"}},
					Args: map[string]parser.NamedArg{
						argRegion:                  {Desc: "Set the public cloud region", Suggest: suggestCloudPlatforms},
						argTags:                    {Desc: "Set tags on the new platform"},
						argSonarEnabled:            {Desc: "Enable sonar health-checks"},
						argSonarURL:                {Desc: "URL to check"},
						argSonarPollInterval:       {Desc: "Seconds between checks"},
						argSonarTimeout:            {Desc: "Timeout for health-check"},
						argSonarMethod:             {Desc: "HTTP method for health-check"},
						argSonarIgnoreSSLErrors:    {Desc: "Accept invalid SSL certs"},
						argSonarMaintenanceMode:    {Desc: "Force state down to Openmix"},
						argSonarHost:               {Desc: "Override host from URL"},
						argSonarMarket:             {Desc: "Source for health-checks", Suggest: suggestSonarMarket},
						argSonarRequestContentType: {Desc: "Request Content-Type header"},
						argSonarResponseBodyMatch:  {Desc: "Any string"},
						argSonarResponseMatchType:  {Desc: "Pass vs fail based on body match"},
					}},
			}},
			"alert": {Desc: "Create a new alert",
				Code:    int(CmdCreateAlert),
				Handler: handleCreateAlert,
				PosArgs: []parser.PosArg{{Name: argName, Desc: "Name of alert"}},
				Args: map[string]parser.NamedArg{
					argType:     {Desc: "The alert type", Suggest: suggestAlertTypes},
					argPlatform: {Desc: "Name of the platform", Suggest: suggestPrivatePlatforms},
					argChange:   {Desc: "Event triggering alert", Suggest: suggestAlertChange},
					argEmails:   {Desc: "Notification targets"},
					argTiming:   {Desc: "Summary or immediate notification", Suggest: suggestAlertTiming},
					argInterval: {Desc: "Alert gap (in minutes)", Suggest: suggestAlertInterval},
				},
			},
		},
	},
	"delete": {Desc: "Deletes platforms, alerts, etc",
		Sub: map[string]parser.CommandFrag{
			"platform": {Desc: "Delete a platform",
				Handler: handleDeletePlatform,
				Code:    int(CmdDeletePlatform),
				PosArgs: []parser.PosArg{{Name: argName, Desc: "Name of platform", Suggest: suggestPrivatePlatforms}},
			},
			"alert": {Desc: "Delete an alert",
				Handler: handleDeleteAlert,
				Code:    int(CmdDeleteAlert),
				PosArgs: []parser.PosArg{{Name: argName, Desc: "Name of alert", Suggest: suggestAlerts}},
			},
		},
	},
	"list": {Desc: "List platforms, alerts, etc",
		Handler: handleList,
		Sub: map[string]parser.CommandFrag{
			"platform": {Desc: "List platforms",
				Args: map[string]parser.NamedArg{argFilter: {Desc: "Regex filter"}},
				Sub: map[string]parser.CommandFrag{
					"community": {Desc: "List community platforms", Code: int(CmdListCommunityPlatforms)},
					"private":   {Desc: "List private platforms", Code: int(CmdListPrivatePlatforms)},
				}},
			"alert":       {Desc: "List alerts", Code: int(CmdListAlerts)},
			"application": {Desc: "List Openmix apps", Code: int(CmdListApplications)},
		},
	},
	"show": {Desc: "Show details",
		Handler: handleShow,
		Sub: map[string]parser.CommandFrag{
			"platform": {Desc: "Show platform",
				Code:    int(CmdShowPlatform),
				PosArgs: []parser.PosArg{{Name: argName, Desc: "Name of platform", Suggest: suggestPrivatePlatforms}},
			},
			"alert": {Desc: "Show alert",
				Code:    int(CmdShowAlert),
				PosArgs: []parser.PosArg{{Name: argName, Desc: "Name of alert", Suggest: suggestAlerts}},
			},
			"application": {Desc: "Show Openmix app",
				Code:    int(CmdShowApplication),
				PosArgs: []parser.PosArg{{Name: argName, Desc: "Name of application", Suggest: suggestApps}},
			},
		},
	},
	"exit": {Desc: "Exit", Code: int(CmdExit)},
}

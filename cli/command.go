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

	// CmdFragExit represents the "exit" command
	CmdFragExit

	// CmdFragPlatform represents the "xxx platform xxx" sub-command
	CmdFragPlatform

	// CmdFragCloud represents the "xxx xxx cloud" sub-command
	CmdFragCloud

	// CmdFragPublic represents the "xxx xxx public" sub-command
	CmdFragPublic

	// CmdFragPrivate represents the "xxx xxx private" sub-command
	CmdFragPrivate
)

const (
	// CmdNone represents the default (none) command
	CmdNone CommandCode = 0

	// CmdCreateCloudPlatform represents command "create platform cloud"
	CmdCreateCloudPlatform CommandCode = CommandCode(int(CmdFragCreate | (CmdFragPlatform << 8) | (CmdFragCloud << 16)))

	// CmdListCommunityPlatforms represents command "list platform public"
	CmdListCommunityPlatforms CommandCode = CommandCode(int(CmdFragList | (CmdFragPlatform << 8) | (CmdFragPublic << 16)))

	// CmdListPrivatePlatforms represents command "list platform private"
	CmdListPrivatePlatforms CommandCode = CommandCode(int(CmdFragList | (CmdFragPlatform << 8) | (CmdFragPrivate << 16)))

	// CmdExit represents "exit" command
	CmdExit CommandCode = CommandCode(int(CmdFragExit))
)

var commandCodeNames = map[CommandCode]string{
	CmdNone:                   "CmdNone",
	CmdCreateCloudPlatform:    "CmdCreateCloudPlatform",
	CmdListCommunityPlatforms: "CmdListPublicPlatforms",
	CmdListPrivatePlatforms:   "CmdListPrivatePlatforms",
	CmdExit:                   "CmdExit",
}

func (c CommandCode) String() string {
	return commandCodeNames[c]
}

var commandSpec = map[string]parser.CommandFrag{
	"create": {Desc: "Creates platforms, etc",
		Args: map[string]parser.NamedArg{
			"shortName":   {Desc: "Set the shortname"},
			"description": {Desc: "Set the description"},
		},
		Sub: map[string]parser.CommandFrag{
			"platform": {Desc: "Create a new platform", Sub: map[string]parser.CommandFrag{
				"cloud": {Desc: "New public cloud platform",
					Code:    int(CmdCreateCloudPlatform),
					PosArgs: []parser.PosArg{{Name: "name", Desc: "Name of platform"}},
					Args: map[string]parser.NamedArg{
						"region": {Desc: "Set the public cloud region", Suggest: suggestCloudPlatforms},
						"tags":   {Desc: "Set tags on the new platform"},
					}},
			}},
		},
	},
	"list": {Desc: "List platforms, etc",
		Args: map[string]parser.NamedArg{"filter": {Desc: "Regex filter"}},
		Sub: map[string]parser.CommandFrag{
			"platform": {Desc: "List platforms", Sub: map[string]parser.CommandFrag{
				"community": {Desc: "List community platforms", Code: int(CmdListCommunityPlatforms)},
				"private":   {Desc: "List private platforms", Code: int(CmdListPrivatePlatforms)},
			}},
		},
	},
	"exit": {Desc: "Exit", Code: int(CmdExit)},
}

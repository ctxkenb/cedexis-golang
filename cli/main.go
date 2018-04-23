package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
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

			if filter, ok := command.Args["filter"]; ok {
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
		shortName := command.Args["shortName"]
		if shortName == "" {
			shortName = strings.Replace(command.Args["name"], " ", "_", -1)
		}

		cat := cedexis.PlatformCategoryCloudComputing
		availPlatforms := getPlatforms(cedexis.PlatformsTypeCommunity, &cat)
		platformID := -1
		for _, p := range availPlatforms {
			if *p.Name == command.Args["region"] {
				platformID = *p.ID
			}
		}
		if platformID == -1 {
			fmt.Printf("Region '%v' not found\n", command.Args["region"])
			return
		}

		tags := strings.Split(command.Args["tags"], ",")

		p := cedexis.NewPublicCloudPrivatePlatform(shortName, command.Args["name"],
			command.Args["description"], platformID, tags)

		//fmt.Printf("%#v\n", p)

		p, err = cClient.CreatePrivatePlatform(p)
		if err != nil {
			fmt.Println(err)
		}

		// reset cache
		platforms = map[cedexis.PlatformType][]*cedexis.PlatformInfo{}
	}

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

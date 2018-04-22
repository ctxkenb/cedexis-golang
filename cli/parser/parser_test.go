package parser

import (
	"testing"
)

const (
	CmdExit int = 1
	CmdCreateCloudPlatform
)

var commandSpec = map[string]CommandFrag{
	"create": {Desc: "Creates platforms, etc",
		Args: map[string]NamedArg{
			"shortName":   {Desc: "Overrides the shortname"},
			"description": {Desc: "Sets the description"},
		},
		Sub: map[string]CommandFrag{
			"platform": {Desc: "Create a new platform", Sub: map[string]CommandFrag{
				"cloud": {Desc: "New public cloud platform",
					Code:    int(CmdCreateCloudPlatform),
					PosArgs: []PosArg{{Name: "name", Desc: "Name of platform"}},
					Args: map[string]NamedArg{
						"region": {Desc: "Sets the public cloud region"},
						"tags":   {Desc: "Sets tags on the new platform"},
						"flag":   {Desc: "Test Flag", Flag: true},
					}},
			}},
		},
	},
	"exit": {Desc: "Exit", Code: int(CmdExit)},
}

var parser = New(commandSpec)

func TestTokenize(t *testing.T) {
	tests := []struct {
		cmd    string
		tokens []string
		ws     bool
	}{
		{"", []string{}, false},
		{"exit", []string{"exit"}, false},
		{"cmd \"a space\"", []string{"cmd", "a space"}, false},
		{"cmd -arg", []string{"cmd", "-arg"}, false},
		{"cmd  -arg", []string{"cmd", "-arg"}, false},
		{" cmd -arg", []string{"cmd", "-arg"}, false},
		{"cmd -arg ", []string{"cmd", "-arg"}, true},
		{"cmd \"ff ", []string{"cmd", "ff "}, false},
	}

	for _, test := range tests {
		tokens, ws := tokenize(test.cmd)

		if len(tokens) != len(test.tokens) {
			t.Errorf("Incorrect number of tokens for \"%v\", got %d, want %d.", test.cmd, len(tokens), len(test.tokens))
		}
		for i := range tokens {
			if tokens[i] != test.tokens[i] {
				t.Errorf("Incorrect token for \"%v\", got \"%v\", want \"%v\"", test.cmd, tokens[i], test.tokens[i])
			}
		}
		if ws != test.ws {
			t.Errorf("Incorrect whitespace flag for \"%v\", got \"%v\", want \"%v\"", test.cmd, ws, test.ws)
		}
	}
}

func TestParseCommand(t *testing.T) {
	tests := []struct {
		cmd     string
		isError bool
		code    int
		args    *map[string]string
	}{
		{"", true, 0, nil},
		{"exit", false, CmdExit, nil},
		{"create platform cloud \"a space\"", false, CmdCreateCloudPlatform,
			&map[string]string{"name": "a space"},
		},
		{"create platform cloud \"a space\" -flag", false, CmdCreateCloudPlatform,
			&map[string]string{"name": "a space", "flag": ""},
		},
		{"create platform", true, 0, nil},
		{"create platform cloud", true, 0, nil},                   // missing pos. param
		{"create platform cloud -shortName f", true, 0, nil},      // missing pos. param
		{"create platform cloud name -shortName", true, 0, nil},   // missing param value
		{"create platform cloud name -unknown val", true, 0, nil}, // unknown param value
		{"create platform cloud name extra", true, 0, nil},        // extra named value
		{"create platform xxx", true, 0, nil},
		{"create platform cloud \"long name\" -shortName f", false, CmdCreateCloudPlatform,
			&map[string]string{"shortName": "f", "name": "long name"},
		}, // missing pos. param
	}

	for _, test := range tests {
		cmd, error := parser.ParseCommand(test.cmd)

		if test.isError != (error != nil) {
			t.Errorf("Incorrect error for '%v', got '%v', wanted error: %v", test.cmd, error, test.isError)
		}
		if error == nil && test.code != cmd.Code {
			t.Errorf("Incorrect code for '%v', got %d, want %d", test.cmd, cmd.Code, test.code)
		}

		if cmd == nil {
			// remaining tests require cmd to be returned
			continue
		}

		if (test.args == nil) != (len(cmd.Args) == 0) {
			want := 0
			if test.args != nil {
				want = len(*test.args)
			}

			t.Errorf("Incorrect arg count returned for '%v', got %d, want %d", test.cmd, len(cmd.Args), want)
		}
		if test.args != nil {
			for k, v := range *test.args {
				tv, ok := cmd.Args[k]
				if !ok {
					t.Errorf("Missing arg for '%v', got <nil>, want %s", test.cmd, k)
				} else if tv != v {
					t.Errorf("Incorrect arg value for '%v', got '%v', want '%v'", test.cmd, tv, v)
				}
			}
		}
	}
}
func TestSuggest(t *testing.T) {
	tests := []struct {
		cmd            string
		hasSuggestions bool
	}{
		{"", true},
		{"cr", true},
		{"junk", false},
		{"create ", true},
		{"create bogus cloud", false},
		{"create platform cloud ", true},
		{"create platform cloud fred", true},
		{"create platform cloud fred ", true},
		{"create platform cloud \"fred ", false}, // spaces break the prompt library
		{"create platform cloud fred -", true},
		{"create platform cloud fred -shortName ", true},
		{"create platform cloud fred -shortName f", true},
		{"create platform cloud fred -shortName \"f ", false}, // spaces break the prompt library
		{"create platform cloud fred -shortName f -ta", true},
		{"create platform cloud fred extra", false},
		{"create platform cloud \"a space\" -fla", true},
		{"create platform cloud \"a space\" -flag ", true},
	}

	for _, test := range tests {
		s := parser.Suggest(test.cmd)

		if (len(s) > 0) != test.hasSuggestions {
			t.Errorf("Incorrect suggestions for '%v'. got '%v', wanted %v", test.cmd, s, test.hasSuggestions)
		}
	}
}

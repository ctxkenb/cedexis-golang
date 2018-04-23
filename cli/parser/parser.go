package parser

import (
	"fmt"
	"strings"
)

// Command represents a parsed command and all arguments
type Command struct {
	// Code is a unique identifier for a command, zero is reserved and must not be used
	Code int
	Args map[string]string
}

// CommandFrag specifies a fragment of a command.
//
// For the command "create plaform cloud", these are all fragments: "create", "platform", "cloud". Some
// parts of the fragment only make sense to specify on the leaf fragments ("cloud" in this example).
type CommandFrag struct {
	// Desc is a description of the fragment, used to assist user
	Desc string

	// PosArgs are the positional arguments for the command (leaf nodes only)
	PosArgs []PosArg
	Sub     map[string]CommandFrag
	Args    map[string]NamedArg
	Code    int
}

// PosArg specifies a positional argument, found immediately after the command itself.
type PosArg struct {
	Name    string
	Desc    string
	Suggest SuggestFn
}

// NamedArg specifies a named argument, which can be in any order after named arguments and start with '-'
type NamedArg struct {
	// Description of the argument
	Desc string

	// Function that suggests argument values
	Suggest SuggestFn

	// Flag indicates the arg has no following value
	Flag bool
}

// Parser implements a shell / command-line parser with facilities for auto-complete
type Parser struct {
	Spec map[string]CommandFrag
}

// New creates a new shell / command-line parser, based on a spec
func New(spec map[string]CommandFrag) *Parser {
	return &Parser{spec}
}

// ParseCommand breaks a command string into the command + args
func (p *Parser) ParseCommand(str string) (*Command, error) {
	tokens, _ := tokenize(str)

	// Blank command
	if len(tokens) < 1 {
		return nil, fmt.Errorf("Empty command: \"%v\"", str)
	}

	ptr := CommandFrag{Sub: p.Spec}
	args := map[string]string{}
	namedArg := ""
	possibleArgs := map[string]NamedArg{}
	for i := range tokens {

		if ptr.Code == 0 {
			// We're still looking for a full command
			c, ok := ptr.Sub[strings.ToLower(tokens[i])]
			if !ok {
				return nil, fmt.Errorf("Unrecognised token '%v' in command '%v'", tokens[i], str)
			}

			// Move pointer to current sub-command
			ptr = c

			// Accumulate all possible named args
			for k, v := range ptr.Args {
				possibleArgs[k] = v
			}

		} else if len(args) < len(ptr.PosArgs) {
			argName := ptr.PosArgs[len(args)].Name

			// We're still looking for positional args
			if strings.HasPrefix(tokens[i], "-") {
				return nil, fmt.Errorf(
					"Unexpected named argument '%v' while looking for positional argument '%v'", tokens[i], argName)
			}

			args[argName] = tokens[i]
		} else if namedArg == "" {
			if !strings.HasPrefix(tokens[i], "-") {
				return nil, fmt.Errorf("Expected '-' looking for named argument, not '%v'", tokens[i])
			}

			namedArg = strings.TrimLeft(tokens[i], "-")

			if _, ok := possibleArgs[namedArg]; !ok {
				return nil, fmt.Errorf("Argument '%v' not recognized", namedArg)
			}

			if ptr.Args[namedArg].Flag {
				args[namedArg] = ""
				namedArg = ""
			}
		} else {
			args[namedArg] = tokens[i]
			namedArg = ""
		}
	}

	if len(args) < len(ptr.PosArgs) {
		return nil, fmt.Errorf("Command incomplete, expecting positional argument '%v'", ptr.PosArgs[len(args)].Name)
	}

	if namedArg != "" {
		return nil, fmt.Errorf("Incomplete argument '%v'", namedArg)
	}

	if ptr.Code != 0 {
		return &Command{
			Code: ptr.Code,
			Args: args,
		}, nil
	}

	return nil, fmt.Errorf("Incomplete command '%v'", str)
}

// Suggest suggests what the user should provide when completing a command
func (p *Parser) Suggest(str string) []Suggestion {
	tokens, wsAtEnd := tokenize(str)

	// Blank command
	if len(tokens) < 1 {
		return cmdSuggestions(p.Spec)
	}

	ptr := CommandFrag{Sub: p.Spec}
	args := map[string]string{}
	namedArg := ""
	possibleArgs := map[string]NamedArg{}
	for i := range tokens {
		isLast := i == len(tokens)-1

		if ptr.Code == 0 {
			// We're still looking for a full command
			c, ok := ptr.Sub[strings.ToLower(tokens[i])]
			if !ok && !isLast {
				return []Suggestion{}
			} else if !ok {
				return FilterHasPrefix(cmdSuggestions(ptr.Sub), tokens[i], true)
			}

			// Move pointer to current sub-command
			ptr = c

			// Accumulate all possible named args
			for k, v := range ptr.Args {
				possibleArgs[k] = v
			}

		} else if len(args) < len(ptr.PosArgs) {
			// We're still looking for positional args
			arg := ptr.PosArgs[len(args)]

			// Spaces break the prompt library
			if isLast && !wsAtEnd && !strings.Contains(tokens[i], " ") {
				if arg.Suggest != nil {
					return arg.Suggest(tokens[i])
				} else {
					return []Suggestion{{Text: tokens[i], Description: "<" + arg.Name + "> " + arg.Desc}}
				}
			}

			args[arg.Name] = tokens[i]
		} else if namedArg == "" {
			namedArg = strings.TrimLeft(tokens[i], "-")

			_, ok := possibleArgs[namedArg]
			if isLast && !ok {
				return FilterHasPrefix(argSuggestions(possibleArgs), tokens[i], true)
			}

			if ok && possibleArgs[namedArg].Flag {
				delete(possibleArgs, namedArg)
				namedArg = ""
			}
		} else {
			if isLast && !wsAtEnd {
				// Spaces break the prompt library
				if !strings.Contains(tokens[i], " ") {
					arg, ok := possibleArgs[namedArg]

					if ok && arg.Suggest != nil {
						return arg.Suggest(tokens[i])
					}

					return []Suggestion{{Text: tokens[i], Description: arg.Desc}}
				}

				return []Suggestion{}
			}

			// Prevent this named arg from being suggested, since it's been used
			delete(possibleArgs, namedArg)

			args[namedArg] = tokens[i]
			namedArg = ""
		}
	}

	if wsAtEnd {
		if ptr.Code == 0 {
			return cmdSuggestions(ptr.Sub)
		}

		if len(args) < len(ptr.PosArgs) {
			arg := ptr.PosArgs[len(args)]
			return []Suggestion{{Text: "", Description: "<" + arg.Name + "> " + arg.Desc}}
		}

		if namedArg == "" {
			return argSuggestions(possibleArgs)
		} else if len(args) >= len(ptr.PosArgs) {
			arg, ok := possibleArgs[namedArg]
			if ok {
				return []Suggestion{{Text: "", Description: "<" + namedArg + "> " + arg.Desc}}
			}
		}
	}

	return []Suggestion{}
}

// tokenize breaks a command string into logical tokens, allowing for quotes and escaping
func tokenize(cmd string) ([]string, bool) {
	var result []string

	token := ""
	inQuote := false
	for _, c := range cmd {
		if c == '"' {
			// Quote indicates start or end of token
			if inQuote || token != "" {
				result = append(result, token)
				token = ""
			}
			inQuote = !inQuote
		} else if c == ' ' && !inQuote {
			if token != "" {
				result = append(result, token)
				token = ""
			}
		} else {
			token += string(c)
		}
	}

	if token != "" {
		result = append(result, token)
	}

	wsAtEnd := !inQuote && strings.HasSuffix(cmd, " ")

	return result, wsAtEnd
}

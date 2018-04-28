package parser

import (
	"sort"
	"strings"
)

// Suggestion is a recommended auto-complete for the user
type Suggestion struct {
	Text        string
	Description string
}

// SuggestFn returns a list of suggestions for a given token the user is currently editing
type SuggestFn func(s string) []Suggestion

func cmdSuggestions(cmds map[string]CommandFrag) []Suggestion {
	result := make([]Suggestion, len(cmds))
	i := 0
	for k, v := range cmds {
		result[i].Text = k
		result[i].Description = v.Desc
		i++
	}
	return sortSuggestions(result)
}

func argSuggestions(args map[string]NamedArg) []Suggestion {
	result := make([]Suggestion, len(args))
	i := 0
	for k, v := range args {
		result[i].Text = "-" + k
		result[i].Description = v.Desc
		i++
	}
	return sortSuggestions(result)
}

// FilterHasPrefix filters out suggestions without the given prefix
func FilterHasPrefix(suggestions []Suggestion, prefix string, ignoreCase bool) []Suggestion {
	result := make([]Suggestion, 0, len(suggestions))

	if ignoreCase {
		lowerPrefix := strings.ToLower(prefix)
		for _, s := range suggestions {
			if strings.HasPrefix(strings.ToLower(s.Text), lowerPrefix) {
				result = append(result, s)
			}
		}
	} else {
		for _, s := range suggestions {
			if strings.HasPrefix(s.Text, prefix) {
				result = append(result, s)
			}
		}
	}

	return result
}

// FilterContains filters out suggestions that don't contain str
func FilterContains(suggestions []Suggestion, str string, ignoreCase bool) []Suggestion {
	result := make([]Suggestion, 0, len(suggestions))

	if ignoreCase {
		lowerPrefix := strings.ToLower(str)
		for _, s := range suggestions {
			if strings.Contains(strings.ToLower(s.Text), lowerPrefix) {
				result = append(result, s)
			}
		}
	} else {
		for _, s := range suggestions {
			if strings.Contains(s.Text, str) {
				result = append(result, s)
			}
		}
	}

	return result
}

func sortSuggestions(suggestions []Suggestion) []Suggestion {
	sort.Slice(suggestions, func(i, j int) bool {
		return suggestions[i].Text < suggestions[j].Text
	})

	return suggestions
}

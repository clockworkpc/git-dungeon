package parser

import (
	"fmt"
	"strings"
	"unicode"
)

type ParsedCommand struct {
	Base      string
	Sub       string
	Args      []string
	Flags     map[string]string
	BoolFlags map[string]bool
	Raw       string
}

func tokenize(input string) []string {
	var tokens []string
	var cur strings.Builder
	inQuote := false
	quoteChar := rune(0)

	for _, r := range input {
		switch {
		case inQuote:
			if r == quoteChar {
				inQuote = false
			} else {
				cur.WriteRune(r)
			}
		case r == '"' || r == '\'':
			inQuote = true
			quoteChar = r
		case unicode.IsSpace(r):
			if cur.Len() > 0 {
				tokens = append(tokens, cur.String())
				cur.Reset()
			}
		default:
			cur.WriteRune(r)
		}
	}
	if cur.Len() > 0 {
		tokens = append(tokens, cur.String())
	}
	return tokens
}

func Parse(input string) (ParsedCommand, error) {
	input = strings.TrimSpace(input)
	cmd := ParsedCommand{
		Raw:       input,
		Flags:     make(map[string]string),
		BoolFlags: make(map[string]bool),
	}

	tokens := tokenize(input)
	if len(tokens) == 0 {
		return cmd, fmt.Errorf("empty command")
	}
	if tokens[0] != "git" {
		return cmd, fmt.Errorf("not a git command")
	}
	cmd.Base = "git"
	if len(tokens) < 2 {
		return cmd, fmt.Errorf("missing git subcommand")
	}
	cmd.Sub = tokens[1]

	i := 2
	for i < len(tokens) {
		tok := tokens[i]
		if strings.HasPrefix(tok, "-") {
			if i+1 < len(tokens) && !strings.HasPrefix(tokens[i+1], "-") {
				cmd.Flags[tok] = tokens[i+1]
				i += 2
			} else {
				cmd.BoolFlags[tok] = true
				i++
			}
		} else {
			cmd.Args = append(cmd.Args, tok)
			i++
		}
	}
	return cmd, nil
}

func IsAllowed(cmd ParsedCommand, allowed []string) bool {
	full := cmd.Base + " " + cmd.Sub
	for _, a := range allowed {
		if strings.EqualFold(a, full) {
			return true
		}
	}
	return false
}

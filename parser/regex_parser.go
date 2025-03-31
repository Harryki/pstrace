package parser

import (
	"regexp"
	"strings"
)

type RegexParser struct{}

var (
	funcDeclRegex = regexp.MustCompile(`(?i)^function\s+([A-Za-z_][A-Za-z0-9_-]*)\s*(\(\))?\s*\{?`)
)

func NewRegexParser() *RegexParser {
	return &RegexParser{}
}

func (p *RegexParser) ParseFunctions(content string) (map[string][]string, []string, error) {
	lines := strings.Split(content, "\n")
	funcBodies := map[string][]string{}
	funcNames := []string{}
	var currentFunc string

	var collecting bool
	var buffer []string
	var counter int
	var ignoring bool

	for _, line := range lines {
		lineTrim := strings.TrimSpace(line)

		// If line matches a function declaration pattern, we will collect function name and its full body
		if matches := funcDeclRegex.FindStringSubmatch(lineTrim); matches != nil {
			currentFunc = matches[1]

			funcNames = append(funcNames, currentFunc)
			// fmt.Printf("Parsed function: %s\n", currentFunc)
			collecting = true
			buffer = nil
		}

		if collecting {
			if strings.Contains(lineTrim, "<#") && strings.Contains(lineTrim, "#>") {
				continue // one-line block comment
			}
			if strings.HasPrefix(lineTrim, "<#") {
				ignoring = true
				continue
			}
			if strings.Contains(lineTrim, "#>") {
				ignoring = false
				continue
			}
			if ignoring || strings.HasPrefix(lineTrim, "#") {
				continue
			}
			if idx := strings.Index(lineTrim, "#"); idx != -1 {
				lineTrim = strings.TrimSpace(lineTrim[:idx])
			}
			if lineTrim == "" {
				continue
			}

			if strings.Contains(line, "{") {
				counter += strings.Count(line, "{")
			}
			if strings.Contains(line, "}") {
				counter -= strings.Count(line, "}")
			}

			if counter > 0 {
				buffer = append(buffer, lineTrim)
			}

			if counter == 0 && buffer != nil {
				funcBodies[currentFunc] = append([]string{}, buffer...)
				buffer = nil
				collecting = false
			}
		}
	}
	if counter != 0 {
		panic("Counter is negative when it should be 0!")
	}
	return funcBodies, funcNames, nil
}

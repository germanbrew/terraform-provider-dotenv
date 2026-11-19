// Package dotenv implements the parsing of the .env format.
//
// The parser is derived from github.com/direnv/direnv/v2/pkg/dotenv
// (MIT License, Copyright (c) 2019 zimbatm and contributors, see NOTICE),
// extended with options to control how variable references are resolved.
//
// There is no formal definition of the format but it has been introduced by
// https://github.com/bkeepers/dotenv which is thus canonical.
package dotenv

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	// ErrInvalidLine is returned when a line can not be parsed.
	ErrInvalidLine = errors.New("invalid line")
	// ErrUnclosedQuote is returned when a quoted value is never closed.
	ErrUnclosedQuote = errors.New("unclosed quoted value in .env file")
)

// line is the regexp matching a single line.
const line = `
\A
\s*
(?:|#.*|          # comment line
(?:export\s+)?    # optional export
([\w\.]+)         # key
(?:\s*=\s*|:\s+?) # separator
(                 # optional value begin
  '(?:\'|[^'])*'  #   single quoted value
  |               #   or
  "(?:\"|[^"])*"  #   double quoted value
  |               #   or
  [^\s#\n]+       #   unquoted value
)?                # value end
\s*
(?:\#.*)?         # optional comment
)
\z
`

var (
	linesRe = regexp.MustCompile(`[\r\n]+`)
	lineRe  = regexp.MustCompile(
		regexp.MustCompile(`\s+`).ReplaceAllLiteralString(
			regexp.MustCompile(`\s+# .*`).ReplaceAllLiteralString(line, ""), ""))
	escRe = regexp.MustCompile(`\\([^$])`)
)

// Options control how a dotenv document is parsed.
type Options struct {
	// LookupEnv resolves variable references (`$VAR`, `${VAR}`) that are not
	// defined in the document itself. If nil, such references resolve to an
	// empty string (or to their `${VAR:-default}` default).
	LookupEnv func(key string) (string, bool)

	// Base contains entries that are already defined before parsing starts,
	// e.g. from a previously parsed file. They can be referenced from the
	// document and are included in the result unless overridden.
	Base map[string]string
}

// Parse reads a string in the .env format and returns a map of the extracted
// key=values. References to variables not defined in the document are
// resolved from the process environment.
func Parse(data string) (map[string]string, error) {
	return ParseWithOptions(data, Options{LookupEnv: os.LookupEnv})
}

// ParseWithOptions reads a string in the .env format and returns a map of the
// extracted key=values.
//
// Ported from https://github.com/bkeepers/dotenv/blob/84f33f48107c492c3a99bd41c1059e7b4c1bb67a/lib/dotenv/parser.rb
func ParseWithOptions(data string, opts Options) (map[string]string, error) {
	p := parser{
		env:       make(map[string]string, len(opts.Base)),
		lookupEnv: opts.LookupEnv,
	}

	for k, v := range opts.Base {
		p.env[k] = v
	}

	var (
		inMultiline    bool
		multilineValue string
		quoteChar      byte
	)

	for _, line := range linesRe.Split(data, -1) {
		// Continue collecting a multi-line value
		if inMultiline {
			multilineValue += "\n" + line

			// Check if this line completes the multi-line value
			if strings.HasSuffix(strings.TrimSpace(line), string(quoteChar)) {
				match := lineRe.FindStringSubmatch(multilineValue)
				if len(match) > 1 && len(match[1]) > 0 {
					p.parseValue(match[1], match[2])
				}

				inMultiline = false
			}

			continue
		}

		if q, ok := startsMultiline(line); ok {
			quoteChar = q
			inMultiline = true
			multilineValue = line

			continue
		}

		match := lineRe.FindStringSubmatch(line)
		if match == nil {
			return nil, fmt.Errorf("%w: %s", ErrInvalidLine, line)
		}

		// commented or empty line
		if len(match[1]) == 0 {
			continue
		}

		p.parseValue(match[1], match[2])
	}

	// If we end with an unclosed multi-line value, return an error
	if inMultiline {
		return nil, ErrUnclosedQuote
	}

	return p.env, nil
}

// startsMultiline reports whether line opens a quoted value that is not closed
// on the same line, and returns the quote character.
func startsMultiline(line string) (byte, bool) {
	// Do not treat comment lines that happen to contain a quote as values.
	if strings.HasPrefix(strings.TrimLeft(line, " \t"), "#") {
		return 0, false
	}

	sepIdx := strings.IndexAny(line, "=:")
	if sepIdx <= 0 || sepIdx+1 >= len(line) {
		return 0, false
	}

	afterSep := strings.TrimLeft(line[sepIdx+1:], " \t")
	if len(afterSep) == 0 || (afterSep[0] != '"' && afterSep[0] != '\'') {
		return 0, false
	}

	if !unclosedQuote(afterSep) {
		return 0, false
	}

	return afterSep[0], true
}

// unclosedQuote reports whether value, which starts with a quote, has no
// unescaped closing quote.
func unclosedQuote(value string) bool {
	quote := value[0]

	for i := 1; i < len(value); i++ {
		switch value[i] {
		case '\\':
			if quote == '"' {
				i++
			}
		case quote:
			return false
		}
	}

	return true
}

type parser struct {
	env       map[string]string
	lookupEnv func(key string) (string, bool)
}

func (p *parser) parseValue(key string, value string) {
	if len(value) <= 1 {
		p.env[key] = value

		return
	}

	switch {
	case value[0] == '\'' && value[len(value)-1] == '\'':
		// single-quoted string, do not expand
		p.env[key] = value[1 : len(value)-1]

		return
	case value[0] == '"' && value[len(value)-1] == '"':
		value = value[1 : len(value)-1]
		value = expandNewLines(value)
		value = unescapeCharacters(value)
	}

	p.env[key] = p.expandEnv(value)
}

func unescapeCharacters(value string) string {
	return escRe.ReplaceAllString(value, "$1")
}

func expandNewLines(value string) string {
	value = strings.ReplaceAll(value, "\\n", "\n")
	value = strings.ReplaceAll(value, "\\r", "\r")

	return value
}

func (p *parser) expandEnv(value string) string {
	expander := func(ref string) string {
		envKey, defaultValue, hasDefault := strings.Cut(ref, ":-")

		if expanded, found := p.env[envKey]; found {
			return expanded
		}

		return p.fromEnvOrDefault(envKey, defaultValue, hasDefault)
	}

	// A backslash before a dollar sign escapes it, as in the canonical
	// implementation. unescapeCharacters deliberately leaves `\$` alone for this
	// step, and os.Expand knows nothing about escaping, so the escape has to be
	// resolved here or the backslash is kept and the variable expands anyway.
	var expanded strings.Builder

	for {
		i := strings.Index(value, `\$`)
		if i < 0 {
			break
		}

		expanded.WriteString(os.Expand(value[:i], expander))
		expanded.WriteString("$")

		value = value[i+2:]
	}

	expanded.WriteString(os.Expand(value, expander))

	return expanded.String()
}

func (p *parser) fromEnvOrDefault(envKey string, defaultValue string, hasDefault bool) string {
	var envValue string

	if p.lookupEnv != nil {
		envValue, _ = p.lookupEnv(envKey)
	}

	if len(envValue) == 0 && hasDefault {
		return defaultValue
	}

	return envValue
}
